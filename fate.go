package fate

import (
	"time"

	"github.com/godcong/chronos"
)

type fate struct {
	name     *Name
	calendar chronos.Calendar
	property *Property
}

func NewFate(lastName string) *fate {
	name := newName(lastName)
	return &fate{name: name}
}

func (f *fate) SetLastName(lastName string) {
	f.name = newName(lastName)
}

func (f *fate) SetLunarData(t time.Time) {
	f.calendar = chronos.New(t)
}

func (f *fate) SetProperty(p *Property) {
	f.property = p
}

func (f *fate) MakeFirstName() *Name {

	return f.name
}

//EightCharacter 计算生辰八字
func (f *fate) EightCharacter() (string, string, string, string) {
	return f.calendar.Lunar().EightCharacter()
}
