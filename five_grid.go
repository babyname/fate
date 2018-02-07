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
func MakeFiveGridFromStrokes(l1, l2, f1, f2 int) FiveGrid {
	return FiveGrid{
		SkyGrid:    skyGrid(l1, l2, f1, f2),
		LandGrid:   landGrid(l1, l2, f1, f2),
		PersonGrid: personGrid(l1, l2, f1, f2),
		OutGrid:    outGrid(l1, l2, f1, f2),
		AllGrid:    allGrid(l1, l2, f1, f2),
	}

}

//skyGrid input the ScienceStrokes with last name
//天格（复姓）姓的笔画相加
//天格（单姓）姓的笔画上加一
func skyGrid(l1, l2, f1, f2 int) int {
	if l2 == 0 {
		return l1 + 1
	}
	return l1 + l2
}

//landGrid input the ScienceStrokes with name
//人格（复姓）姓氏的第二字的笔画加名的第一字
//人格（复姓单名）姓的第二字加名
//人格（单姓单名）姓加名
// 人格（单姓复名）姓加名的第一字
func personGrid(l1, l2, f1, f2 int) int {
	//人格（复姓）姓氏的第二字的笔画加名的第一字
	//人格（复姓单名）姓的第二字加名
	if l2 != 0 {
		return l2 + f1
	} else {
		return l1 + f1
	}
}

//personGrid input the ScienceStrokes with name
//地格（复姓复名，单姓复名）名字相加
//地格（复姓单名，单姓单名）名字+1
func landGrid(l1, l2, f1, f2 int) int {
	if f2 == 0 {
		return f1 + 1
	}
	return f1 + f2
}

//personGrid input the ScienceStrokes with name
//外格（复姓单名）姓的第一字加笔画数一
//外格（复姓复名）姓的第一字和名的最后一定相加的笔画数
//外格（单姓复名）一加名的最后一个字
//外格（单姓单名）一加一
func outGrid(l1, l2, f1, f2 int) (n int) {
	//单姓单名
	if l2 == 0 && f2 == 0 {
		n = 1 + 1
	}
	//单姓复名
	if l2 == 0 && f2 != 0 {
		n = 1 + f2
	}
	//复姓单名
	if l2 != 0 && f2 == 0 {
		n = l1 + 1
	}
	//复姓复名
	if l2 != 0 && f2 != 0 {
		n = l1 + f2
	}
	return n
}

//allGrid input the ScienceStrokes with name
//总格，姓加名的笔画总数  数理五行分类
func allGrid(l1, l2, f1, f2 int) int {
	return l1 + l2 + f1 + f2
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
