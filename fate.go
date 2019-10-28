package fate

import (
	"math/rand"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/godcong/chronos"
	"github.com/godcong/fate/mongo"
	"github.com/godcong/yi"
	_ "github.com/mattn/go-sqlite3"
)

var DefaultDatabase = "fate.db"

type Fate interface {
	MakeName() (e error)
	XiYong() *XiYong
	//SetCharDB(engine *xorm.Engine)
	//GetLastCharacter() error
}

type fateImpl struct {
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
	baZi     *BaZi
	zodiac   *Zodiac
}

type Generating struct {
	martial *Martial
	current interface{} //当前对象
	step    int         //当前
	number  int         //生成数
	fate    *fateImpl
	//stroke    []*Stroke
	character []*mongo.Character
}

type Options func(f *fateImpl)

//NewFate 所有的入口,新建一个fate对象
func NewFate(lastName string, born time.Time, options ...Options) Fate {
	f := &fateImpl{
		last:     strings.Split(lastName, ""),
		born:     chronos.New(born),
		nameType: mongo.KangXi,
	}
	f.lastChar = make([]*Character, len(f.last))
	if len(f.last) > 2 {
		panic("last name could not bigger than 2 characters")
	}

	for _, op := range options {
		op(f)
	}

	f.init()

	return f
}

func Database(engine *xorm.Engine) Options {
	return func(f *fateImpl) {
		f.db = engine
	}
}

func CharacterDatabase(engine *xorm.Engine) Options {
	return func(f *fateImpl) {
		f.chardb = engine
	}
}

//func (f *fateImpl) SetCharDB(engine *xorm.Engine) {
//	f.chardb = engine
//}
//
//func (f *fateImpl) SetDB(engine *xorm.Engine) {
//	f.db = engine
//}

func (f *fateImpl) RandomName() {
	//filterWuGe(f.db, f.last...)
}

func (f *fateImpl) getLastCharacter() error {
	for i, c := range f.last {
		character, e := getCharacter(f, Char(c))
		if e != nil {
			return e
		}
		log.With("index", i, "char", c, "character", character).Info("last name")
		f.lastChar[i] = character
	}
	return nil
}

func (f *fateImpl) MakeName() (e error) {
	n, e := CountWuGeLucky(f.db)
	if e != nil {
		return Wrap(e, "count total error")
	}
	f.isFirst = n == 0
	if f.isFirst {
		lucky := make(chan *WuGeLucky)
		go initWuGe(lucky)
		for la := range lucky {
			_, e = InsertOrUpdateWuGeLucky(f.db, la)
			if e != nil {
				return Wrap(e, "insert failed")
			}
		}
	}

	e = f.getLastCharacter()
	if e != nil {
		return Wrap(e, "get char faile")
	}
	return f.getCharacterWugeLucky()

}

func (f *fateImpl) XiYong() *XiYong {
	if f.baZi == nil {
		b := NewBazi(f.born)
		b.XiYong()
		f.baZi = b
	}
	return f.baZi.xiyong
}

func (f *fateImpl) init() {
	var e error
	if f.db == nil {
		f.db, e = NewSQLite3(DefaultDatabase)
		if e != nil {
			panic(e)
		}
	}

	//use the same db when char db not set
	if f.chardb == nil {
		f.chardb = f.db
	}

	e = f.db.Sync2(WuGeLucky{})
	if e != nil {
		panic(e)
	}
}

//SetBornData 设定生日
func (f *fateImpl) SetBornData(t time.Time) {
	f.born = chronos.New(t)
}

func (f *fateImpl) getCharacterWugeLucky() (e error) {
	lucky := make(chan *WuGeLucky)
	go func() {
		e = filterWuGe(lucky, f)
		if e != nil {
			log.Error(e)
			return
		}
	}()
	var f1s []*Character
	var f2s []*Character
	for l := range lucky {
		log.With("first1", l.FirstStroke1, "first2", l.FirstStroke2).Info("lucky")
		f1s, e = getCharacters(f, Stoker(l.FirstStroke1))
		if e != nil {
			return e
		}
		f2s, e = getCharacters(f, Stoker(l.FirstStroke2))
		if e != nil {
			return e
		}

		for _, f1 := range f1s {
			for _, f2 := range f2s {
				name := createName(f, f1, f2)
				log.Info(name)
			}
		}
	}
	return nil
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
