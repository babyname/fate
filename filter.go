package fate

import (
	"errors"
	"strings"

	"github.com/babyname/fate/core"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/internal/dayan"
)

var defaultFilter = newFilter()

// Filter prefix with Query are used db search, others are the default data check
type Filter interface {
	FilterType() CharacterFilterType
	CustomFilter(name string, v any) error
	QueryCharacterFilter(query *ent.CharacterQuery) *ent.CharacterQuery
	QueryRegularFilter(query *ent.CharacterQuery) *ent.CharacterQuery
	QueryStrokeFilter(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery
	CheckSkipSexFilter(lucky *ent.WuGeLucky) bool
	CheckSkipDaYanFilter(lucky *ent.WuGeLucky) bool
	CheckSkipWuXingFilter(ge int, ge2 int, ge3 int) bool
	CheckSkipStrokeNumberScope(stroke ...int) bool
}

type CharacterType interface {
}

type QueryFilter interface {
	Query(name string, v ...any) func(query *ent.CharacterQuery) *ent.CharacterQuery
}

type CheckFilter interface {
	Check(name string, v ...any) bool
}

type filter struct {
	characterFilterType        CharacterFilterType
	checkSkipSexFilter         func(lucky *ent.WuGeLucky) bool
	checkSkipDaYanFilter       func(lucky *ent.WuGeLucky) bool
	checkSkipWuXingFilter      func(ge int, ge2 int, ge3 int) bool
	checkSkipStrokeNumberScope func(stroke []int) bool
	queryCharacterFilter       func(query *ent.CharacterQuery) *ent.CharacterQuery
	queryRegularFilter         func(query *ent.CharacterQuery) *ent.CharacterQuery
	queryStrokeFilter          func(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery
}

func (f *filter) FilterType() CharacterFilterType {
	return f.characterFilterType
}

func (f *filter) CheckSkipStrokeNumberScope(stroke ...int) bool {
	return f.checkSkipStrokeNumberScope(stroke)
}

func (f *filter) QueryStrokeFilter(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery {
	return f.queryStrokeFilter(stroke)
}

func (f *filter) QueryRegularFilter(query *ent.CharacterQuery) *ent.CharacterQuery {
	return f.queryRegularFilter(query)
}

func (f *filter) GetCharacterStroke(c *ent.Character) int {
	switch f.characterFilterType {
	case CharacterFilterTypeChs:
		return c.SimpleTotalStroke
	case CharacterFilterTypeCht:
		return c.TraditionalTotalStroke
	case CharacterFilterTypeKangxi:
		return c.KangXiStroke
	case CharacterFilterTypeDefault:
		fallthrough
	default:
		return c.ScienceStroke
	}
}

func (f *filter) CustomFilter(name string, v any) error {
	if v == nil {
		return errors.New("not implements")
	}
	switch name {
	//case "CheckSkipSexFilter":
	//	fn, ok := v.(func(lucky *ent.WuGeLucky) bool)
	//	if ok {
	//		f.checkSkipSexFilter = fn
	//	}
	//case "CheckSkipWuXingFilter":
	//	fn, ok := v.(func(ge int, ge2 int, ge3 int) bool)
	//	if ok {
	//		f.checkSkipWuXingFilter = fn
	//	}
	//case "CheckSkipDaYanFilter":
	//case "QueryCharacterFilter":
	default:
		return errors.New("not implements")
	}
}

func (f *filter) CheckSkipWuXingFilter(ge int, ge2 int, ge3 int) bool {
	return f.checkSkipWuXingFilter(ge, ge2, ge3)
}

func (f *filter) CheckSkipDaYanFilter(lucky *ent.WuGeLucky) bool {
	return f.checkSkipDaYanFilter(lucky)
}

func (f *filter) CheckSkipSexFilter(lucky *ent.WuGeLucky) bool {
	return f.checkSkipSexFilter(lucky)
}

func (f *filter) QueryCharacterFilter(query *ent.CharacterQuery) *ent.CharacterQuery {
	return f.queryCharacterFilter(query)
}

// DefaultFilter ...
// @return Filter
func DefaultFilter() Filter {
	return defaultFilter
}

func newFilter() *filter {
	return &filter{
		characterFilterType: CharacterFilterTypeDefault,
		checkSkipSexFilter: func(lucky *ent.WuGeLucky) bool {
			return false
		},
		checkSkipDaYanFilter: func(lucky *ent.WuGeLucky) bool {
			return false
		},
		checkSkipWuXingFilter: func(ge int, ge2 int, ge3 int) bool {
			return false
		},
		checkSkipStrokeNumberScope: func(stroke []int) bool {
			return false
		},
		queryCharacterFilter: func(query *ent.CharacterQuery) *ent.CharacterQuery {
			return query
		},
		queryRegularFilter: func(query *ent.CharacterQuery) *ent.CharacterQuery {
			return query
		},
		queryStrokeFilter: func(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery {
			return func(query *ent.CharacterQuery) *ent.CharacterQuery {
				return query.Where(character.StrokeEQ(stroke))
			}
		},
	}
}

// NewFilter creates a new Filter based on the given FilterOption.
//
// Parameters:
// - fo: The FilterOption used to configure the Filter.
//
// Returns:
// - Filter: The newly created Filter.
func NewFilter(fo FilterOption) Filter {
	f := newFilter()
	if fo.SexFilter {
		f.checkSkipSexFilter = func(lucky *ent.WuGeLucky) bool {
			return lucky.ZongLucky == false
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

	return f
}

func characterTypeFilterDefault(min, max int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.ScienceStrokeGTE(min)).Where(character.And(character.ScienceStrokeLTE(max)))
	}
}

func strokeFilterDefault(s int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.ScienceStrokeEQ(s))
	}
}

func characterTypeFilterCHS(min, max int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.SimpleTotalStrokeGTE(min)).Where(character.And(character.SimpleTotalStrokeLTE(max)))
	}
}

func strokeFilterCHS(s int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.SimpleTotalStrokeEQ(s))
	}
}

func characterTypeFilterCHT(min, max int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.TraditionalTotalStrokeGTE(min)).Where(character.And(character.TraditionalTotalStrokeLTE(max)))
	}
}

func strokeFilterCHT(s int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.TraditionalTotalStrokeEQ(s))
	}
}

func characterTypeFilterKX(min, max int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.KangXiStrokeGTE(min)).Where(character.And(character.KangXiStrokeLTE(max)))
	}
}

func strokeFilterKX(s int) func(*ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.KangXiStrokeEQ(s))
	}
}

func regularFilter(query *ent.CharacterQuery) *ent.CharacterQuery {
	return query.Where(character.RegularEQ(true))
}

func wuXingFilter(ge int, ge2 int, ge3 int) bool {
	sc := core.NewSanCai(ge, ge2, ge3)
	return sc.Check(5) == false
}

func daYanFilter(lucky *ent.WuGeLucky) bool {
	// simplified logic
	return isLucky(dayan.Find(lucky.DiGe).Lucky) &&
		isLucky(dayan.Find(lucky.RenGe).Lucky) &&
		isLucky(dayan.Find(lucky.WaiGe).Lucky) &&
		isLucky(dayan.Find(lucky.ZongGe).Lucky)
}

func isLucky(s string) bool {
	return strings.Index(s, "吉") != -1
}
