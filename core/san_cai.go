package core

var sanCai = []rune("水木木火火土土金金水")

var yinYang = []rune("阴阳")

// SanCai 三才
type SanCai struct {
	tianCai        string
	tianCaiYinYang string
	renCai         string
	renCaiYinYang  string
	diCai          string
	diCaiYingYang  string
	fortune        string
	comment        string
}

// NewSanCai 创建三才
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

// Check 检查三才属性
func (s *SanCai) Check(point int) bool {
	wx, exist := NewWuXing(s.String())
	if !exist {
		return false
	}
	if wx.Point() >= point {
		return true
	}
	return false
}

func (s *SanCai) String() string {
	return s.tianCai + s.renCai + s.diCai
}

// GenerateThreeTalent 计算字符的三才属性
// 1-2木：1为阳木，2为阴木   3-4火：3为阳火，4为阴火   5-6土：5为阳土，6为阴土   7-8金：7为阳金，8为阴金   9-10水：9为阳水，10为阴水
func sanCaiAttr(i int) string {
	return string([]rune(sanCai)[i%10])
}

// yinYangAttr returns a string that represents the Yin or Yang attribute based on the given integer.
//
// The function takes an integer i as a parameter and returns a string. The integer is used to determine whether to return the Yin or Yang attribute from the yinYang string. The Yin or Yang attribute is selected by using the modulo operator on i with 2. If i is an even number, the Yin attribute is returned. If i is an odd number, the Yang attribute is returned.
func yinYangAttr(i int) string {
	return string([]rune(yinYang)[i%2])
}
