package zhouyi

import (
	"github.com/babyname/yi"
)

// QiGua 根据下卦和上卦数起卦。
func QiGua(xia, shang int) *yi.ZhouYi {
	return yi.DivineByNumber(shang, xia)
}
