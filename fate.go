package fate

import (
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

const DefaultDatabase = "fate.db"

var log = trait.NewZapSugar()

type Fate interface {
	MakeName() (e error)
	SetDB(engine *xorm.Engine)
	SetCharDB(engine *xorm.Engine)
	//GetLastCharacter() error
}

type fate struct {
	chardb   *xorm.Engine
	db       *xorm.Engine
	born     chronos.Calendar
	last     []string
	lastChar []*Character
	names    []*Name
	nameType int
	sex      string
	isFirst  bool
	Limit    int
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

type Options func(f *fate)

//NewFate 所有的入口,新建一个fate对象
func NewFate(lastName string, born time.Time, options ...Options) Fate {
	f := &fate{
		last:     strings.Split(lastName, ""),
		born:     chronos.New(born),
		nameType: mongo.KangXi,
	}
	f.lastChar = make([]*Character, len(f.last))
	for _, op := range options {
		op(f)
	}

	f.init()

	return f
}

func Database(engine *xorm.Engine) Options {
	return func(f *fate) {
		f.db = engine
	}
}

func CharacterDatabase(engine *xorm.Engine) Options {
	return func(f *fate) {
		f.chardb = engine
	}
}

func (f *fate) SetCharDB(engine *xorm.Engine) {
	f.chardb = engine
}

func (f *fate) SetDB(engine *xorm.Engine) {
	f.db = engine
}

func (f *fate) RandomName() {
	//filterWuGe(f.db, f.last...)
}

func (f *fate) getLastCharacter() error {
	for i, c := range f.last {
		character, e := getCharacter(f, Char(c))
		if e != nil {
			return e
		}
		log.With("index", i, "char", c).Info("last")
		f.lastChar[i] = character
	}
	return nil
}

func (f *fate) MakeName() (e error) {
	e = f.db.Sync2(WuGeLucky{})
	if e != nil {
		return e
	}
	if f.chardb == nil {
		f.chardb = f.db
	}

	n, e := CountWuGeLucky(f.db)
	if e != nil {
		return e
	}
	f.isFirst = n == 0
	if f.isFirst {
		lucky := make(chan *WuGeLucky)
		go initWuGe(lucky)
		for la := range lucky {
			_, e = InsertOrUpdateWuGeLucky(f.db, la)
			if e != nil {
				return e
			}
		}
	}

	return f.getLastCharacter()
}

func (f *fate) init() {
	var e error
	if f.db == nil {
		f.db, e = NewSQLite3(DefaultDatabase)
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
