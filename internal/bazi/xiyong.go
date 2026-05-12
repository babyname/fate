package bazi

import (
	"math"
	"strings"

	"github.com/babyname/fate/ent"
)

type XiYong struct {
	WuXingFen          map[string]int
	Similar            []string
	SimilarPoint       int
	Heterogeneous      []string
	HeterogeneousPoint int
}

var sheng = []string{"木", "火", "土", "金", "水"}
var ke = []string{"木", "土", "水", "火", "金"}

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

func (xy *XiYong) Shen() string {
	if !xy.QiangRuo() {
		return xy.minFenWuXing(xy.Similar...)
	}
	return xy.minFenWuXing(xy.Heterogeneous...)
}

func (xy *XiYong) QiangRuo() bool {
	return xy.SimilarPoint > xy.HeterogeneousPoint
}

func filterXiYong(yong string, cs ...*ent.Character) (b bool) {
	for _, c := range cs {
		if strings.Contains(yong, c.WuXing) {
			return true
		}
	}
	return false
}
