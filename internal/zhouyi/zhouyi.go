package zhouyi

import (
	"github.com/godcong/yi"
)

func QiGua(xia, shang int) *yi.Yi {
	return yi.NumberQiGua(shang, xia)
}
