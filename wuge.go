package fate

import (
	"fmt"

	"github.com/godcong/yi"
)

// WuGe ...
type WuGe struct {
	tianGe int
	renGe  int
	diGe   int
	waiGe  int
	zongGe int
}

//五格笔画缓存
var wuGes map[int][]*WuGe = map[int][]*WuGe{}

func ClearWuGe(nk NameStroke) {
	hash := nk.hash()

	if wuGes[hash] != nil {
		delete(wuGes, hash)
	}
}

//从笔画获取五格
func GetWuGe(nk NameStroke, filter_hard bool) []*WuGe {
	hash := nk.hash()

	wuGe := wuGes[hash]
	if wuGe == nil {
		wuGes[hash] = []*WuGe{calcWuGe(nk.LastStroke1, nk.LastStroke2, nk.FirstStroke1, nk.FirstStroke2)}
		if filter_hard && nk.LastStroke2 != 0 && nk.FirstStroke2 != 0 {
			// 复名复姓算法2
			wuGes[hash] = append(wuGes[hash], &WuGe{
				tianGe: tianGe(nk.LastStroke1, nk.LastStroke2, nk.FirstStroke1, nk.FirstStroke2),
				renGe:  renGe2(nk.LastStroke1, nk.LastStroke2, nk.FirstStroke1, nk.FirstStroke2),
				diGe:   diGe(nk.LastStroke1, nk.LastStroke2, nk.FirstStroke1, nk.FirstStroke2),
				waiGe:  waiGe2(nk.LastStroke1, nk.LastStroke2, nk.FirstStroke1, nk.FirstStroke2),
				zongGe: zongGe(nk.LastStroke1, nk.LastStroke2, nk.FirstStroke1, nk.FirstStroke2),
			})
		}
		return wuGes[hash]
	}

	return wuGe
}

func (wg *WuGe) getSanCai() SanCai {
	return GetSanCai(wg.tianGe, wg.renGe, wg.diGe)
}

// ZongGe ...
func (ge *WuGe) ZongGe() int {
	return ge.zongGe
}

// WaiGe ...
func (ge *WuGe) WaiGe() int {
	return ge.waiGe
}

// DiGe ...
func (ge *WuGe) DiGe() int {
	return ge.diGe
}

// RenGe ...
func (ge *WuGe) RenGe() int {
	return ge.renGe
}

// TianGe ...
func (ge *WuGe) TianGe() int {
	return ge.tianGe
}

//calcWuGe 计算五格
func calcWuGe(l1, l2, f1, f2 int) *WuGe {
	return &WuGe{
		tianGe: tianGe(l1, l2, f1, f2),
		renGe:  renGe(l1, l2, f1, f2),
		diGe:   diGe(l1, l2, f1, f2),
		waiGe:  waiGe(l1, l2, f1, f2),
		zongGe: zongGe(l1, l2, f1, f2),
	}
}

//tianGe input the ScienceStrokes with last name
//天格（复姓）姓的笔画相加
//天格（单姓）姓的笔画上加一
func tianGe(l1, l2, _, _ int) int {
	if l2 == 0 {
		return l1 + 1
	}
	return l1 + l2
}

//renGe input the ScienceStrokes with name
//人格（复姓）姓氏的第二字的笔画加名的第一字
//人格（复姓单名）姓的第二字加名
//人格（单姓单名）姓加名
// 人格（单姓复名）姓加名的第一字
func renGe(l1, l2, f1, _ int) int {
	//人格（复姓）姓氏的第二字的笔画加名的第一字
	//人格（复姓单名）姓的第二字加名
	if l2 != 0 {
		return l2 + f1
	}
	return l1 + f1
}

//复名复姓人格算法2
func renGe2(l1, l2, f1, _ int) int {
	return l1 + l2 + f1
}

//diGe input the ScienceStrokes with name
//地格（复姓复名，单姓复名）名字相加
//地格（复姓单名，单姓单名）名字+1
func diGe(_, _, f1, f2 int) int {
	if f2 == 0 {
		return f1 + 1
	}
	return f1 + f2
}

//waiGe input the ScienceStrokes with name
//外格（复姓单名）姓的第一字加笔画数一
//外格（复姓复名）姓的第一字和名的最后一定相加的笔画数
//外格（单姓复名）一加名的最后一个字
//外格（单姓单名）一加一
func waiGe(l1, l2, _, f2 int) (n int) {
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

//复名复姓外格算法2
func waiGe2(_, _, _, f2 int) (n int) {
	return 1 + f2
}

//zongGe input the ScienceStrokes with name
//总格，姓加名的笔画总数  数理五行分类
func zongGe(l1, l2, f1, f2 int) int {
	//归1
	zg := (l1 + l2 + f1 + f2) - 1
	if zg < 1 {
		fmt.Println(l1, l2, f1, f2)
		panic("bad input")
	}
	return zg%81 + 1
}

//Check 格检查
func (ge *WuGe) Check(ss ...string) bool {
	v := map[string]bool{}
	if ss == nil {
		ss = append(ss, "吉")
	}
	//ignore:tianGe
	v[yi.GetDaYan(ge.diGe).Lucky] = false
	v[yi.GetDaYan(ge.renGe).Lucky] = false
	v[yi.GetDaYan(ge.waiGe).Lucky] = false
	v[yi.GetDaYan(ge.zongGe).Lucky] = false

	for l := range v {
		for i := range ss {
			if ss[i] == l {
				v[l] = true
				break
			}
		}
	}
	for l := range v {
		if !v[l] {
			return false
		}
	}
	return true
}
