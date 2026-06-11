package fate

import (
	"context"
	"time"

	"github.com/babyname/fate/v4/internal/naming"
	namesvc "github.com/babyname/fate/v4/internal/service/naming"
	"github.com/babyname/fate/v4/ent"
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
	// Filter is the filter option. If nil, DefaultFilter() is used.
	Filter *FilterOption
}

// GenerateResult is the output of synchronous name generation.
type GenerateResult struct {
	TopNames       []NameResult
	ExcellentTable ExcellentTable
	CharMap        map[string]CharInfo
	FateData       *FateData
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

	// 2. Build filter from option or use default.
	var flt Filter
	if req.Filter != nil {
		flt = NewFilter(*req.Filter)
	} else {
		flt = DefaultFilter()
	}

	// 3. Create pipeline and run generation.
	pipe := namesvc.NewPipeline(f.db)
	result, err := pipe.Generate(ctx, namesvc.GenerateRequest{
		LastName: lastName,
		Born:     req.Born,
		Sex:      naming.Sex(req.Sex),
		Filter:   flt,
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
		ExcellentTable: *result.ExcellentTable,
		CharMap:        charMap,
		FateData:       result.FateData,
	}, nil
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
