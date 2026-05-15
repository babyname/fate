// Package filter provides character filtering for name generation.
package filter

import (
	"errors"
	"strings"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/internal/wuge"
	"github.com/babyname/fate/internal/wuxing"
)

var defaultFilter = newFilter()

// Filter defines the interface for name filtering operations including
// character query filtering, stroke-based filtering, WuGe luck checks,
// gender filtering, and display name formatting.
type Filter interface {
	FilterType() CharacterFilterType
	PoetryMode() int
	CustomFilter(name string, v any) error
	// QueryCharacterFilter applies character-level predicates to a CharacterQuery.
	QueryCharacterFilter(query *ent.CharacterQuery) *ent.CharacterQuery
	// QueryRegularFilter applies regular-character predicates to a CharacterQuery.
	QueryRegularFilter(query *ent.CharacterQuery) *ent.CharacterQuery
	// QueryStrokeFilter returns a function that filters characters by the given stroke count.
	QueryStrokeFilter(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery
	// CheckSkipSexFilter returns true if the WuGe result should be skipped based on sex luck.
	CheckSkipSexFilter(lucky *wuge.WuGeResult) bool
	// CheckSkipDaYanFilter returns true if the WuGe result should be skipped based on DaYan luck.
	CheckSkipDaYanFilter(lucky *wuge.WuGeResult) bool
	// CheckSkipWuXingFilter returns true if the SanCai WuXing combination should be skipped.
	CheckSkipWuXingFilter(ge int, ge2 int, ge3 int) bool
	// CheckSkipStrokeNumberScope returns true if any of the given stroke values fall outside scope.
	CheckSkipStrokeNumberScope(stroke ...int) bool
	// CheckCharacter returns true if the character passes all character-level filters.
	CheckCharacter(c *ent.Character) bool
	// GetCharacterStroke returns the stroke count for a character based on the active StrokeMode
	// or CharacterFilterType.
	GetCharacterStroke(c *ent.Character) int
	// CheckGenderFilter returns true if the character's gender hint matches the gender filter.
	CheckGenderFilter(c *ent.Character) bool
	// GetDisplayName returns the character string formatted according to the active NameStyle.
	GetDisplayName(c *ent.Character) string
}

// filter is the concrete implementation of the Filter interface.
type filter struct {
	characterFilterType        CharacterFilterType
	strokeMode                 StrokeMode
	nameStyle                  NameStyle
	genderFilter               string
	poetryMode                 int
	requireCharacters          []string
	avoidPinyin                []string
	checkSkipSexFilter         func(lucky *wuge.WuGeResult) bool
	checkSkipDaYanFilter       func(lucky *wuge.WuGeResult) bool
	checkSkipWuXingFilter      func(ge int, ge2 int, ge3 int) bool
	checkSkipStrokeNumberScope func(stroke []int) bool
	checkCharacter             func(c *ent.Character) bool
	checkGenderFilter          func(c *ent.Character) bool
	getDisplayName             func(c *ent.Character) string
	queryCharacterFilter       func(query *ent.CharacterQuery) *ent.CharacterQuery
	queryRegularFilter         func(query *ent.CharacterQuery) *ent.CharacterQuery
	queryStrokeFilter          func(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery
}

// FilterType returns the active CharacterFilterType.
func (f *filter) FilterType() CharacterFilterType {
	return f.characterFilterType
}

func (f *filter) PoetryMode() int {
	return f.poetryMode
}

// CheckSkipStrokeNumberScope checks if any stroke value falls outside the configured scope.
func (f *filter) CheckSkipStrokeNumberScope(stroke ...int) bool {
	return f.checkSkipStrokeNumberScope(stroke)
}

// QueryStrokeFilter returns a function that filters a CharacterQuery by the given stroke count.
func (f *filter) QueryStrokeFilter(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery {
	return f.queryStrokeFilter(stroke)
}

// QueryRegularFilter applies regular-character predicates to the query.
func (f *filter) QueryRegularFilter(query *ent.CharacterQuery) *ent.CharacterQuery {
	return f.queryRegularFilter(query)
}

// GetCharacterStroke returns the stroke count for the character based on the
// active StrokeMode. If StrokeMode is not explicitly set (StrokeModeScience is
// the zero value and also the default), it falls back to CharacterFilterType.
func (f *filter) GetCharacterStroke(c *ent.Character) int {
	switch f.strokeMode {
	case StrokeModeSimplified:
		return c.SimplifiedStroke
	case StrokeModeTraditional:
		return c.TraditionalStroke
	case StrokeModeKangxi:
		return c.KangxiStroke
	case StrokeModeScience:
		return c.ScienceStroke
	default:
		switch f.characterFilterType {
		case CharacterFilterTypeChs:
			return c.SimplifiedStroke
		case CharacterFilterTypeCht:
			return c.TraditionalStroke
		case CharacterFilterTypeKangxi:
			return c.KangxiStroke
		case CharacterFilterTypeDefault:
			fallthrough
		default:
			return c.ScienceStroke
		}
	}
}

// CustomFilter applies a named custom filter. Returns ErrNotImplemented if the
// filter name is not recognized or the value is nil.
func (f *filter) CustomFilter(name string, v any) error {
	if v == nil {
		return ErrNotImplemented
	}
	switch name {
	default:
		return ErrNotImplemented
	}
}

// CheckSkipWuXingFilter returns true if the SanCai WuXing combination should be skipped.
func (f *filter) CheckSkipWuXingFilter(ge int, ge2 int, ge3 int) bool {
	return f.checkSkipWuXingFilter(ge, ge2, ge3)
}

// CheckSkipDaYanFilter returns true if the WuGe result should be skipped based on DaYan luck.
func (f *filter) CheckSkipDaYanFilter(lucky *wuge.WuGeResult) bool {
	return f.checkSkipDaYanFilter(lucky)
}

// CheckSkipSexFilter returns true if the WuGe result should be skipped based on sex luck.
func (f *filter) CheckSkipSexFilter(lucky *wuge.WuGeResult) bool {
	return f.checkSkipSexFilter(lucky)
}

// QueryCharacterFilter applies character-level predicates to the query.
func (f *filter) QueryCharacterFilter(query *ent.CharacterQuery) *ent.CharacterQuery {
	return f.queryCharacterFilter(query)
}

// CheckCharacter returns true if the character passes all character-level filters.
func (f *filter) CheckCharacter(c *ent.Character) bool {
	return f.checkCharacter(c)
}

// CheckGenderFilter returns true if the character's gender hint is compatible
// with the configured gender filter.
func (f *filter) CheckGenderFilter(c *ent.Character) bool {
	return f.checkGenderFilter(c)
}

// GetDisplayName returns the character string formatted according to the
// active NameStyle. For NameStyleSimplified it returns the character as-is
// (characters are stored in standard form). For NameStyleTraditional it
// returns the character as-is; proper simplified/traditional conversion
// requires edge traversal at a higher level.
func (f *filter) GetDisplayName(c *ent.Character) string {
	return f.getDisplayName(c)
}

// ErrNotImplemented is returned by CustomFilter when the requested filter
// name is not recognized or the value is nil.
var ErrNotImplemented = errors.New("not implemented")

// DefaultFilter returns the shared default Filter instance.
func DefaultFilter() Filter {
	return defaultFilter
}

// newFilter creates a filter with no-op defaults for all filter functions.
func newFilter() *filter {
	return &filter{
		characterFilterType: CharacterFilterTypeDefault,
		strokeMode:          StrokeModeScience,
		nameStyle:           NameStyleSimplified,
		genderFilter:        "",
		checkSkipSexFilter: func(_ *wuge.WuGeResult) bool {
			return false
		},
		checkSkipDaYanFilter: func(_ *wuge.WuGeResult) bool {
			return false
		},
		checkSkipWuXingFilter: func(_ int, _ int, _ int) bool {
			return false
		},
		checkSkipStrokeNumberScope: func(_ []int) bool {
			return false
		},
		checkCharacter: func(_ *ent.Character) bool {
			return true
		},
		checkGenderFilter: func(_ *ent.Character) bool {
			return true
		},
		getDisplayName: func(c *ent.Character) string {
			return c.Char
		},
		queryCharacterFilter: func(query *ent.CharacterQuery) *ent.CharacterQuery {
			return query
		},
		queryRegularFilter: func(query *ent.CharacterQuery) *ent.CharacterQuery {
			return query
		},
		queryStrokeFilter: func(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery {
			return func(query *ent.CharacterQuery) *ent.CharacterQuery {
				return query.Where(character.ScienceStrokeEQ(stroke))
			}
		},
	}
}

// NewFilter creates a new Filter from the given FilterOption. All filter
// functions are composed from the option fields so that calling the Filter
// methods applies the corresponding logic.
func NewFilter(fo FilterOption) Filter {
	f := newFilter()

	if fo.StrokeMode != StrokeModeScience {
		f.strokeMode = fo.StrokeMode
	}
	f.nameStyle = fo.NameStyle
	f.genderFilter = fo.GenderFilter
	f.requireCharacters = fo.RequireCharacters
	f.avoidPinyin = fo.AvoidPinyin
	f.poetryMode = fo.PoetryMode

	if fo.SexFilter {
		f.checkSkipSexFilter = func(lucky *wuge.WuGeResult) bool {
			return !lucky.ZongLucky
		}
	}

	if fo.CharacterFilter {
		f.characterFilterType = fo.CharacterFilterType
		switch fo.CharacterFilterType {
		case CharacterFilterTypeChs:
			f.queryCharacterFilter = characterTypeFilterCHS(fo.MinStroke, fo.MaxStroke)
			f.queryStrokeFilter = strokeFilterCHS
		case CharacterFilterTypeCht:
			f.queryCharacterFilter = characterTypeFilterCHT(fo.MinStroke, fo.MaxStroke)
			f.queryStrokeFilter = strokeFilterCHT
		case CharacterFilterTypeKangxi:
			f.queryCharacterFilter = characterTypeFilterKX(fo.MinStroke, fo.MaxStroke)
			f.queryStrokeFilter = strokeFilterKX
		case CharacterFilterTypeDefault:
			fallthrough
		default:
			f.queryCharacterFilter = characterTypeFilterDefault(fo.MinStroke, fo.MaxStroke)
			f.queryStrokeFilter = strokeFilterDefault
		}
	}

	if fo.MinStroke != 0 || fo.MaxStroke != 0 {
		f.checkSkipStrokeNumberScope = func(strokes []int) bool {
			for i := range strokes {
				if fo.MinStroke != 0 && strokes[i] < fo.MinStroke {
					return true
				}
				if fo.MaxStroke != 0 && strokes[i] > fo.MaxStroke {
					return true
				}
			}
			return false
		}
	}

	if fo.DaYanFilter {
		f.checkSkipDaYanFilter = daYanFilter
	}

	if fo.WuXingFilter {
		f.checkSkipWuXingFilter = wuXingFilter
	}

	if fo.RegularFilter {
		f.queryRegularFilter = regularFilter
	}

	if fo.GenderFilter != "" {
		f.checkGenderFilter = func(c *ent.Character) bool {
			if c.GenderHint == "neutral" {
				return true
			}
			return c.GenderHint == fo.GenderFilter
		}
	}

	f.getDisplayName = func(c *ent.Character) string {
		switch f.nameStyle {
		case NameStyleTraditional:
			if c.IsTraditional {
				return c.Char
			}
			return c.Char
		case NameStyleSimplified:
			fallthrough
		default:
			if c.IsSimplified {
				return c.Char
			}
			return c.Char
		}
	}

	originalCheck := f.checkCharacter
	f.checkCharacter = func(c *ent.Character) bool {
		if !originalCheck(c) {
			return false
		}

		for _, avoidChar := range fo.AvoidCharacters {
			if c.Char == avoidChar {
				return false
			}
		}

		if fo.CommonLevelFilter {
			if c.CommonLevel < fo.MinCommonLevel || c.CommonLevel > fo.MaxCommonLevel {
				return false
			}
		}

		if len(fo.PreferredWuXing) > 0 {
			found := false
			for _, wx := range fo.PreferredWuXing {
				if c.WuXing == wx {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		if len(fo.AvoidWuXing) > 0 {
			for _, wx := range fo.AvoidWuXing {
				if c.WuXing == wx {
					return false
				}
			}
		}

		if len(fo.PreferredRadicals) > 0 {
			found := false
			for _, rad := range fo.PreferredRadicals {
				if c.Radical == rad {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		if len(fo.AvoidRadicals) > 0 {
			for _, rad := range fo.AvoidRadicals {
				if c.Radical == rad {
					return false
				}
			}
		}

		if len(fo.AvoidPinyin) > 0 {
			for _, p := range c.Pinyin {
				for _, ap := range fo.AvoidPinyin {
					if strings.EqualFold(p, ap) {
						return false
					}
				}
			}
		}

		if len(fo.RequireCharacters) > 0 {
			for _, rc := range fo.RequireCharacters {
				if c.Char == rc {
					return true
				}
			}
		}

		return true
	}

	originalQuery := f.queryCharacterFilter
	f.queryCharacterFilter = func(query *ent.CharacterQuery) *ent.CharacterQuery {
		q := originalQuery(query)

		if fo.CommonLevelFilter {
			q = q.Where(
				character.CommonLevelGTE(fo.MinCommonLevel),
				character.CommonLevelLTE(fo.MaxCommonLevel),
			)
		}

		if len(fo.AvoidCharacters) > 0 {
			q = q.Where(character.CharNotIn(fo.AvoidCharacters...))
		}

		if len(fo.PreferredWuXing) > 0 {
			q = q.Where(character.WuXingIn(fo.PreferredWuXing...))
		}

		if len(fo.AvoidWuXing) > 0 {
			q = q.Where(character.WuXingNotIn(fo.AvoidWuXing...))
		}

		if fo.GenderFilter != "" {
			q = q.Where(
				character.Or(
					character.GenderHintEQ(fo.GenderFilter),
					character.GenderHintEQ("neutral"),
				),
			)
		}

		return q
	}

	return f
}

// characterTypeFilterDefault returns a query filter that restricts characters
// by ScienceStroke range.
func characterTypeFilterDefault(minStroke, maxStroke int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.ScienceStrokeGTE(minStroke), character.ScienceStrokeLTE(maxStroke))
	}
}

// strokeFilterDefault returns a query filter that matches characters with the
// given ScienceStroke value.
func strokeFilterDefault(s int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.ScienceStrokeEQ(s))
	}
}

// characterTypeFilterCHS returns a query filter that restricts characters by
// SimplifiedStroke range.
func characterTypeFilterCHS(minStroke, maxStroke int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.SimplifiedStrokeGTE(minStroke), character.SimplifiedStrokeLTE(maxStroke))
	}
}

// strokeFilterCHS returns a query filter that matches characters with the
// given SimplifiedStroke value.
func strokeFilterCHS(s int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.SimplifiedStrokeEQ(s))
	}
}

// characterTypeFilterCHT returns a query filter that restricts characters by
// TraditionalStroke range.
func characterTypeFilterCHT(minStroke, maxStroke int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.TraditionalStrokeGTE(minStroke), character.TraditionalStrokeLTE(maxStroke))
	}
}

