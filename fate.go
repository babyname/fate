package fate

import (
	"math/rand"
	"strings"
	"time"

	"github.com/godcong/chronos"
	"github.com/godcong/fate/mongo"
	"github.com/godcong/yi"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xormsharp/xorm"
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
		return Wrap(e, "get char failed")
	}
	name := make(chan *Name)
	go func() {
		e := f.getCharacterWugeLucky(name)
		if e != nil {
			log.Error(e)
		}
	}()

	var tmpChar []*Character
	for n := range name {
		tmpChar = n.FirstName
		tmpChar = append(tmpChar, n.LastName...)
		if filterXiYong(f.XiYong().Shen(), tmpChar...) {
			log.With("wuxing", n.WuXing(), "xi", f.XiYong()).Info(n)
		}
	}
	return nil
}

func (f *fateImpl) XiYong() *XiYong {
	if f.baZi == nil {
		f.baZi = NewBazi(f.born)
	}
	return f.baZi.XiYong()
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

func (f *fateImpl) getCharacterWugeLucky(name chan<- *Name) (e error) {
	defer func() {
		close(name)
	}()
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
		log.With("l1", l.LastStroke1, "l2", l.LastStroke2, "f1", l.FirstStroke1, "f2", l.FirstStroke2).Info("lucky")
		f1s, e = getCharacters(f, Stoker(l.FirstStroke1))
		if e != nil {
			return Wrap(e, "first stroke1 error")
		}
		f2s, e = getCharacters(f, Stoker(l.FirstStroke2))
		if e != nil {
			return Wrap(e, "first stoke2 error")
		}

		for _, f1 := range f1s {
			for _, f2 := range f2s {
				n := createName(f, f1, f2)
				name <- n
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
