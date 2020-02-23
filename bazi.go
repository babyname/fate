package fate

import (
	"github.com/godcong/chronos"
	"strings"
)

var diIndex = map[string]int{
	"子": 0, "丑": 1, "寅": 2, "卯": 3, "辰": 4, "巳": 5, "午": 6, "未": 7, "申": 8, "酉": 9, "戌": 10, "亥": 11,
}

var tianIndex = map[string]int{
	"甲": 0, "乙": 1, "丙": 2, "丁": 3, "戊": 4, "己": 5, "庚": 6, "辛": 7, "壬": 8, "癸": 9,
}

//天干强度表
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

//地支强度表
var dizhi = []map[string][]int{
	{
		"癸": {1200, 1100, 1000, 1000, 1040, 1060, 1000, 1000, 1200, 1200, 1060, 1140},
	}, {
		"癸": {360, 330, 300, 300, 312, 318, 300, 300, 360, 360, 318, 342},
		"辛": {200, 228, 200, 200, 230, 212, 200, 220, 228, 248, 232, 200},
		"己": {500, 550, 530, 500, 550, 570, 600, 580, 500, 500, 570, 500},
	}, {
		"丙": {300, 300, 360, 360, 318, 342, 360, 330, 300, 300, 342, 318},
		"甲": {840, 742, 798, 840, 770, 700, 700, 728, 742, 700, 700, 840},
	}, {
		"乙": {1200, 1060, 1140, 1200, 1100, 1000, 1000, 1040, 1060, 1000, 1000, 1200},
	}, {
		"乙": {360, 318, 342, 360, 330, 300, 300, 312, 318, 300, 300, 360},
		"癸": {240, 220, 200, 200, 208, 200, 200, 200, 240, 240, 212, 228},
		"戊": {500, 550, 530, 500, 550, 600, 600, 580, 500, 500, 570, 500},
	}, {
		"庚": {300, 342, 300, 300, 330, 300, 300, 330, 342, 360, 348, 300},
		"丙": {700, 700, 840, 840, 742, 840, 840, 798, 700, 700, 728, 742},
	}, {
		"丁": {1000, 1000, 1200, 1200, 1060, 1140, 1200, 1100, 1000, 1000, 1040, 1060},
	}, {
		"丁": {300, 300, 360, 360, 318, 342, 360, 330, 300, 300, 312, 318},
		"乙": {240, 212, 228, 240, 220, 200, 200, 208, 212, 200, 200, 240},
		"己": {500, 550, 530, 500, 550, 570, 600, 580, 500, 500, 570, 500},
	}, {
		"壬": {360, 330, 300, 300, 312, 318, 300, 300, 360, 360, 318, 342},
		"庚": {700, 798, 700, 700, 770, 742, 700, 770, 798, 840, 812, 700},
	}, {
		"辛": {1000, 1140, 1000, 1000, 1100, 1060, 1000, 1100, 1140, 1200, 1160, 1000},
	}, {
		"辛": {300, 342, 300, 300, 330, 318, 300, 330, 342, 360, 348, 300},
		"丁": {200, 200, 240, 240, 212, 228, 240, 220, 200, 200, 208, 212},
		"戊": {500, 550, 530, 500, 550, 570, 600, 580, 500, 500, 570, 500},
	}, {
		"甲": {360, 318, 342, 360, 330, 300, 300, 312, 318, 300, 300, 360},
		"壬": {840, 770, 700, 700, 728, 742, 700, 700, 840, 840, 724, 798},
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

//WuXingTianGan 五行天干
func WuXingTianGan(s string) string {
	return wuXingTianGan[s]
}

//WuXingDiZhi 五行地支
func WuXingDiZhi(s string) string {
	return wuXingDiZhi[s]
}

// BaZi ...
type BaZi struct {
	baZi   []string
	wuXing []string
	xiyong *XiYong
}

//NewBazi 创建八字
func NewBazi(calendar chronos.Calendar) *BaZi {
	ec := calendar.Lunar().EightCharacter()
	return &BaZi{
		baZi:   ec,
		wuXing: baziToWuXing(ec),
	}
}

// String ...
func (z *BaZi) String() string {
	return strings.Join(z.baZi, "")
}

//RiZhu 日主
func (z *BaZi) RiZhu() string {
	return z.baZi[4]
}

func (z *BaZi) calcXiYong() {
	z.xiyong = &XiYong{}
	//TODO:need fix
	z.point().calcSimilar().calcHeterogeneous() //.yongShen().xiShen()
}

//XiYong 喜用神
func (z *BaZi) XiYong() *XiYong {
	if z.xiyong == nil {
		z.calcXiYong()
	}
	return z.xiyong
}

//XiYongShen 平衡用神
func (z *BaZi) XiYongShen() string {
	return z.XiYong().Shen()
}

//func (z *BaZi) yongShen() *BaZi {
//	z.xiyong.YongShen = z.xiyong.Similar[0]
//	return z
//}
//func (z *BaZi) xiShen() *BaZi {
//	rt := sheng
//	if z.QiangRuo() {
//		rt = ke
//	}
//	for i := range rt {
//		if rt[i] == z.xiyong.YongShen {
//			if i == len(rt) {
//				i = -1
//			}
//			z.xiyong.XiShen = rt[i-1]
//			break
//		}
//	}
//	return z
//}
func (z *BaZi) point() *BaZi {
	di := diIndex[z.baZi[3]]
	for idx, v := range z.baZi {
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

//计算同类
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

//计算异类
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
