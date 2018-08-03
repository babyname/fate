package fate

import (
	"math/rand"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/godcong/chronos"
	"github.com/godcong/fate/config"
	"github.com/godcong/fate/mongo"
)

type fate struct {
	nameType int
	name     *Name
	sex      string
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
	wuge    []*WuGe
	sancai  []*SanCai
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
	return &fate{
		nameType: mongo.KangXi,
		name:     name,
	}
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

func filterWuGe(f *fate) []*WuGe {
	var rltWuge []*WuGe
	l1 := f.name.lastChar[0].GetStrokeByType(f.nameType)
	l2 := 0
	if len(f.name.firstChar) == 2 {
		l2 = f.name.lastChar[1].GetStrokeByType(f.nameType)
	}

	var dy []*mongo.DaYan
	mongo.C("dayan").Find(nil).Sort("index").All(&dy)
	for f1, f2 := 1, 1; 30 >= f1; f2++ {
		wuge := MakeWuGe(l1, l2, f1, f2)

		zg := checkWuGe(dy, wuge.ZongGe)
		wg := checkWuGe(dy, wuge.WaiGe)
		rg := checkWuGe(dy, wuge.RenGe)
		dg := checkWuGe(dy, wuge.DiGe)
		if zg && wg && rg && dg {
			rltWuge = append(rltWuge, wuge)
		}

		if f2 >= 30 {
			f2 = 0
			f1++
		}
	}

	return rltWuge
}

func filterSanCai(wuge []*WuGe) []*SanCai {
	var scs []*mongo.WuXing
	if wuge == nil {
		return nil
	}
	for idx := range wuge {
		sc := MakeSanCai(wuge[idx])
		mongo.C("wuxing").Find(bson.M{
			"wu_xing": []string{sc.TianCai, sc.RenCai, sc.DiCai},
		}).One(&scs)
	}
	return nil
}

func (g *Generating) CurrentStep() int {
	return g.step
}

func (g *Generating) Continue() *Generating {
	f := g.fate
	//过滤五格
	if g.step == 0 {
		if f.martial.BiHua {
			g.wuge = filterWuGe(g.fate)
		}
	}

	//过滤三才
	if g.step == 1 {
		if f.martial.SanCai {
			g.sancai = filterSanCai(g.wuge)
		}
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
	g.step++
	return g
}
