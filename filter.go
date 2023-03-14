package fate

import (
	"strings"

	"github.com/babyname/fate/dayan"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/wuxing"
)

type Filter interface {
	CustomFilter(name string, v any)
	SexFilter(lucky *ent.WuGeLucky) bool
	CharacterFilter(query *ent.CharacterQuery) *ent.CharacterQuery
	RegularFilter(query *ent.CharacterQuery) *ent.CharacterQuery
	StrokeFilter(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery
	DaYanFilter(lucky *ent.WuGeLucky) bool
	WuXingFilter(ge int, ge2 int, ge3 int) bool
	GetCharacterStroke(e *ent.Character) int
}

type filter struct {
	characterFilterType CharacterFilterType
	sexFilter           func(lucky *ent.WuGeLucky) bool
	daYanFilter         func(lucky *ent.WuGeLucky) bool
	wuXingFilter        func(ge int, ge2 int, ge3 int) bool
	characterFilter     func(query *ent.CharacterQuery) *ent.CharacterQuery
	regularFilter       func(query *ent.CharacterQuery) *ent.CharacterQuery
	strokeFilter        func(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery
}

func (f *filter) StrokeFilter(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery {
	return f.strokeFilter(stroke)
}

func (f *filter) RegularFilter(query *ent.CharacterQuery) *ent.CharacterQuery {
	return f.regularFilter(query)
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

func (f *filter) CustomFilter(name string, v any) {
	if v == nil {
		return
	}
	switch name {
	case "SexFilter":
		fn, ok := v.(func(lucky *ent.WuGeLucky) bool)
		if ok {
			f.sexFilter = fn
		}
	case "WuXingFilter":
		fn, ok := v.(func(ge int, ge2 int, ge3 int) bool)
		if ok {
			f.wuXingFilter = fn
		}
	case "DaYanFilter":
	case "CharacterFilter":
	default:
		return
	}
}

func (f *filter) WuXingFilter(ge int, ge2 int, ge3 int) bool {
	return f.wuXingFilter(ge, ge2, ge3)
}

func (f *filter) DaYanFilter(lucky *ent.WuGeLucky) bool {
	return f.daYanFilter(lucky)
}

func (f *filter) SexFilter(lucky *ent.WuGeLucky) bool {
	return f.sexFilter(lucky)
}

func (f *filter) CharacterFilter(query *ent.CharacterQuery) *ent.CharacterQuery {
	return f.characterFilter(query)
}

// DefaultFilter ...
// @return Filter
func DefaultFilter() Filter {
	return newFilter()
}

func newFilter() *filter {
	return &filter{
		characterFilterType: CharacterFilterTypeDefault,
		sexFilter: func(lucky *ent.WuGeLucky) bool {
			return false
		},
		daYanFilter: func(lucky *ent.WuGeLucky) bool {
			return false
		},
		wuXingFilter: func(ge int, ge2 int, ge3 int) bool {
			return false
		},
		characterFilter: func(query *ent.CharacterQuery) *ent.CharacterQuery {
			return query
		},
		regularFilter: func(query *ent.CharacterQuery) *ent.CharacterQuery {
			return query
		},
		strokeFilter: func(stroke int) func(query *ent.CharacterQuery) *ent.CharacterQuery {
			return func(query *ent.CharacterQuery) *ent.CharacterQuery {
				return query.Where(character.StrokeEQ(stroke))
			}
		},
	}
}

func NewFilter(fo FilterOption) Filter {
	f := newFilter()
	if fo.SexFilter {
		f.sexFilter = func(lucky *ent.WuGeLucky) bool {
			return lucky.ZongLucky == false
		}
	}
	if fo.CharacterFilter {
		f.characterFilterType = fo.CharacterFilterType
		switch fo.CharacterFilterType {
		case CharacterFilterTypeChs:
			f.characterFilter = characterTypeFilterCHS(fo.MinCharacter, fo.MaxCharacter)
			f.strokeFilter = strokeFilterCHS
		case CharacterFilterTypeCht:
			f.characterFilter = characterTypeFilterCHT(fo.MinCharacter, fo.MaxCharacter)
			f.strokeFilter = strokeFilterCHT
		case CharacterFilterTypeKangxi:
			f.characterFilter = characterTypeFilterKX(fo.MinCharacter, fo.MaxCharacter)
			f.strokeFilter = strokeFilterKX
		case CharacterFilterTypeDefault:
			fallthrough
		default:
			f.characterFilter = characterTypeFilterDefault(fo.MinCharacter, fo.MaxCharacter)
			f.strokeFilter = strokeFilterDefault
		}
	}
	if fo.DaYanFilter {
		f.daYanFilter = daYanFilter
	}
	if fo.WuXingFilter {
		f.wuXingFilter = wuXingFilter
	}
	if fo.RegularFilter {
		f.regularFilter = regularFilter
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
