package fate

import (
	"fmt"

	"github.com/godcong/fate/debug"
	"github.com/godcong/fate/model"
)

type FiveGrid struct {
	SkyGrid    int //天格
	LandGrid   int //地格
	PersonGrid int //人格
	OutGrid    int //外格
	AllGrid    int //总格
}

// MakeFiveGridFromStrokes 按照姓名笔画创建一个五格属性
// 五格是进行大衍计算的基础
func MakeFiveGridFromStrokes(l1, l2, f1, f2 int) *FiveGrid {
	return &FiveGrid{
		SkyGrid:    skyGrid(l1, l2, f1, f2),
		LandGrid:   landGrid(l1, l2, f1, f2),
		PersonGrid: personGrid(l1, l2, f1, f2),
		OutGrid:    outGrid(l1, l2, f1, f2),
		AllGrid:    allGrid(l1, l2, f1, f2),
	}
}

func (fg *FiveGrid) PrintBigYan(filter bool) bool {
	sg := model.GetBigYanByIndex(fg.SkyGrid)
	lg := model.GetBigYanByIndex(fg.LandGrid)
	pg := model.GetBigYanByIndex(fg.PersonGrid)
	og := model.GetBigYanByIndex(fg.OutGrid)
	ag := model.GetBigYanByIndex(fg.AllGrid)
	v := fmt.Sprintf("天格(%d):%s,地格(%d):%s,人格(%d):%s,外格(%d):%s,总格(%d):%s",
		sg.Index, sg.Goil,
		lg.Index, lg.Goil,
		pg.Index, pg.Goil,
		og.Index, og.Goil,
		ag.Index, ag.Goil,
	)

	if filter {
		if (lg.Goil == "吉" || lg.Goil == "半吉") &&
			(pg.Goil == "吉" || pg.Goil == "半吉") &&
			(og.Goil == "吉" || og.Goil == "半吉") &&
			(ag.Goil == "吉" || ag.Goil == "半吉") {
			debug.Println(v)
			return true
		}

	} else {
		debug.Println(v)
		return true
	}
	return false
}

func (fg *FiveGrid) ContainBest(rule ...string) bool {
	b := true

	validate := []model.BigYan{
		model.GetBigYanByIndex(fg.LandGrid),
		model.GetBigYanByIndex(fg.PersonGrid),
		model.GetBigYanByIndex(fg.OutGrid),
		model.GetBigYanByIndex(fg.AllGrid),
	}

	switch len(rule) {
	case 0:
	case 1:
		for _, v := range validate {
			if v.Goil != rule[0] {
				b = false
			}
		}
	case 2:
		for _, v := range validate {
			if v.Goil != rule[0] && v.Goil != rule[1] {
				b = false
			}
		}
	case 3:
		for _, v := range validate {
			if v.Goil != rule[0] && v.Goil != rule[1] && v.Goil != rule[2] {
				b = false
			}
		}
	default:
		b = false
	}
	return b
}
