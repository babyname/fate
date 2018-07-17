package fate

import (
	"time"

	"github.com/godcong/chronos"
)

type fate struct {
	name     *Name
	calendar chronos.Calendar
}

//MaxStokers 超过32划的字不易书写,过滤
const MaxStokers = 32

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

//EightCharacter 计算生辰八字(需要SetLunarData),按年柱,月柱,日柱,时柱 输出
func (f *fate) EightCharacter() (string, string, string, string) {
	if f.calendar != nil {
		return f.calendar.Lunar().EightCharacter()
	}
	return "", "", "", ""
}
