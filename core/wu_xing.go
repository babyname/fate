package core

// WuXing 五行
type WuXing struct {
	wuXing string
	point  int
}

// String returns the string representation of the WuXing object.
//
// No parameters.
// Returns a string.
func (wx WuXing) String() string {
	return string(wx.wuXing)
}

// Point returns the lucky point of the WuXing object.
//
// No parameters.
// Returns an integer.
func (wx WuXing) Point() int {
	return wx.point
}

var wuXingTables = map[string]string{
	"木木木": "大吉",
	"木木火": "大吉",
	"木木土": "大吉",
	"木木金": "凶多吉少",
	"木木水": "吉多凶少",
	"木火木": "大吉",
	"木火火": "吉",
	"木火土": "大吉",
	"木火金": "凶多吉少",
	"木火水": "大凶",
	"木土木": "大凶",
	"木土火": "吉",
	"木土土": "吉",
	"木土金": "吉多凶少",
	"木土水": "大凶",
	"木金木": "大凶",
	"木金火": "大凶",
	"木金土": "凶多吉少",
	"木金金": "大凶",
	"木金水": "大凶",
	"木水木": "大吉",
	"木水火": "凶多吉少",
	"木水土": "凶多吉少",
	"木水金": "大吉",
	"木水水": "大吉",
	"火木木": "大吉",
	"火木火": "大吉",
	"火木土": "大吉",
	"火木金": "凶多吉少",
	"火木水": "吉",
	"火火木": "大吉",
	"火火火": "吉",
	"火火土": "大吉",
	"火火金": "大凶",
	"火火水": "大凶",
	"火土木": "吉多凶少",
	"火土火": "大吉",
	"火土土": "大吉",
	"火土金": "大吉",
	"火土水": "吉多凶少",
	"火金木": "大凶",
	"火金火": "大凶",
	"火金土": "吉凶参半",
	"火金金": "大凶",
	"火金水": "大凶",
	"火水木": "凶多吉少",
	"火水火": "大凶",
	"火水土": "大凶",
	"火水金": "大凶",
	"火水水": "大凶",
	"土木木": "吉",
	"土木火": "吉",
	"土木土": "凶多吉少",
	"土木金": "大凶",
	"土木水": "凶多吉少",
	"土火木": "大吉",
	"土火火": "大吉",
	"土火土": "大吉",
	"土火金": "吉多凶少",
	"土火水": "大凶",
	"土土木": "吉",
	"土土火": "大吉",
	"土土土": "大吉",
	"土土金": "大吉",
	"土土水": "凶多吉少",
	"土金木": "凶多吉少",
	"土金火": "凶多吉少",
	"土金土": "大吉",
	"土金金": "大吉",
	"土金水": "大吉",
	"土水木": "凶多吉少",
	"土水火": "大凶",
	"土水土": "大凶",
	"土水金": "吉凶参半",
	"土水水": "大凶",
	"金木木": "凶多吉少",
	"金木火": "凶多吉少",
	"金木土": "凶多吉少",
	"金木金": "大凶",
	"金木水": "凶多吉少",
	"金火木": "凶多吉少",
	"金火火": "吉凶参半",
	"金火土": "吉凶参半",
	"金火金": "大凶",
	"金火水": "大凶",
	"金土木": "吉",
	"金土火": "大吉",
	"金土土": "大吉",
	"金土金": "大吉",
	"金土水": "吉多凶少",
	"金金木": "大凶",
	"金金土": "大吉",
	"金金金": "吉",
	"金金水": "吉",
	"金水木": "大吉",
	"金水火": "凶多吉少",
	"金水土": "吉",
	"金水金": "大吉",
	"金水水": "吉",
	"水木木": "大吉",
	"水木火": "大吉",
	"水木土": "大吉",
	"水木金": "凶多吉少",
	"水木水": "大吉",
	"水火木": "吉",
	"水火火": "大凶",
	"水火土": "凶多吉少",
	"水火金": "大凶",
	"水火水": "大凶",
	"水土木": "大凶",
	"水土火": "吉",
	"水土土": "吉",
	"水土金": "吉",
	"水土水": "大凶",
	"水金木": "凶多吉少",
	"水金火": "凶多吉少",
	"水金土": "大吉",
	"水金金": "吉",
	"水金水": "大吉",
	"水水木": "大吉",
	"水水火": "大凶",
	"水水土": "大凶",
	"水水金": "大吉",
	"水水水": "吉",
}

var wuXingPoint = map[string]int{
	"大凶": 1, "凶": 2, "凶多吉少": 3, "吉凶参半": 4, "吉多凶少": 5, "吉": 6, "大吉": 7,
}

// NewWuXing creates a new WuXing object based on the given string.
//
// The parameter 's' is the string used to create the WuXing object.
// The function returns a pointer to the created WuXing object and a boolean value
// indicating whether the string exists in the wuXingTables map.
func NewWuXing(s string) (wx *WuXing, exist bool) {
	wxp, exist := wuXingTables[s]
	if !exist {
		return &WuXing{}, exist
	}
	p, exist := wuXingPoint[wxp]
	if !exist {
		return &WuXing{}, exist
	}
	return &WuXing{
		wuXing: s,
		point:  p,
	}, true
}
