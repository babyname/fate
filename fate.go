package fate

import (
	"log"
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
	firstOne []*mongo.Character
	calendar chronos.Calendar
}

type Generating struct {
	martial   *Martial
	current   interface{} //当前对象
	step      int         //当前
	number    int         //生成数
	fate      *fate
	stroke    []*Stroke
	character []*mongo.Character
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

func filterWuGe(f *fate) []*Stroke {
	var rltS []*Stroke
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
			rltS = append(rltS, &Stroke{
				LastStroke:  []int{l1, l2},
				FirstStroke: []int{f1, f2},
				wuge:        wuge,
				sancai:      nil,
			})
		}

		if f2 >= 30 {
			f2 = 0
			f1++
		}
	}

	return rltS
}

func filterSanCai(s []*Stroke) []*Stroke {
	var strokes []*Stroke
	var wx mongo.WuXing
	if s == nil {
		return nil
	}
	for idx := range s {
		sc := MakeSanCai(s[idx].wuge)
		mongo.C("wuxing").Find(bson.M{
			"wu_xing": []string{sc.TianCai, sc.RenCai, sc.DiCai},
		}).One(&wx)
		switch wx.Fortune {
		case "吉", "中吉", "大吉", "吉多于凶":
			strokes = append(strokes, s[idx])
		}
	}
	return strokes
}

func (g *Generating) CurrentStep() int {
	return g.step
}

func (g *Generating) SetMartial(martial *Martial) {
	g.martial = martial
}

func (g *Generating) GetMartial() *Martial {
	if g.martial == nil {
		return &Martial{}
	}
	return g.martial
}

func (g *Generating) Strokes() []*Stroke {
	return g.stroke
}

func (g *Generating) Character() []*mongo.Character {
	return nil
}

func (g *Generating) PreStroke() *Generating {
	//过滤五格
	if g.step == 0 {
		if g.martial.BiHua {
			g.stroke = filterWuGe(g.fate)
		}
	}

	//过滤三才
	if g.step == 1 {
		if g.martial.SanCai {
			g.stroke = filterSanCai(g.stroke)
		}
	}
	return g
}

func (g *Generating) Continue() *Generating {

	//过滤生肖
	if g.martial.ShengXiao {
		g.character = filterShengXiao(g.stroke)
	}

	log.Printf("stroke %+v", g.stroke)

	//过滤八字
	if g.martial.BaZi {

	}

	//过滤天运
	if g.martial.TianYun {

	}

	//过滤卦象
	if g.martial.GuaXiang {

	}
	g.step++
	return g
}

func filterShengXiao(strokes []*Stroke) []*mongo.Character {
	return nil
}
