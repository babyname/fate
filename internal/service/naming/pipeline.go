package naming

import (
	"context"
	"sync"
	"time"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/analysis"
	"github.com/babyname/fate/v4/internal/chronosfate"
	"github.com/babyname/fate/v4/internal/filter"
	"github.com/babyname/fate/v4/internal/log"
	"github.com/babyname/fate/v4/internal/naming"
	"github.com/babyname/fate/v4/internal/rating"
	"github.com/babyname/fate/v4/internal/repository"
	"github.com/babyname/fate/v4/internal/wuge"
	v2 "github.com/godcong/chronos/v2"
)

// GenerateRequest carries all inputs required for a name generation run.
// FateData may be nil (the pipeline will compute it from Born/Sex).
// If FateData is nil, Born must be set. Otherwise Born may be zero.
type GenerateRequest struct {
	LastName [2]*ent.Character
	Born     time.Time
	Sex      naming.Sex
	Filter   filter.Filter
	FateData *chronosfate.FateData
}

// GenerateResult holds the output of a completed generation run.
type GenerateResult struct {
	TopNames        []analysis.NameResult
	ExcellentTable  *ExcellentTable
	CharMap         map[string]*ent.Character
	FateData        *chronosfate.FateData
	LastNameStrokes [2]int
}

// StreamBatch is delivered by GenerateStream as top results become available.
type StreamBatch struct {
	TopResults []analysis.NameResult
	Total      int
	IsFinal    bool
	// Result is set only on the final batch (IsFinal=true) and provides the
	// complete GenerationResult for pagination.
	Result *GenerateResult
}

// Pipeline orchestrates the full name generation workflow.
type Pipeline struct {
	db *repository.Repository
}

// NewPipeline creates a new name generation pipeline.
func NewPipeline(db *repository.Repository) *Pipeline {
	return &Pipeline{db: db}
}

// Generate runs the complete name generation pipeline synchronously.
func (p *Pipeline) Generate(ctx context.Context, req GenerateRequest) (*GenerateResult, error) {
	// 1. Calculate surname strokes from filter.
	strokes := surnameStrokes(req.Filter, req.LastName)

	// 2. Use pre-calculated fate data or compute it from birth info.
	fateData := p.resolveFateData(req)

	// 3. Get lucky wuge combinations for this surname.
	lucky := wuge.GetLuckyByLastName(strokes[0], strokes[1])

	// 4. Preload characters needed by lucky combos.
	chars := preloadChars(ctx, p.db, lucky, req.Filter)

	// 5. Query poetry set.
	poetrySet := p.queryPoetrySet(ctx)

	// 6. Create rater with fate data and surname strokes.
	rater := rating.NewRaterWithStrokes(fateData, strokes[0], strokes[1])

	// 7. Score all candidates with 3-tier filter strictness fallback.
	table := NewExcellentTable()
	table, chars = p.runScoringWithFallback(ctx, lucky, rater, table, poetrySet, chars, req)

	// 8. Finalize the heap into sorted entries.
	table.Finalize()

	// 9. Build character lookup map.
	charMap := buildCharMap(chars)

	// 10. Convert top entries to analysis.NameResult slices.
	result := p.buildResults(table, charMap, strokes, fateData, surnameStr(req.LastName))
	return result, nil
}

