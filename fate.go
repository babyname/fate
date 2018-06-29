package fate

import (
	"time"

	"github.com/godcong/chronos"
)

type fate struct {
	name     *Name
	calendar chronos.Calendar
	property *Property
	fn1      int
	fn2      int
}

func NewFate(lastName string) *fate {
	name := newName(lastName)
	return &fate{name: name}
}

func (f *fate) SetLastName(lastName string) {
	f.name = newName(lastName)
}

//SetLunarData 设定生日
func (f *fate) SetLunarData(t time.Time) {
	f.calendar = chronos.New(t)
}

func (f *fate) SetProperty(p *Property) {
	f.property = p
}

func (f *fate) MakeFirstName() *Name {

	return f.name
}

//EightCharacter 计算生辰八字(需先设定生日),按年柱,月柱,日柱,时柱 输出
func (f *fate) EightCharacter() (string, string, string, string) {
	if f.calendar != nil {
		return f.calendar.Lunar().EightCharacter()
	}
	return "", "", "", ""
}

func (f *fate) ThreeTalent() {

	//tt := NewThreeTalent(fg)
}

func (f *fate) FiveGrid() *FiveGrid {
	var fg *FiveGrid
	if len(f.name.cLast) > 1 {
		fg = MakeFiveGridFromStrokes(f.name.cLast[0].ScienceStrokes, f.name.cLast[1].ScienceStrokes, f.fn1, f.fn2)
	} else {
		fg = MakeFiveGridFromStrokes(f.name.cLast[0].ScienceStrokes, 0, f.fn1, f.fn2)
	}
	return fg
}

func (f *fate) IteratorFiveGrid() *FiveGrid {
	var fg *FiveGrid
	if len(f.name.cLast) > 1 {
		fg = MakeFiveGridFromStrokes(f.name.cLast[0].ScienceStrokes, f.name.cLast[1].ScienceStrokes, f.fn1, f.fn2)
	} else {
		fg = MakeFiveGridFromStrokes(f.name.cLast[0].ScienceStrokes, 0, f.fn1, f.fn2)
	}
	if f.fn2 > LenMax {
		f.fn2 = 0
		f.fn1++
		return fg
	}
	f.fn2++
	return fg
}
