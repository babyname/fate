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
	DaYanFilter(lucky *ent.WuGeLucky) bool
	WuXingFilter(ge int, ge2 int, ge3 int) bool
}

type filter struct {
	sexFilter       func(lucky *ent.WuGeLucky) bool
	daYanFilter     func(lucky *ent.WuGeLucky) bool
	wuXingFilter    func(ge int, ge2 int, ge3 int) bool
	characterFilter func(query *ent.CharacterQuery) *ent.CharacterQuery
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
	return &filter{
		sexFilter: func(lucky *ent.WuGeLucky) bool {
			return true
		},
		daYanFilter: func(lucky *ent.WuGeLucky) bool {
			return true
		},
		wuXingFilter: func(ge int, ge2 int, ge3 int) bool {
			return true
		},
		characterFilter: func(query *ent.CharacterQuery) *ent.CharacterQuery {
			return query
		},
	}
}

func NewFilter(fo FilterOption) Filter {
	var f filter
	if fo.SexFilter {
		f.sexFilter = func(lucky *ent.WuGeLucky) bool {
			return lucky.ZongLucky == false
		}
	}
	if fo.CharacterFilter {
		f.characterFilter = func(query *ent.CharacterQuery) *ent.CharacterQuery {
			return query.Where(character.StrokeGTE(fo.MinCharacter)).Where(character.And(character.StrokeLTE(fo.MaxCharacter)))
		}
	}
	if fo.DaYanFilter {
		f.daYanFilter = daYanFilter
	}
	if fo.DaYanFilter && fo.WuXingFilter {
		f.wuXingFilter = wuXingFilter
	}

	return &f
}

func wuXingFilter(ge int, ge2 int, ge3 int) bool {
	sc := wuxing.NewSanCai(ge, ge2, ge3)
	return sc.Check(5) == false
}

func daYanFilter(lucky *ent.WuGeLucky) bool {
	if !isLucky(dayan.Find(lucky.DiGe).Lucky) ||
		!isLucky(dayan.Find(lucky.RenGe).Lucky) ||
		!isLucky(dayan.Find(lucky.WaiGe).Lucky) ||
		!isLucky(dayan.Find(lucky.ZongGe).Lucky) {
		return false
	}
	return true
}

func isLucky(s string) bool {
	if strings.Compare(s, "吉") == 0 || strings.Compare(s, "半吉") == 0 {
		return true
	}
	return false
}
