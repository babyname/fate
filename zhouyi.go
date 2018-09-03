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

var fu = []string{"☰", "☱", "☲", "☳", "☴", "☵", "☶", "☷"}

var gua = [...]string{
	GuaXiangKun:  "坤",
	GuaXiangGen:  "艮",
	GuaXiangKan:  "坎",
	GuaXiangXun:  "巽",
	GuaXiangZhen: "震",
	GuaXiangLi:   "离",
	GuaXiangDui:  "兑",
	GuaXiangQian: "乾",
}

//ZhouYi 周易卦象
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
	hu := ben
	if ben.ShangShu == GuaXiangKun && ben.XiaShu == GuaXiangKun ||
		ben.ShangShu == GuaXiangQian && ben.XiaShu == GuaXiangQian {
		hu = bian
	}
	hu = huGua(hu)
	return &ZhouYi{
		gua: [3]*mongo.GuaXiang{
			BenGua:  ben,
			BianGua: bian,
			HuGua:   hu,
		},
	}

	return nil
}
//Set 设定卦象
func (y *ZhouYi) Set(idx int, xiang *mongo.GuaXiang) {
	y.gua[idx] = xiang
}

func getGua(i int) string {
	if i = i % 8; i != 0 {
		return gua[i-1]
	}
	return gua[7]
}
//取爻
func getYao(i int) int {
	if i = i % 6; i != 0 {
		return i - 1
	}
	return 5
}

//本卦
func benGua(x, m int) *mongo.GuaXiang {
	bg := strings.Join([]string{getGua(x), getGua(m)}, "")
	gx := mongo.GetGuaXiang()
	if v, b := gx[bg]; b {
		return v
	}
	return nil
}
//变卦
func bianGua(ben *mongo.GuaXiang, b int) *mongo.GuaXiang {
	gx := mongo.GetGuaXiang()
	bz := getYao(b)
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

//变
func bian(gua, bian int) int {
	idx := 1 << uint(bian)
	if gua&idx == 0 {
		gua = gua | idx
	} else {
		gua = gua ^ idx
	}
	return gua
}

//互
func hu(shang, xia int) int {
	huXia := 0
	er := 1 << 1
	san := 1 << 0
	si := 1 << 2
	if xia&er > 0 {
		huXia |= 1 << 2
	}
	if xia&san > 0 {
		huXia |= 1 << 1
	}
	if shang&si > 0 {
		huXia |= 1 << 0
	}
	return huXia
}

//交
func jiao(shang, xia int) int {
	jiaoShang := 0
	san := 1 << 0
	si := 1 << 2
	wu := 1 << 1
	if xia&san > 0 {
		//位移2
		jiaoShang |= 1 << 2
	}
	if shang&si > 0 {
		//位移1
		jiaoShang |= 1 << 1
	}
	if shang&wu > 0 {
		//位不动
		jiaoShang |= 1 << 0
	}
	return jiaoShang
}

//错
func cuo(gua int) int {
	gua ^= 0x7
	return gua
}

//综
func zong(shang, xia int) (int, int) {
	zShang := 0
	zXia := 0
	if (shang & (1 << 2)) > 0 {
		zXia |= 1 << 0
	}
	if (shang & (1 << 1)) > 0 {
		zXia |= 1 << 1
	}
	if (shang & (1 << 0)) > 0 {
		zXia |= 1 << 2
	}
	if (xia & (1 << 2)) > 0 {
		zShang |= 1 << 0
	}
	if (xia & (1 << 1)) > 0 {
		zShang |= 1 << 1
	}
	if (xia & (1 << 0)) > 0 {
		zShang |= 1 << 2
	}
	return zShang, zXia
}
//互卦
func huGua(ben *mongo.GuaXiang) *mongo.GuaXiang {
	bg := strings.Join([]string{getGua(jiao(ben.ShangShu, ben.XiaShu)), getGua(hu(ben.ShangShu, ben.XiaShu))}, "")
	gx := mongo.GetGuaXiang()
	if v, b := gx[bg]; b {
		return v
	}
	return nil
}
