package fate

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
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

type Generating struct {
	current interface{} //当前对象
	step    int         //当前
	number  int         //生成数
	fate    *fate
	wuge    *WuGe
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

//MaxStokers 超过30划的字不易书写,过滤
const MaxStokers = 30

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

func randomInt32(max uint32, t time.Time) uint32 {
	r := rand.NewSource(t.UnixNano())
	return rand.New(r).Uint32()
}

func (f *fate) Generate(number int) *Generating {
	g := &Generating{
		step:   0,
		number: number,
		fate:   f,
	}
	return g

}

func filterWuGe(wg *WuGe) bool {
	return false
}

func (g *Generating) Continue() *Generating {
	f := g.fate
	if g.step == 0 && f.martial.BiHua {
		fc1, _ := strconv.Atoi(f.name.firstChar[0].KangxiStrokes)
		fc2 := 0
		if len(f.name.firstChar) == 2 {
			fc2, _ = strconv.Atoi(f.name.firstChar[1].KangxiStrokes)
		}
		MakeWuGe(fc1, fc2, 1, 1)
	}
	var names []*Name
	//获取笔画列表
	//var strokes []*Stroke

	//过滤五格
	if f.martial.BiHua {

	}
	//过滤三才
	if f.martial.SanCai {

	}

	//过滤生肖
	if f.martial.ShengXiao {

	}

	//过滤八字
	if f.martial.BaZi {

	}

	//过滤天运
	if f.martial.TianYun {

	}

	//过滤卦象
	if f.martial.GuaXiang {

	}

	return names
}
