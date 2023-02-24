package fate

import (
	"github.com/godcong/chronos/v2"
)

type NameDetail struct {
	date chronos.Calendar
}

func (n *NameDetail) EightChar() chronos.EightChar {
	return n.date.Lunar().GetEightChar()
}

//func (n NameDetail) name()  {
//
//}
