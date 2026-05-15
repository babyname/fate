package zhouyi

import (
	"github.com/godcong/yi"
)

// QiGua 根据下卦和上卦数起卦。
func QiGua(xia, shang int) *yi.Yi {
	return yi.NumberQiGua(shang, xia)
}
