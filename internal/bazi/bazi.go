package bazi

import (
	v2 "github.com/godcong/chronos/v2"
)

var diIndex = map[string]int{
	"子": 0, "丑": 1, "寅": 2, "卯": 3, "辰": 4, "巳": 5, "午": 6, "未": 7, "申": 8, "酉": 9, "戌": 10, "亥": 11,
}

var tianIndex = map[string]int{
	"甲": 0, "乙": 1, "丙": 2, "丁": 3, "戊": 4, "己": 5, "庚": 6, "辛": 7, "壬": 8, "癸": 9,
}

var tiangan = [][]int{
	{1200, 1200, 1000, 1000, 1000, 1000, 1000, 1000, 1200, 1200},
	{1060, 1060, 1000, 1000, 1100, 1100, 1140, 1140, 1100, 1100},
	{1140, 1140, 1200, 1200, 1060, 1060, 1000, 1000, 1000, 1000},
	{1200, 1200, 1200, 1200, 1000, 1000, 1000, 1000, 1000, 1000},
	{1100, 1100, 1060, 1060, 1100, 1100, 1100, 1100, 1040, 1040},
	{1000, 1000, 1140, 1140, 1140, 1140, 1060, 1060, 1060, 1060},
	{1000, 1000, 1200, 1200, 1200, 1200, 1000, 1000, 1000, 1000},
	{1040, 1040, 1100, 1100, 1160, 1160, 1100, 1100, 1000, 1000},
	{1060, 1060, 1000, 1000, 1000, 1000, 1140, 1140, 1200, 1200},
	{1000, 1000, 1000, 1000, 1000, 1000, 1200, 1200, 1200, 1200},
	{1000, 1000, 1040, 1040, 1140, 1140, 1160, 1160, 1060, 1060},
	{1200, 1200, 1000, 1000, 1000, 1000, 1000, 1000, 1140, 1140},
}

var dizhi = []map[string][]int{
	{
		"癸": {1200, 1100, 1000, 1000, 1040, 1060, 1000, 1000, 1200, 1200, 1060, 1140},
	}, {
		"癸": {360, 330, 300, 300, 312, 318, 300, 300, 360, 360, 318, 342},
		"辛": {200, 228, 200, 200, 220, 212, 200, 220, 228, 240, 232, 200},
		"己": {500, 550, 530, 500, 550, 570, 600, 580, 500, 500, 570, 500},
	}, {
		"甲": {600, 530, 570, 600, 550, 500, 500, 520, 530, 500, 500, 600},
		"丙": {300, 300, 360, 360, 318, 342, 360, 330, 300, 300, 312, 300},
		"戊": {200, 220, 212, 200, 220, 228, 240, 232, 200, 200, 228, 200},
	}, {
		"乙": {1200, 1060, 1140, 1200, 1100, 1000, 1000, 1040, 1060, 1000, 1000, 1200},
	}, {
		"乙": {360, 318, 342, 360, 330, 300, 300, 312, 318, 300, 300, 360},
		"癸": {240, 220, 200, 200, 208, 212, 200, 200, 240, 240, 212, 228},
		"戊": {500, 550, 530, 500, 550, 570, 600, 580, 500, 500, 570, 500},
	}, {
		"丙": {500, 500, 600, 600, 530, 570, 600, 550, 500, 500, 520, 500},
		"庚": {300, 342, 300, 300, 330, 318, 300, 330, 342, 360, 348, 300},
		"戊": {200, 220, 212, 200, 220, 228, 240, 232, 200, 200, 228, 200},
	}, {
		"丁": {700, 700, 840, 840, 742, 798, 840, 770, 700, 700, 728, 700},
		"己": {300, 330, 318, 300, 330, 342, 360, 348, 300, 300, 342, 300},
	}, {
		"丁": {300, 300, 360, 360, 318, 342, 360, 330, 300, 300, 312, 300},
		"乙": {240, 212, 228, 240, 220, 200, 200, 208, 212, 200, 200, 240},
		"己": {500, 550, 530, 500, 550, 570, 600, 580, 500, 500, 570, 500},
	}, {
		"庚": {500, 570, 500, 500, 550, 530, 500, 550, 570, 600, 580, 500},
		"壬": {360, 330, 300, 300, 312, 318, 300, 300, 360, 360, 318, 342},
		"戊": {200, 220, 212, 200, 220, 228, 240, 232, 200, 200, 228, 200},
	}, {
		"辛": {1000, 1140, 1000, 1000, 1100, 1060, 1000, 1100, 1140, 1200, 1160, 1000},
	}, {
		"辛": {300, 342, 300, 300, 330, 318, 300, 330, 342, 360, 348, 300},
		"丁": {200, 200, 240, 240, 212, 228, 240, 220, 200, 200, 208, 200},
		"戊": {500, 550, 530, 500, 550, 570, 600, 580, 500, 500, 570, 500},
	}, {
		"甲": {360, 318, 342, 360, 330, 300, 300, 312, 318, 300, 300, 360},
		"壬": {840, 770, 700, 700, 728, 742, 700, 700, 840, 840, 742, 798},
	},
}

