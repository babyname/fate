package fate

import (
	"github.com/godcong/fate/mongo"
)

const (
	BenGua = iota
	BianGua
	HuGua
)

const (
	ShangQian = 0x00
	ShangDui  = 0x01
	ShangLi   = 0x02
	ShangZhen = 0x03
	ShangXun  = 0x04
	ShangKan  = 0x05
	ShangGen  = 0x06
	ShangKun  = 0x07
	XiaQian   = 0x00
	XiaDui    = 0x10
	XiaLi     = 0x20
	XiaZhen   = 0x30
	XiaXun    = 0x40
	XiaKan    = 0x50
	XiaGen    = 0x60
	XiaKun    = 0x70
)

type ZhouYi struct {
	gua [3]mongo.GuaXiang
}

//QiGua 起卦
func QiGua(name *Name) *ZhouYi {
	x := CountStroke(name.lastChar...)
	m := CountStroke(name.firstChar...)
	getGua(x)
	getGua(m)
	return nil
}

func getGua(i int) {
	panic("")
}

func ShangGua(){
	panic("")
}

func XiaGua()  {
	panic("")
}
