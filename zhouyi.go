package fate

import (
	"github.com/godcong/yi"
)

//QiGua 起卦
func QiGua(xia, shang int) *yi.Yi {
	return yi.NumberQiGua(shang, xia)
}
