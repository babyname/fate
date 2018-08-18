package fate

import (
	"strings"

	"github.com/godcong/fate/mongo"
)

const (
	BenGua = iota
	BianGua
	HuGua
)

const (
	GuaXiangQian = 0x00
	GuaXiangDui  = 0x01
	GuaXiangLi   = 0x02
	GuaXiangZhen = 0x03
	GuaXiangXun  = 0x04
	GuaXiangKan  = 0x05
	GuaXiangGen  = 0x06
	GuaXiangKun  = 0x07
	//ShangQian = 0x00
	//ShangDui  = 0x01
	//ShangLi   = 0x02
	//ShangZhen = 0x03
	//ShangXun  = 0x04
	//ShangKan  = 0x05
	//ShangGen  = 0x06
	//ShangKun  = 0x07
	//XiaQian   = 0x00
	//XiaDui    = 0x10
	//XiaLi     = 0x20
	//XiaZhen   = 0x30
	//XiaXun    = 0x40
	//XiaKan    = 0x50
	//XiaGen    = 0x60
	//XiaKun    = 0x70
)

var gua = []string{
	GuaXiangKun:  "坤",
	GuaXiangGen:  "艮",
	GuaXiangKan:  "坎",
	GuaXiangXun:  "巽",
	GuaXiangZhen: "震",
	GuaXiangLi:   "离",
	GuaXiangDui:  "兑",
	GuaXiangQian: "乾",
}

type ZhouYi struct {
	gua [3]*mongo.GuaXiang
}

//QiGua 起卦
func QiGua(name *Name) *ZhouYi {
	x := CountStroke(name.lastChar...)
	m := CountStroke(name.firstChar...)
	b := CountStroke(append(name.lastChar, name.firstChar...)...)
	ben := benGua(x, m)
	bian := bianGua(ben, b)
	return &ZhouYi{
		gua: [3]*mongo.GuaXiang{
			BenGua:  ben,
			BianGua: bian,
		},
	}

	return nil
}

func (y *ZhouYi) Set(idx int, xiang *mongo.GuaXiang) {
	y.gua[idx] = xiang
}

func getGua(i int) string {
	if i = i % 8; i != 0 {
		return gua[i-1]
	}
	return gua[7]
}

func getZhou(i int) int {
	if i = i % 6; i != 0 {
		return i - 1
	}
	return 5
}

func benGua(x, m int) *mongo.GuaXiang {
	bg := strings.Join([]string{getGua(x), getGua(m)}, "")
	gx := mongo.GetGuaXiang()
	if v, b := gx[bg]; b {
		return v
	}
	return nil
}

func bianGua(ben *mongo.GuaXiang, b int) *mongo.GuaXiang {
	gx := mongo.GetGuaXiang()
	bz := getZhou(b)
	sg := ben.ShangGua
	xg := ben.XiaGua
	if b > 2 {
		sg = gua[bian(ben.ShangShu, bz-3)]
	} else {
		xg = gua[bian(ben.XiaShu, bz-3)]
	}
	gua := strings.Join([]string{sg, xg}, "")
	return gx[gua]
}

func bian(gua, bian int) int {
	idx := 1 << uint(bian)
	if gua&idx == 0 {
		gua = gua | idx
	} else {
		gua = gua ^ idx
	}
	return gua
}

func huGua(ben *mongo.GuaXiang, ) {
	panic("")
}
