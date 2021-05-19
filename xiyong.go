package fate

import (
	"math"
	"strings"

	"github.com/godcong/yi"
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

//过滤喜用神，目前支持倒序前两个
//输入的喜用神序列必须是有倾向的，如果是若干种五行平分，则不能进行硬过滤
func fixedFilterXiYong(yong string, hard_filter bool, cs ...*Character) (b bool) {
	wuxing_count_yong := yi.CountWuxing(nil, yong)

	wuxing_count_yong = yi.SortWuxingCount(wuxing_count_yong)

	if wuxing_count_yong.Content[0].Count == 0 {
		panic("没有输入喜用神")
	}

	wuxing_count_got := yi.CountWuxing(nil, "")

	for _, c := range cs {
		wuxing_count_got = yi.CountWuxing(wuxing_count_got, c.WuXing)
	}

	wuxing_count_got = yi.SortWuxingCount(wuxing_count_got)

	if wuxing_count_yong.Content[1].Count == 0 {
		if wuxing_count_yong.Content[0].WuXingChr == wuxing_count_got.Content[0].WuXingChr {
			return true
		}
	} else {
		if wuxing_count_yong.Content[1].WuXingChr == wuxing_count_got.Content[1].WuXingChr {
			return true
		} else {
			if !hard_filter {
				if wuxing_count_yong.Content[0].WuXingChr == wuxing_count_got.Content[1].WuXingChr &&
					wuxing_count_yong.Content[1].WuXingChr == wuxing_count_got.Content[0].WuXingChr {
					return true
				}
			}
		}
	}

	return false
}
