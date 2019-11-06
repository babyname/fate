package fate

import (
	"github.com/godcong/chronos"
	"github.com/godcong/yi"
)

//Name 姓名
type Name struct {
	FirstName   []*Character //名姓
	LastName    []*Character
	born        *chronos.Calendar
	baGua       *yi.Yi //周易八卦
	zodiac      *Zodiac
	zodiacPoint int
}

func (n Name) String() string {
	var s string
	for _, l := range n.LastName {
		s += l.Ch
	}
	for _, f := range n.FirstName {
		s += f.Ch
	}
	return s
}

func (n Name) WuXing() string {
	var s string
	for _, l := range n.LastName {
		s += l.WuXing
	}
	for _, f := range n.FirstName {
		s += f.WuXing
	}
	return s
}

func createName(impl *fateImpl, f1 *Character, f2 *Character) *Name {
	lastSize := len(impl.lastChar)
	last := make([]*Character, lastSize, lastSize)
	copy(last, impl.lastChar)
	ff1 := *f1
	ff2 := *f2
	first := []*Character{&ff1, &ff2}

	return &Name{
		FirstName: first,
		LastName:  last,
	}
}

//BaGua
func (n *Name) BaGua() *yi.Yi {
	if n.baGua == nil {
		lastSize := len(n.LastName)
		shang := getStroke(n.LastName[0])
		if lastSize > 1 {
			shang += getStroke(n.LastName[1])
		}
		xia := getStroke(n.FirstName[0]) + getStroke(n.FirstName[1])
		n.baGua = yi.NumberQiGua(xia, shang, shang+xia)
	}

	return n.baGua
}
