package wuxing

const sanCai = "水木木火火土土金金水"
const yinYang = "阴阳"

type SanCai struct {
	tianCai        string `bson:"tian_cai"`
	tianCaiYinYang string `bson:"tian_cai_yin_yang"`
	renCai         string `bson:"ren_cai"`
	renCaiYinYang  string `bson:"ren_cai_yin_yang"`
	diCai          string `bson:"di_cai"`
	diCaiYingYang  string `bson:"di_cai_ying_yang"`
	fortune        string `bson:"fortune"`
	comment        string `bson:"comment"`
}

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

func (s *SanCai) Check(point int) bool {
	if GetLuckyPoint(s.String()) >= point {
		return true
	}
	return false
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
