package fate

import (
	"time"

	"github.com/godcong/chronos"
)

type fate struct {
	name        *Name
	fiveGrid    *FiveGrid
	threeTalent *ThreeTalent
	calendar    chronos.Calendar
	property    *Property
	fn1         int
	fn2         int
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

//EightCharacter 计算生辰八字(需要SetLunarData),按年柱,月柱,日柱,时柱 输出
func (f *fate) EightCharacter() (string, string, string, string) {
	if f.calendar != nil {
		return f.calendar.Lunar().EightCharacter()
	}
	return "", "", "", ""
}

//GenerateName 通过规则生成姓名
func (f *fate) GenerateName() {

}

func (f *fate) ThreeTalent() *ThreeTalent {
	f.threeTalent = NewThreeTalent(f.fiveGrid)
	return f.threeTalent
}

func (f *fate) FiveGrid() *FiveGrid {
	f.fiveGrid = currentFiveGrid(f.name, f.fn1, f.fn2)
	return f.fiveGrid
}

func currentFiveGrid(name *Name, fn1, fn2 int) *FiveGrid {
	var fg *FiveGrid
	if len(name.cLast) > 1 {
		fg = MakeFiveGridFromStrokes(name.cLast[0].ScienceStrokes, name.cLast[1].ScienceStrokes, fn1, fn2)
	} else {
		fg = MakeFiveGridFromStrokes(name.cLast[0].ScienceStrokes, 0, fn1, fn2)
	}
	return fg
}
