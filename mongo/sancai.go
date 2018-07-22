//三才五格
package mongo

import "gopkg.in/mgo.v2/bson"


type SanCai struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	TianCai        string        `bson:"tian_cai"`
	TianCaiYinYang string        `bson:"tian_cai_yin_yang"`
	RenCai         string        `bson:"ren_cai"`
	RenCaiYinYang  string        `bson:"ren_cai_yin_yang"`
	DiCai          string        `bson:"di_cai"`
	DiCaiYingYang  string        `bson:"di_cai_ying_yang"`
	Fortune        string        `bson:"fortune"` //吉凶
	Comment        string        `bson:"comment"` //说明
}


func MakeSanCai(wuGe *WuGe) *SanCai {
	return &SanCai{
		TianCai:        sanCaiAttr(wuGe.TianGe),
		TianCaiYinYang: sanCaiYinYang(wuGe.TianGe),
		RenCai:         sanCaiAttr(wuGe.RenGe),
		RenCaiYinYang:  sanCaiYinYang(wuGe.RenGe),
		DiCai:          sanCaiAttr(wuGe.DiGe),
		DiCaiYingYang:  sanCaiYinYang(wuGe.DiGe),
	}
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
