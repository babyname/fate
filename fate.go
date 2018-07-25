package fate

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/godcong/chronos"
	"github.com/godcong/fate/config"
	"github.com/godcong/fate/mongo"
)

type fate struct {
	name     *Name
	martial  *Martial
	strokes  []*Stroke
	firstOne []*mongo.Character
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

func (f *fate) BestFirstOne() *mongo.Character {
	var cs []*mongo.Character

	f.strokes = f.BestStrokes()
	firsts := make(map[int][]byte)
	for idx := range f.strokes {
		firsts[f.strokes[idx].FirstStroke[0]] = nil
	}

	var charStroke []string
	for i := range firsts {
		charStroke = append(charStroke, strconv.Itoa(i))
	}

	err := mongo.C("character").Find(bson.M{
		"total_strokes":            bson.M{"$in": charStroke},
		"folk.commonly_characters": bson.M{"$in": []string{"是", ""}},
	}).All(&cs)
	if err != nil {
		return nil
	}
	//var chars []string
	//for idx := range cs {
	//	chars = append(chars, cs[idx].Character)
	//}
	if int32(len(cs)) <= 0 {
		return nil
	}
	log.Println("num:", int32(len(cs)))
	return cs[randomInt32(int32(len(cs)), f.calendar.Lunar().Time)]
}

func (f *fate) BestFirstTwo(c *mongo.Character) *mongo.Character {
	var cs []*mongo.Character

	if f.strokes == nil || c == nil {
		return nil
	}
	st, _ := strconv.Atoi(c.TotalStrokes)
	seconds := make(map[int][]byte)
	for idx := range f.strokes {
		log.Println(f.strokes [idx].FirstStroke[0] == st)
		if f.strokes [idx].FirstStroke[0] == st {
			log.Println(f.strokes[idx].FirstStroke[1])
			seconds[f.strokes[idx].FirstStroke[1]] = nil
		}
	}

	var charStroke []string
	for i := range seconds {
		charStroke = append(charStroke, strconv.Itoa(i))
	}

	err := mongo.C("character").Find(bson.M{
		"total_strokes":            bson.M{"$in": charStroke},
		"folk.commonly_characters": bson.M{"$in": []string{"是", ""}},
	}).All(&cs)
	if err != nil {
		return nil
	}
	//var chars []string
	//for idx := range cs {
	//	chars = append(chars, cs[idx].Character)
	//}
	if int32(len(cs)) <= 0 {
		return nil
	}
	log.Println("num:", int32(len(cs)))
	return cs[randomInt32(int32(len(cs)), f.calendar.Lunar().Time)]
}

func randomInt32(max int32, t time.Time) int32 {
	r := rand.NewSource(t.UnixNano())
	return rand.New(r).Int31n(max)
}
