package bazi

import (
	"math"
)

// XiYong 表示八字喜用神的计算结果，包含五行分数及同类/异类信息。
type XiYong struct {
	WuXingFen          map[string]int
	Similar            []string
	SimilarPoint       int
	Heterogeneous      []string
	HeterogeneousPoint int
}

var sheng = []string{"木", "火", "土", "金", "水"}

// AddFen 为指定五行累加分数。
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

// GetFen 获取指定五行的累计分数。
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
	minFen := math.MaxInt32
	for _, s := range ss {
		if xy.WuXingFen[s] < minFen {
			minFen = xy.WuXingFen[s]
			wx = s
		} else if xy.WuXingFen[s] == minFen {
			wx += s
		}
	}
	return
}

// Shen 返回喜用神对应的五行属性。
func (xy *XiYong) Shen() string {
	if !xy.QiangRuo() {
		return xy.minFenWuXing(xy.Similar...)
	}
	return xy.minFenWuXing(xy.Heterogeneous...)
}

// QiangRuo 判断日主强弱，同类分数大于异类则为强。
func (xy *XiYong) QiangRuo() bool {
	return xy.SimilarPoint > xy.HeterogeneousPoint
}
