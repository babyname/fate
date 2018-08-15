package fate

import "github.com/godcong/chronos"

//QiGua 起卦
func QiGua(name *Name, c chronos.Calendar) {

}

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