// GenerateStream runs the name generation pipeline in parallel and sends
// batches of top results via callback as workers complete.
func (p *Pipeline) GenerateStream(ctx context.Context, req GenerateRequest, onBatch func(StreamBatch)) error {
	strokes := surnameStrokes(req.Filter, req.LastName)
	fateData := p.resolveFateData(req)
	lucky := wuge.GetLuckyByLastName(strokes[0], strokes[1])
	chars := preloadChars(ctx, p.db, lucky, req.Filter)
	poetrySet := p.queryPoetrySet(ctx)
	rater := rating.NewRaterWithStrokes(fateData, strokes[0], strokes[1])
	surname := surnameStr(req.LastName)

	table := NewExcellentTable()

	var batchMu sync.Mutex
	batchSeq := 0

	cb := func() {
		batchMu.Lock()
		seq := batchSeq + 1
		batchSeq = seq
		batchMu.Unlock()

		if seq == 1 || seq%2 == 0 {
			snapshot := heapPreview(table, 10, surname, strokes)
			onBatch(StreamBatch{
				TopResults: snapshot,
				Total:      table.HeapLen(),
				IsFinal:    false,
			})
		}
	}

	runScoringParallel(ctx, lucky, rater, table, poetrySet, chars, req.Filter, cb)

	// Fallback tiers if no results.
	strictness := req.Filter.FilterStrictness()
	if table.HeapLen() == 0 && strictness != "relaxed" {
		req.Filter.SetFilterStrictness("moderate")
		chars = preloadChars(ctx, p.db, lucky, req.Filter)
		runScoringParallel(ctx, lucky, rater, table, poetrySet, chars, req.Filter, cb)
	}
	if table.HeapLen() == 0 && strictness != "relaxed" {
		req.Filter.SetFilterStrictness("relaxed")
		chars = preloadChars(ctx, p.db, lucky, req.Filter)
		runScoringParallel(ctx, lucky, rater, table, poetrySet, chars, req.Filter, cb)
	}

	table.Finalize()
	charMap := buildCharMap(chars)
	result := p.buildResults(table, charMap, strokes, fateData, surname)

	onBatch(StreamBatch{
		TopResults: result.TopNames,
		Total:      table.Len(),
		IsFinal:    true,
		Result:     result,
	})
	return nil
}

// resolveFateData returns the fate data, computing it from birth info if needed.
func (p *Pipeline) resolveFateData(req GenerateRequest) *chronosfate.FateData {
	if req.FateData != nil {
		return req.FateData
	}
	method := chronosfate.XiYongMethodBalance
	if req.Filter.XiYongMethod() == "geju" {
		method = chronosfate.XiYongMethodGeJu
	}
	fateData, err := chronosfate.GetFateData(chronosfate.FateInput{
		Calendar:     v2.NewSolarCalendar(req.Born),
		Gender:       int(req.Sex),
		XiYongMethod: method,
	})
	if err != nil {
		log.Error("get fate data", err)
	}
	return fateData
}

// runScoringWithFallback runs scoring with up to 3 filter strictness tiers.
// Returns the populated table and the final set of preloaded characters.
func (p *Pipeline) runScoringWithFallback(
	ctx context.Context,
	lucky []wuge.WuGeResult,
	rater *rating.Rater,
	table *ExcellentTable,
	poetrySet map[string]bool,
	chars map[int][]*ent.Character,
	req GenerateRequest,
) (*ExcellentTable, map[int][]*ent.Character) {
	runScoring(ctx, lucky, rater, table, poetrySet, chars, req.Filter)

	strictness := req.Filter.FilterStrictness()
	if table.HeapLen() == 0 && strictness != "relaxed" {
		req.Filter.SetFilterStrictness("moderate")
		chars = preloadChars(ctx, p.db, lucky, req.Filter)
		runScoring(ctx, lucky, rater, table, poetrySet, chars, req.Filter)
	}
	if table.HeapLen() == 0 && strictness != "relaxed" {
		req.Filter.SetFilterStrictness("relaxed")
		chars = preloadChars(ctx, p.db, lucky, req.Filter)
		runScoring(ctx, lucky, rater, table, poetrySet, chars, req.Filter)
	}

	return table, chars
}