var wuXingTianGan = map[string]string{
	"甲": "木",
	"乙": "木",
	"丙": "火",
	"丁": "火",
	"戊": "土",
	"己": "土",
	"庚": "金",
	"辛": "金",
	"壬": "水",
	"癸": "水",
}

var wuXingDiZhi = map[string]string{
	"子": "水",
	"丑": "土",
	"寅": "木",
	"卯": "木",
	"辰": "土",
	"巳": "火",
	"午": "火",
	"未": "土",
	"申": "金",
	"酉": "金",
	"戌": "土",
	"亥": "水",
}

func WuXingTianGan(s string) string {
	return wuXingTianGan[s]
}

func WuXingDiZhi(s string) string {
	return wuXingDiZhi[s]
}

type BaZi struct {
	baZi   v2.EightChar
	wuXing [4]string
	xiyong *XiYong
	bridge *v2.Bridge
}

func NewBazi(calendar v2.Calendar) *BaZi {
	ec := calendar.Lunar().GetEightChar()
	return &BaZi{
		baZi:   ec,
		wuXing: ec.GetWuXing(),
	}
}

func NewBaziFromBridge(bridge *v2.Bridge) *BaZi {
	ec := bridge.EightChar()
	return &BaZi{
		baZi:   ec,
		wuXing: ec.GetWuXing(),
		bridge: bridge,
	}
}

func (z *BaZi) String() string {
	siZhu := z.baZi.GetSiZhu()
	return siZhu[0] + siZhu[1] + siZhu[2] + siZhu[3]
}

func (z *BaZi) RiZhu() string {
	return z.baZi.GetSiZhu()[3]
}

func (z *BaZi) calcXiYong() {
	z.xiyong = &XiYong{}
	z.point().calcSimilar().calcHeterogeneous()
}

func (z *BaZi) XiYong() *XiYong {
	if z.xiyong == nil {
		z.calcXiYong()
	}
	return z.xiyong
}

func (z *BaZi) XiYongShen() string {
	return z.XiYong().Shen()
}

func (z *BaZi) point() *BaZi {
	di := diIndex[z.baZi.GetSiZhu()[2]]
	for idx, v := range z.baZi.GetSiZhu() {
		if idx%2 == 0 {
			z.xiyong.AddFen(WuXingTianGan(v), tiangan[di][tianIndex[v]])
		} else {
			dz := dizhi[diIndex[v]]
			for k := range dz {
				z.xiyong.AddFen(WuXingTianGan(k), dz[k][di])
			}
		}
	}
	return z
}

func baziToWuXing(bazi []string) []string {
	var wx []string
	for idx, v := range bazi {
		if idx%2 == 0 {
			wx = append(wx, WuXingTianGan(v))
		} else {
			wx = append(wx, WuXingDiZhi(v))
		}
	}
	return wx
}

func (z *BaZi) calcSimilar() *BaZi {
	for i := range sheng {
		if wuXingTianGan[z.RiZhu()] == sheng[i] {
			z.xiyong.Similar = append(z.xiyong.Similar, sheng[i])
			z.xiyong.SimilarPoint = z.xiyong.GetFen(sheng[i])
			if i == 0 {
				i = len(sheng) - 1
				z.xiyong.Similar = append(z.xiyong.Similar, sheng[i])
				z.xiyong.SimilarPoint += z.xiyong.GetFen(sheng[i])
			} else {
				z.xiyong.Similar = append(z.xiyong.Similar, sheng[i-1])
				z.xiyong.SimilarPoint += z.xiyong.GetFen(sheng[i-1])
			}
			break
		}
	}
	return z
}

func (z *BaZi) calcHeterogeneous() *BaZi {
	for i := range sheng {
		for ti := range z.xiyong.Similar {
			if z.xiyong.Similar[ti] == sheng[i] {
				goto EndSimilar
			}
		}
		z.xiyong.Heterogeneous = append(z.xiyong.Heterogeneous, sheng[i])
		z.xiyong.HeterogeneousPoint += z.xiyong.GetFen(sheng[i])
	EndSimilar:
		continue
	}
	return z
}
