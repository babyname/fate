package fate

import (
	"strings"
	"time"

	"github.com/godcong/chronos"
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

	isFirst      bool
	Limit        int
	baZi         *BaZi
	zodiac       *Zodiac
	supplyFilter bool //补八字
	zodiacFilter bool //生肖
	baguaFilter  bool //卦象
}

type Options func(f *fateImpl)

func BaGuaFilter() Options {
	return func(f *fateImpl) {
		f.baguaFilter = true
	}
}

func ZodiacFilter() Options {
	return func(f *fateImpl) {
		f.zodiacFilter = true
	}
}

func SupplyFilter() Options {
	return func(f *fateImpl) {
		f.supplyFilter = true
	}
}

//NewFate 所有的入口,新建一个fate对象
func NewFate(lastName string, born time.Time, options ...Options) Fate {
	f := &fateImpl{
		last: strings.Split(lastName, ""),
		born: chronos.New(born),
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
	//supplyFilter := false
	for n := range name {
		tmpChar = n.FirstName
		tmpChar = append(tmpChar, n.LastName...)
		//filter bazi
		if f.supplyFilter && !filterXiYong(f.XiYong().Shen(), tmpChar...) {
			continue
		}
		//filter zodiac
		if f.zodiacFilter && !filterZodiac(f.born, n.FirstName...) {
			continue
		}

		if f.baguaFilter && !filterYao(n.baGua, "凶", "平") {
			continue
		}
		log.With("name", n.String()).Info("name")
		log.With("born", f.born.LunarDate(), "time", f.born.Lunar().EightCharacter()).Info("bazi")
		log.With("wuxing", n.WuXing(), "god", f.XiYong().Shen(), "pinheng", f.XiYong()).Info("xiyong")
		ben := n.BaGua().Get(yi.BenGua)
		log.With("ming", ben.GuaMing, "chu", ben.ChuYaoJiXiong, "er", ben.ErYaoJiXiong, "san", ben.SanYaoJiXiong, "si", ben.SiYaoJiXiong, "wu", ben.WuYaoJiXiong, "liu", ben.ShangYaoJiXiong).Info("ben")
		bian := n.BaGua().Get(yi.BianGua)
		log.With("ming", bian.GuaMing, "chu", bian.ChuYaoJiXiong, "er", bian.ErYaoJiXiong, "san", bian.SanYaoJiXiong, "si", bian.SiYaoJiXiong, "wu", bian.WuYaoJiXiong, "liu", bian.ShangYaoJiXiong).Info("bian")
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
