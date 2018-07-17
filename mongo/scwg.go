//三才五格
package mongo

type WuGe struct {
	TianGe int `json:"tian_ge"`
	RenGe  int `json:"ren_ge"`
	DiGe   int `json:"di_ge"`
	WaiGe  int `json:"wai_ge"`
	ZongGe int `json:"zong_ge"`
}

type SanCai struct {
	TianCai        string `json:"tian_cai"`
	TianCaiYinYang string `json:"tian_cai_yin_yang"`
	RenCai         string `json:"ren_cai"`
	RenCaiYinYang  string `json:"ren_cai_yin_yang"`
	DiCai          string `json:"di_cai"`
	DiCaiYingYang  string `json:"di_cai_ying_yang"`
}

type SanCaiWuGe struct {
	*WuGe
	*SanCai
	Fortune string `json:"fortune"` //吉凶
	Comment string `json:"comment"` //说明
}

func MakeSanCaiWuGe(l1, l2, f1, f2 int) *SanCaiWuGe {
	wuGe := WuGe{
		TianGe: tianGe(l1, l2, f1, f2),
		RenGe:  renGe(l1, l2, f1, f2),
		DiGe:   diGe(l1, l2, f1, f2),
		WaiGe:  waiGe(l1, l2, f1, f2),
		ZongGe: zongGe(l1, l2, f1, f2),
	}

	sanCai := SanCai{
		TianCai:        sanCaiAttr(wuGe.TianGe),
		TianCaiYinYang: sanCaiYinYang(wuGe.TianGe),
		RenCai:         sanCaiAttr(wuGe.RenGe),
		RenCaiYinYang:  sanCaiYinYang(wuGe.RenGe),
		DiCai:          sanCaiAttr(wuGe.DiGe),
		DiCaiYingYang:  sanCaiYinYang(wuGe.DiGe),
	}

	return &SanCaiWuGe{
		WuGe:   &wuGe,
		SanCai: &sanCai,
	}
}

//tianGe input the ScienceStrokes with last name
//天格（复姓）姓的笔画相加
//天格（单姓）姓的笔画上加一
func tianGe(l1, l2, f1, f2 int) int {
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
func renGe(l1, l2, f1, f2 int) int {
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
func diGe(l1, l2, f1, f2 int) int {
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
func waiGe(l1, l2, f1, f2 int) (n int) {
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

// GenerateThreeTalent 计算字符的三才属性
// 1-2木：1为阳木，2为阴木   3-4火：3为阳火，4为阴火   5-6土：5为阳土，6为阴土   7-8金：7为阳金，8为阴金   9-10水：9为阳水，10为阴水
func sanCaiAttr(i int) string {
	var attr string
	switch i % 10 {
	case 1, 2:
		attr = "木"
	case 3, 4:
		attr = "火"
	case 5, 6:
		attr = "土"
	case 7, 8:
		attr = "金"
	case 9, 0:
		fallthrough
	default:
		attr = "水"
	}

	return attr
}

func sanCaiYinYang(i int) string {
	var yy string
	if i%2 == 1 {
		yy = "阳"
	} else {
		yy = "阴"
	}
	return yy
}
