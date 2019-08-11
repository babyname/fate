package fate

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/godcong/chronos"
	"github.com/godcong/fate/mongo"
	"github.com/godcong/go-trait"
	"github.com/godcong/yi"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"strings"
	"time"
)

var log = trait.NewZapSugar()

type Fate interface {
	FirstRunInit()
	SetDB(engine *xorm.Engine)
}

type fate struct {
	db       *xorm.Engine
	born     chronos.Calendar
	last     []string
	names    []*Name
	nameType int

	sex      string
	firstOne []*mongo.Character
}

type Generating struct {
	martial *Martial
	current interface{} //当前对象
	step    int         //当前
	number  int         //生成数
	fate    *fate
	//stroke    []*Stroke
	character []*mongo.Character
}

//NewFate 所有的入口,新建一个fate对象
func NewFate(lastName string, born time.Time) Fate {
	f := &fate{
		last:     strings.Split(lastName, ""),
		born:     chronos.New(born),
		nameType: mongo.KangXi,
	}
	f.preInit()

	return f
}

func (f *fate) SetDB(engine *xorm.Engine) {
	f.db = engine
}

func (f *fate) MakeName() {

}

func (f *fate) FirstRunInit() {
	e := f.db.Sync2(WuGeLucky{})
	if e != nil {
		return
	}
	ge := initWuGe()
init:
	for {
		select {
		case lu := <-ge:
			if lu == nil {
				break init
			}
			_, e = InsertOrUpdate(f.db.NewSession(), lu)
			if e != nil {
				panic(e)
			}
		}
	}
}

func (f *fate) preInit() {
	var e error
	if f.db == nil {
		f.db, e = InitSQLite3("fate.db")
		if e != nil {
			panic(e)
		}
	}
}

//SetBornData 设定生日
func (f *fate) SetBornData(t time.Time) {
	f.born = chronos.New(t)
}

func randomInt32(max uint32, t time.Time) uint32 {
	r := rand.NewSource(t.UnixNano())
	return rand.New(r).Uint32()
}

func (g *Generating) Character() []*mongo.Character {
	return nil
}

func filterGuaXiang(characters []*mongo.Character) []*mongo.Character {
	gua := yi.NumberQiGua(0, 0, 0)
	log.Info(gua)
	return nil
}

func filterTianYun() []*mongo.Character {
	//TODO:
	return nil
}

func filterBaZi(character []*mongo.Character, wuxing []string) []*mongo.Character {
	//TODO:
	//计算平衡用神
	return nil
}

// sqlite3DB ...
func sqlite3DB(name string) string {
	return fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", name)
}

// InitSQLite3 ...
func InitSQLite3(name string) (engine *xorm.Engine, e error) {
	engine, e = xorm.NewEngine("sqlite3", name)
	if e != nil {
		return nil, e
	}
	return engine, nil
}
