package fate

import (
	"errors"
	"strings"

	"github.com/babyname/fate/dayan"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/wuxing"
)

var defaultFilter Filter

// Filter
// prefix with Query are used db search
type Filter interface {
	CustomFilter(name string, v any) error
	QueryCharacterFilter(query *ent.CharacterQuery) *ent.CharacterQuery
	QueryRegularFilter(query *ent.CharacterQuery) *ent.CharacterQuery
	QueryStrokeFilter(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery
	CheckSexFilter(lucky *ent.WuGeLucky) bool
	CheckDaYanFilter(lucky *ent.WuGeLucky) bool
	CheckWuXingFilter(ge int, ge2 int, ge3 int) bool
	CheckStrokeNumber(stroke int) bool
	GetCharacterStroke(c *ent.Character) int
}

type filter struct {
	characterFilterType  CharacterFilterType
	checkSexFilter       func(lucky *ent.WuGeLucky) bool
	checkDaYanFilter     func(lucky *ent.WuGeLucky) bool
	checkWuXingFilter    func(ge int, ge2 int, ge3 int) bool
	checkStrokeNumber    func(stroke int) bool
	queryCharacterFilter func(query *ent.CharacterQuery) *ent.CharacterQuery
	queryRegularFilter   func(query *ent.CharacterQuery) *ent.CharacterQuery
	queryStrokeFilter    func(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery
}

func (f *filter) CheckStrokeNumber(stroke int) bool {
	return f.checkStrokeNumber(stroke)
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
	//case "CheckSexFilter":
	//	fn, ok := v.(func(lucky *ent.WuGeLucky) bool)
	//	if ok {
	//		f.checkSexFilter = fn
	//	}
	//case "CheckWuXingFilter":
	//	fn, ok := v.(func(ge int, ge2 int, ge3 int) bool)
	//	if ok {
	//		f.checkWuXingFilter = fn
	//	}
	//case "CheckDaYanFilter":
	//case "QueryCharacterFilter":
	default:
		return errors.New("not implements")
	}
}

func (f *filter) CheckWuXingFilter(ge int, ge2 int, ge3 int) bool {
	return f.checkWuXingFilter(ge, ge2, ge3)
}

func (f *filter) CheckDaYanFilter(lucky *ent.WuGeLucky) bool {
	return f.checkDaYanFilter(lucky)
}

func (f *filter) CheckSexFilter(lucky *ent.WuGeLucky) bool {
	return f.checkSexFilter(lucky)
}

func (f *filter) QueryCharacterFilter(query *ent.CharacterQuery) *ent.CharacterQuery {
	return f.queryCharacterFilter(query)
}

// DefaultFilter ...
// @return Filter
func DefaultFilter() Filter {
	if defaultFilter == nil {
		defaultFilter = newFilter()
	}
	return defaultFilter
}

func newFilter() *filter {
	return &filter{
		characterFilterType: CharacterFilterTypeDefault,
		checkSexFilter: func(lucky *ent.WuGeLucky) bool {
			return false
		},
		checkDaYanFilter: func(lucky *ent.WuGeLucky) bool {
			return false
		},
		checkWuXingFilter: func(ge int, ge2 int, ge3 int) bool {
			return false
		},
		checkStrokeNumber: func(stroke int) bool {
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

func NewFilter(fo FilterOption) Filter {
	f := newFilter()
	if fo.SexFilter {
		f.checkSexFilter = func(lucky *ent.WuGeLucky) bool {
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

	f.checkStrokeNumber = func(stroke int) bool {
		if fo.MinStroke != 0 && stroke < fo.MinStroke {
			return true
		}
		if fo.MaxStroke != 0 && stroke > fo.MaxStroke {
			return true
		}
		return false
	}

	if fo.DaYanFilter {
		f.checkDaYanFilter = daYanFilter
	}
	if fo.WuXingFilter {
		f.checkWuXingFilter = wuXingFilter
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
	sc := wuxing.NewSanCai(ge, ge2, ge3)
	return sc.Check(5) == false
}

func daYanFilter(lucky *ent.WuGeLucky) bool {
	if isLucky(dayan.Find(lucky.DiGe).Lucky) &&
		isLucky(dayan.Find(lucky.RenGe).Lucky) &&
		isLucky(dayan.Find(lucky.WaiGe).Lucky) &&
		isLucky(dayan.Find(lucky.ZongGe).Lucky) {
		return false
	}
	log.Info("dayan filter", "lucky", lucky,
		"dige", isLucky(dayan.Find(lucky.DiGe).Lucky),
		"renge", isLucky(dayan.Find(lucky.RenGe).Lucky),
		"waige", isLucky(dayan.Find(lucky.WaiGe).Lucky),
		"zongge", isLucky(dayan.Find(lucky.ZongGe).Lucky))
	return true
}

func isLucky(s string) bool {
	if strings.Index(s, "吉") != -1 {
		return true
	}
	return false
}
