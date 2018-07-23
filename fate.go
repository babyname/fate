package fate

import (
	"time"

	"github.com/godcong/chronos"
	"github.com/godcong/fate/config"
	"github.com/godcong/fate/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
)

type fate struct {
	name     *Name
	martial  *Martial
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

func (f *fate) GetName() *Name {
	return f.name
}

func (f *fate) SetMartial(martial *Martial) {
	f.martial = martial
}

func (f *fate) GetMartial() *Martial {
	if f.martial == nil {
		return &Martial{}
	}
	return f.martial
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
	s, _ := calculatorBestStroke(f, f.name.lastChar)
	return s
}

func (f *fate) BestCharacters() []*mongo.Character {
	var cs []*mongo.Character

	strokes := f.BestStrokes()
	var firsts []string
	for idx := range strokes {
		firsts = append(firsts, strconv.Itoa(strokes[idx].FirstStroke[0]))
	}

	err := mongo.C("character").Find(bson.M{
		"total_trokes": bson.M{"$in": firsts},
	}).All(&cs)
	log.Println(err)
	return cs
}
