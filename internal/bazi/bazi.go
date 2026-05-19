// Package bazi 实现八字命理计算，包括五行分析、喜用神推算等核心功能。
package bazi

import (
	v2 "github.com/babyname/chronos/v2"
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

// WuXingTianGan 根据天干返回对应的五行属性。
func WuXingTianGan(s string) string {
	return wuXingTianGan[s]
}

// WuXingDiZhi 根据地支返回对应的五行属性。
func WuXingDiZhi(s string) string {
	return wuXingDiZhi[s]
}

// BaZi 表示一个八字命盘，包含四柱、五行及喜用神信息。
type BaZi struct {
	baZi   v2.EightChar
	wuXing [4]string
	xiyong *XiYong
	bridge *v2.Bridge
}

// NewBazi 根据日历创建八字命盘实例。
func NewBazi(calendar v2.Calendar) *BaZi {
	ec := calendar.Lunar().GetEightChar()
	return &BaZi{
		baZi:   ec,
		wuXing: ec.GetWuXing(),
	}
}

// NewBaziFromBridge 根据桥接对象创建八字命盘实例。
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

// RiZhu 返回日柱的天干地支字符串。
func (z *BaZi) RiZhu() string {
	return z.baZi.GetSiZhu()[3]
}

func (z *BaZi) calcXiYong() {
	z.xiyong = &XiYong{}
	z.point().calcSimilar().calcHeterogeneous()
}

// XiYong 返回喜用神信息，若尚未计算则自动触发计算。
func (z *BaZi) XiYong() *XiYong {
	if z.xiyong == nil {
		z.calcXiYong()
	}
	return z.xiyong
}

// XiYongShen 返回喜用神对应的五行属性。
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

func (z *BaZi) calcSimilar() *BaZi {
	riZhuWX := wuXingTianGan[z.RiZhu()]
	for i, wx := range sheng {
		if riZhuWX == wx {
			z.xiyong.Similar = append(z.xiyong.Similar, wx)
			z.xiyong.SimilarPoint = z.xiyong.GetFen(wx)
			shengWX := sheng[(i+4)%5]
			z.xiyong.Similar = append(z.xiyong.Similar, shengWX)
			z.xiyong.SimilarPoint += z.xiyong.GetFen(shengWX)
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
