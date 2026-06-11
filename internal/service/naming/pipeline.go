package naming

import (
	"context"
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
	TopNames       []analysis.NameResult
	ExcellentTable *ExcellentTable
	CharMap        map[string]*ent.Character
	FateData       *chronosfate.FateData
	// LastNameStrokes records the stroke values computed for the surname.
	LastNameStrokes [2]int
}

// Pipeline orchestrates the full name generation workflow:
// fate calculation → lucky wuge combos → character preload → scoring → output.
type Pipeline struct {
	db *repository.Repository
}

// NewPipeline creates a new name generation pipeline.
func NewPipeline(db *repository.Repository) *Pipeline {
	return &Pipeline{db: db}
}

// Generate runs the complete name generation pipeline.
func (p *Pipeline) Generate(ctx context.Context, req GenerateRequest) (*GenerateResult, error) {
	// 1. Calculate surname strokes from filter.
	strokes := surnameStrokes(req.Filter, req.LastName)

	// 2. Use pre-calculated fate data or compute it from birth info.
	fateData := req.FateData
	if fateData == nil {
		method := chronosfate.XiYongMethodBalance
		if req.Filter.XiYongMethod() == "geju" {
			method = chronosfate.XiYongMethodGeJu
		}
		var err error
		fateData, err = chronosfate.GetFateData(chronosfate.FateInput{
			Calendar:     v2.NewSolarCalendar(req.Born),
			Gender:       int(req.Sex),
			XiYongMethod: method,
		})
		if err != nil {
			log.Error("get fate data", err)
		}
	}

	// 3. Get lucky wuge combinations for this surname.
	lucky := wuge.GetLuckyByLastName(strokes[0], strokes[1])

	// 4. Preload characters needed by lucky combos.
	chars := preloadChars(ctx, p.db, lucky, req.Filter)

	// 5. Query poetry set (HasPoetry flag from DB).
	poetrySet := p.queryPoetrySet(ctx)

	// 6. Create rater with fate data and surname strokes.
	rater := rating.NewRaterWithStrokes(fateData, strokes[0], strokes[1])

	// 7. Score all candidates with 3-tier filter strictness fallback.
	table := NewExcellentTable()
	runScoring(ctx, lucky, rater, table, poetrySet, chars, req.Filter)

	if table.HeapLen() == 0 && req.Filter.FilterStrictness() != "relaxed" {
		req.Filter.SetFilterStrictness("moderate")
		chars = preloadChars(ctx, p.db, lucky, req.Filter)
		runScoring(ctx, lucky, rater, table, poetrySet, chars, req.Filter)
	}
	if table.HeapLen() == 0 && req.Filter.FilterStrictness() != "relaxed" {
		req.Filter.SetFilterStrictness("relaxed")
		chars = preloadChars(ctx, p.db, lucky, req.Filter)
		runScoring(ctx, lucky, rater, table, poetrySet, chars, req.Filter)
	}

	// 8. Finalize the heap into sorted entries.
	table.Finalize()

	// 9. Build character lookup map.
	charMap := buildCharMap(chars)

	// 10. Convert top entries to analysis.NameResult slices.
	surname := surnameStr(req.LastName)
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
	}, nil
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
