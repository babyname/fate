package fate

import (
	"github.com/godcong/fate/mongo"
)

//WuGe
type WuGe struct {
	TianGe int `bson:"tian_ge"`
	RenGe  int `bson:"ren_ge"`
	DiGe   int `bson:"di_ge"`
	WaiGe  int `bson:"wai_ge"`
	ZongGe int `bson:"zong_ge"`
}

//NewWuGe 计算五格
func NewWuGe(l1, l2, f1, f2 int) *WuGe {
	return &WuGe{
		TianGe: tianGe(l1, l2, f1, f2),
		RenGe:  renGe(l1, l2, f1, f2),
		DiGe:   diGe(l1, l2, f1, f2),
		WaiGe:  waiGe(l1, l2, f1, f2),
		ZongGe: zongGe(l1, l2, f1, f2),
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
	} else {
		return l1 + f1
	}
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

//zongGe input the ScienceStrokes with name
//总格，姓加名的笔画总数  数理五行分类
func zongGe(l1, l2, f1, f2 int) int {
	return l1 + l2 + f1 + f2
}

func check(dy []*mongo.DaYan, idx int) bool {
	switch dy[idx-1].Fortune {
	case "吉", "半吉":
		return true
	}
	return false
}

//Check 格检查
func (ge *WuGe) Check() bool {
	dy := mongo.GetDaYan()
	if dy == nil {
		return false
	}

	for _, v := range []int{ge.DiGe, ge.RenGe, ge.WaiGe, ge.ZongGe} {
		if len(dy) < v || !check(dy, v) {
			return false
		}
	}
	return true
}
