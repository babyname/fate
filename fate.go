package fate

import (
	"time"

	"github.com/godcong/chronos"
	"github.com/godcong/fate/config"
	"github.com/godcong/fate/mongo"
	"gopkg.in/mgo.v2"
)

type fate struct {
	name     *Name
	calendar chronos.Calendar
}

func init() {
	initDial()
}

func initDial() {
	mongo.Dial(config.Default().GetString("mongodb.url"), &mgo.Credential{
		Username:    config.Default().GetString("mongodb.username"),
		Password:    config.Default().GetString("mongodb.password"),
		Source:      "",
		Service:     "",
		ServiceHost: "",
		Mechanism:   "",
	})
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

func (f *fate) BestStrokes() []*Stroke {
	s, _ := calculatorBestStroke(f.name.lastChar)
	return s
}