// strokeFilterCHT returns a query filter that matches characters with the
// given TraditionalStroke value.
func strokeFilterCHT(s int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.TraditionalStrokeEQ(s))
	}
}

// characterTypeFilterKX returns a query filter that restricts characters by
// KangxiStroke range.
func characterTypeFilterKX(minStroke, maxStroke int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.KangxiStrokeGTE(minStroke), character.KangxiStrokeLTE(maxStroke))
	}
}

// strokeFilterKX returns a query filter that matches characters with the
// given KangxiStroke value.
func strokeFilterKX(s int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.KangxiStrokeEQ(s))
	}
}

// regularFilter applies predicates that restrict the query to regular
// (common-use) and nameable characters.
func regularFilter(query *ent.CharacterQuery) *ent.CharacterQuery {
	return query.Where(character.RegularEQ(true), character.NameableEQ(true))
}

// wuXingFilter returns true if the SanCai WuXing combination should be skipped
// (i.e. the combination does not pass the level-5 check).
func wuXingFilter(ge int, ge2 int, ge3 int) bool {
	sc := wuxing.NewSanCai(ge, ge2, ge3)
	return !sc.Check(5)
}

// daYanFilter returns true if the WuGe result should be skipped because not all
// four Ge positions (Di, Ren, Wai, Zong) are lucky.
func daYanFilter(lucky *wuge.WuGeResult) bool {
	return !isLucky(wuge.Find(lucky.DiGe).Lucky) ||
		!isLucky(wuge.Find(lucky.RenGe).Lucky) ||
		!isLucky(wuge.Find(lucky.WaiGe).Lucky) ||
		!isLucky(wuge.Find(lucky.ZongGe).Lucky)
}

// isLucky returns true if the luck description contains the character 吉.
func isLucky(s string) bool {
	return strings.Contains(s, "吉")
}
