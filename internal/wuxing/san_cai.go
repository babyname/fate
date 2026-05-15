// Package wuxing 实现五行与三才的属性计算及吉凶判断。
package wuxing

const sanCai = "水木木火火土土金金水"
const yinYang = "阴阳"

// SanCai 表示三才（天、人、地）的五行属性及阴阳信息。
type SanCai struct {
	tianCai        string `bson:"tian_cai"`
	tianCaiYinYang string `bson:"tian_cai_yin_yang"`
	renCai         string `bson:"ren_cai"`
	renCaiYinYang  string `bson:"ren_cai_yin_yang"`
	diCai          string `bson:"di_cai"`
	diCaiYingYang  string `bson:"di_cai_ying_yang"`
}

// NewSanCai 根据天格、人格、地格数理创建三才实例。
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

// Check 判断三才吉凶分数是否达到指定阈值。
func (s *SanCai) Check(point int) bool {
	return GetLuckyPoint(s.String()) >= point
}

func (s *SanCai) String() string {
	return s.tianCai + s.renCai + s.diCai
}

func sanCaiAttr(i int) string {
	return string([]rune(sanCai)[i%10])
}

func yinYangAttr(i int) string {
	return string([]rune(yinYang)[i%2])
}
