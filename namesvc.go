package fate

import (
	"context"
	"time"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/naming"
	namesvc "github.com/babyname/fate/v4/internal/service/naming"
	"github.com/babyname/fate/v4/ranking"
)

// GenerateRequest is the public input for synchronous name generation.
// It is a simpler alternative to the Session state machine.
type GenerateRequest struct {
	// LastName is the surname, e.g. "张" or "欧阳".
	LastName string
	// Born is the birth time.
	Born time.Time
	// Sex is SexBoy or SexGirl.
	Sex Sex
	// Filter is the filter option. If nil, DefaultFilterOption() is applied:
	// a bounded, Web-equivalent configuration. DefaultFilter() is never used
	// by Generate because its unbounded stroke range does not converge.
	Filter *FilterOption
	// RecommendationWeight opts in to weighting traditional and sound-and-meaning scores.
	RecommendationWeight *ranking.Weight
}

// GenerateResult is the output of synchronous name generation.
type GenerateResult struct {
	TopNames       []NameResult
	ExcellentTable ExcellentTable
	CharMap        map[string]CharInfo
	FateData       *FateData
	catalogue      *namesvc.GenerateResult
}

// CandidateCount returns the number of generated names available for paging.
func (r *GenerateResult) CandidateCount() int {
	if r.catalogue == nil {
		return len(r.TopNames)
	}
	return r.catalogue.ExcellentTable.Len()
}

// Candidates returns up to limit scored candidates starting at the zero-based offset.
func (r *GenerateResult) Candidates(offset, limit int) []NameResult {
	if r.catalogue != nil {
		return r.catalogue.Candidates(offset, limit)
	}
	if offset < 0 || offset >= len(r.TopNames) || limit <= 0 {
		return []NameResult{}
	}
	return r.TopNames[offset : offset+min(limit, len(r.TopNames)-offset)]
}

// Candidate looks up the full scored result for a generated pair of given-name characters.
func (r *GenerateResult) Candidate(first, second string) (NameResult, bool) {
	for _, name := range r.TopNames {
		if name.Char1.Char == first && name.Char2.Char == second {
			return name, true
		}
	}
	if r.catalogue != nil {
		return r.catalogue.Candidate(first, second)
	}
	return NameResult{}, false
}

// Generate runs the full name generation pipeline synchronously.
// It is a convenience method that replaces the 5-step async dance
// (create Session → Start → poll State → Wait → read Output).
func (f *fateImpl) Generate(ctx context.Context, req GenerateRequest) (*GenerateResult, error) {
	// 1. Lookup surname characters in DB.
	last := splitSurname(req.LastName)
	lastName, err := f.db.QueryLastName(ctx, last)
	if err != nil {
		return nil, err
	}

	// 2. Build filter. A nil option resolves to the bounded default.
	flt := NewFilter(resolveFilterOption(req.Filter))

	// 3. Create pipeline and run generation.
	pipe := namesvc.NewPipeline(f.db)
	result, err := pipe.Generate(ctx, namesvc.GenerateRequest{
		LastName:             lastName,
		Born:                 req.Born,
		Sex:                  naming.Sex(req.Sex),
		Filter:               flt,
		RecommendationWeight: req.RecommendationWeight,
	})
	if err != nil {
		return nil, err
	}

	// 4. Convert internal CharMap to public CharInfo type.
	charMap := make(map[string]CharInfo, len(result.CharMap))
	for ch, c := range result.CharMap {
		charMap[ch] = charInfoFromCharacter(c)
	}

	return &GenerateResult{
		TopNames:       result.TopNames,
		ExcellentTable: result.ExcellentTable,
		CharMap:        charMap,
		FateData:       result.FateData,
		catalogue:      result,
	}, nil
}

// StreamBatch is delivered during streaming generation.
type StreamBatch = namesvc.StreamBatch

// GenerateStream runs the pipeline in parallel, sending batches of top results
// via onBatch as workers complete. The last call always has IsFinal=true.
func (f *fateImpl) GenerateStream(ctx context.Context, req GenerateRequest, onBatch func(StreamBatch)) error {
	last := splitSurname(req.LastName)
	lastName, err := f.db.QueryLastName(ctx, last)
	if err != nil {
		return err
	}

	flt := NewFilter(resolveFilterOption(req.Filter))

	pipe := namesvc.NewPipeline(f.db)

	// Wrap the public onBatch to convert internal StreamBatch to public type.
	// Since StreamBatch is a type alias, no conversion is needed.
	return pipe.GenerateStream(ctx, namesvc.GenerateRequest{
		LastName:             lastName,
		Born:                 req.Born,
		Sex:                  naming.Sex(req.Sex),
		Filter:               flt,
		RecommendationWeight: req.RecommendationWeight,
	}, func(batch namesvc.StreamBatch) {
		onBatch(batch)
	})
}

// charInfoFromCharacter converts an ent.Character to the public CharInfo type.
func charInfoFromCharacter(c *ent.Character) CharInfo {
	pinyin := ""
	if len(c.Pinyin) > 0 {
		pinyin = c.Pinyin[0]
	}
	traditionalChar := ""
	if c.IsTraditional {
		traditionalChar = c.Char
	}
	return CharInfo{
		Char:              c.Char,
		TraditionalChar:   traditionalChar,
		Pinyin:            pinyin,
		WuXing:            c.WuXing,
		SimplifiedStroke:  c.SimplifiedStroke,
		TraditionalStroke: c.TraditionalStroke,
		ScienceStroke:     c.ScienceStroke,
		KangxiStroke:      c.KangxiStroke,
		Radical:           c.Radical,
		Meaning:           c.Meaning,
		IsXiYong:          false, // XiYong is per-name-triad, not per-character
		HasPoetry:         c.HasPoetry,
	}
}

// splitSurname splits a surname string into a [2]string array.
// For single-character surnames like "张", returns ["张", ""].
// For compound surnames like "欧阳", returns ["欧", "阳"].
func splitSurname(surname string) [2]string {
	var result [2]string
	runes := []rune(surname)
	if len(runes) >= 1 {
		result[0] = string(runes[0])
	}
	if len(runes) >= 2 {
		result[1] = string(runes[1])
	}
	return result
}