// buildResults converts the final excellent table into public result types.
func (p *Pipeline) buildResults(table *ExcellentTable, charMap map[string]*ent.Character, strokes [2]int, fateData *chronosfate.FateData, surname string) *GenerateResult {
	topEntries := table.TopN(10)
	topResults := make([]analysis.NameResult, 0, len(topEntries))
	for i, e := range topEntries {
		c1 := charMap[e.Char1]
		c2 := charMap[e.Char2]
		if c1 == nil || c2 == nil {
			continue
		}
		nr := analysis.BuildNameResult(i+1, surname, c1, c2, strokes[0], strokes[1], fateData)
		topResults = append(topResults, nr)
	}

	return &GenerateResult{
		TopNames:        topResults,
		ExcellentTable:  table,
		CharMap:         charMap,
		FateData:        fateData,
		LastNameStrokes: strokes,
	}
}

// heapPreview returns a lightweight preview of top N entries from the heap
// before finalization. Uses only ExcellentEntry-level data (not full NameResult).
func heapPreview(table *ExcellentTable, n int, surname string, strokes [2]int) []analysis.NameResult {
	entries := table.SortedPreview(n)
	if len(entries) == 0 {
		return nil
	}
	results := make([]analysis.NameResult, 0, len(entries))
	for i, e := range entries {
		results = append(results, analysis.NameResult{
			Rank:     i + 1,
			Surname:  surname,
			FullName: surname + e.Char1 + e.Char2,
			Score:    e.Score,
			Grade:    e.Grade,
			HasPoetry: e.HasPoetry,
			Char1: analysis.CharInfo{
				Char:      e.Char1,
				WuXing:    e.WuXing1,
				HasPoetry: e.HasPoetry,
			},
			Char2: analysis.CharInfo{
				Char:      e.Char2,
				WuXing:    e.WuXing2,
				HasPoetry: e.HasPoetry,
			},
		})
	}
	return results
}

// queryPoetrySet loads character IDs with HasPoetry=true from the database.
func (p *Pipeline) queryPoetrySet(ctx context.Context) map[string]bool {
	poetrySet := make(map[string]bool)
	chars, err := p.db.QueryPoetryChars(ctx)
	if err != nil {
		log.Error("query poetry chars failed", err)
	} else {
		log.Info("query poetry chars", "count", len(chars))
		for _, ch := range chars {
			poetrySet[ch] = true
		}
	}
	return poetrySet
}

// surnameStr joins the surname character(s) into a string.
func surnameStr(lastName [2]*ent.Character) string {
	s := ""
	if lastName[0] != nil {
		s = lastName[0].Char
	}
	if lastName[1] != nil && lastName[1].Char != "" {
		s += lastName[1].Char
	}
	return s
}

// surnameStrokes computes surname stroke values using the active filter mode.
func surnameStrokes(flt filter.Filter, lastName [2]*ent.Character) [2]int {
	sg := strokeGetFromFilterType(flt.FilterType())
	var strokes [2]int
	strokes[0] = sg(lastName[0])
	if lastName[1] != nil {
		strokes[1] = sg(lastName[1])
	}
	return strokes
}

func strokeGetFromFilterType(ct filter.CharacterFilterType) func(c *ent.Character) int {
	if ct == filter.CharacterFilterTypeDefault {
		return func(c *ent.Character) int { return c.ScienceStroke }
	}
	if ct.HasType(filter.CharacterFilterTypeChs) {
		return func(c *ent.Character) int { return c.SimplifiedStroke }
	}
	if ct.HasType(filter.CharacterFilterTypeCht) {
		return func(c *ent.Character) int { return c.TraditionalStroke }
	}
	if ct.HasType(filter.CharacterFilterTypeKangxi) {
		return func(c *ent.Character) int { return c.KangxiStroke }
	}
	return func(c *ent.Character) int { return c.ScienceStroke }
}

// buildCharMap creates a deduplicated lookup map from preloaded characters.
func buildCharMap(chars map[int][]*ent.Character) map[string]*ent.Character {
	charMap := make(map[string]*ent.Character)
	for _, charList := range chars {
		for _, c := range charList {
			if _, ok := charMap[c.Char]; !ok {
				charMap[c.Char] = c
			}
		}
	}
	return charMap
}
