package fate

import "github.com/xormsharp/xorm"

const sanCai = "水木木火火土土金金水"
const yinYang = "阴阳"

// SanCai ...
type SanCai struct {
	tianCai        string `bson:"tian_cai"`
	tianCaiYinYang string `bson:"tian_cai_yin_yang"`
	renCai         string `bson:"ren_cai"`
	renCaiYinYang  string `bson:"ren_cai_yin_yang"`
	diCai          string `bson:"di_cai"`
	diCaiYingYang  string `bson:"di_cai_ying_yang"`
	fortune        string `bson:"fortune"` //吉凶
	comment        string `bson:"comment"` //说明
}

//NewSanCai 新建一个三才对象
func NewSanCai(tian, ren, di int) *SanCai {
	return &SanCai{
		tianCai:        sanCaiAttr(tian),
		tianCaiYinYang: yinYangAttr(tian),
		renCai:         sanCaiAttr(ren),
		renCaiYinYang:  yinYangAttr(ren),
		diCai:          sanCaiAttr(di),
		diCaiYingYang:  yinYangAttr(di),
	}
}

//Check 检查三才属性
func Check(engine *xorm.Engine, cai *SanCai, point int) bool {
	wx := FindWuXing(engine, cai.tianCai, cai.renCai, cai.diCai)
	if wx.Luck.Point() >= point {
		return true
	}
	return false
}

// GenerateThreeTalent 计算字符的三才属性
// 1-2木：1为阳木，2为阴木   3-4火：3为阳火，4为阴火   5-6土：5为阳土，6为阴土   7-8金：7为阳金，8为阴金   9-10水：9为阳水，10为阴水
func sanCaiAttr(i int) string {
	return string([]rune(sanCai)[i%10])
}

func yinYangAttr(i int) string {
	return string([]rune(yinYang)[i%2])
}
