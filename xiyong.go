package fate

import (
	"fmt"
	"math"
	"strings"
)

//XiYong 喜用神
type XiYong struct {
	WuXingFen          map[string]int
	Similar            []string //同类
	SimilarPoint       int
	Heterogeneous      []string //异类
	HeterogeneousPoint int
}

var sheng = []string{"木", "火", "土", "金", "水"}
var ke = []string{"木", "土", "水", "火", "金"}

//AddFen 五行分
func (xy *XiYong) AddFen(s string, point int) {
	if xy.WuXingFen == nil {
		xy.WuXingFen = make(map[string]int)
	}

	if v, b := xy.WuXingFen[s]; b {
		xy.WuXingFen[s] = v + point
	} else {
		xy.WuXingFen[s] = point
	}
}

//GetFen 取得分
func (xy *XiYong) GetFen(s string) (point int) {
	if xy.WuXingFen == nil {
		return 0
	}
	if v, b := xy.WuXingFen[s]; b {
		return v
	}
	return 0
}

func (xy *XiYong) minFenWuXing(ss ...string) (wx string) {
	min := math.MaxInt32
	for _, s := range ss {
		if xy.WuXingFen[s] < min {
			min = xy.WuXingFen[s]
			wx = s
		} else if xy.WuXingFen[s] == min {
			wx += s
		}
	}
	return
}

//Shen 喜用神
func (xy *XiYong) Shen() string {
	if !xy.QiangRuo() {
		return xy.minFenWuXing(xy.Similar...)
	}
	return xy.minFenWuXing(xy.Heterogeneous...)
}

//QiangRuo 八字偏强（true)弱（false）
func (xy *XiYong) QiangRuo() bool {
	return xy.SimilarPoint > xy.HeterogeneousPoint
}

func filterXiYong(yong string, cs ...*Character) (b bool) {
	for _, c := range cs {
		if strings.Contains(yong, c.WuXing) {
			return true
		}
	}
	return false
}

func fixedFilterXiYong(yong string, cs ...*Character) (b bool) {
	jin := 0
	mu := 0
	shui := 0
	huo := 0
	tu := 0

	yong_runes := []rune(yong)

	for idx, yong_rune := range yong_runes {
		yong_str := string(yong_rune)
		switch yong_str {
		case "金":
			jin++
		case "木":
			mu++
		case "水":
			shui++
		case "火":
			huo++
		case "土":
			tu++
		default:
			panic(fmt.Sprintf("喜用神不支持：%d:%s", idx, yong_str))
		}
	}

	jin_got := 0
	mu_got := 0
	shui_got := 0
	huo_got := 0
	tu_got := 0

	for _, c := range cs {
		if strings.Contains("金", c.WuXing) {
			jin_got++
		}

		if strings.Contains("木", c.WuXing) {
			mu_got++
		}

		if strings.Contains("水", c.WuXing) {
			shui_got++
		}

		if strings.Contains("火", c.WuXing) {
			huo_got++
		}

		if strings.Contains("土", c.WuXing) {
			tu_got++
		}
	}

	if jin == jin_got && mu == mu_got && shui == shui_got && huo == huo_got && tu == tu_got {
		return true
	}

	return false
}
