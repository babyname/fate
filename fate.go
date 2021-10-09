package fate

import (
	"context"
	"fmt"
	"strings"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/model"

	"github.com/godcong/chronos"
)

// HandleOutputFunc ...
type HandleOutputFunc func(name Name)

// HelpContent ...
const HelpContent = "正在使用Fate生成姓名列表，如遇到问题请访问项目地址：https://github.com/babyname/fate获取帮助!"

// Fate ...
type Fate interface {
	Initialize(ctx context.Context) error
	LoadRules(rules ...RuleOption) Maker
	Debug() bool
	DB() *model.Model
}

type fateImpl struct {
	cfg      *config.Config
	db       *model.Model
	rule     *Rule
	out      Information
	born     chronos.Calendar
	last     []string
	lastChar []*Character
	names    []*Name
	nameType int
	sex      Sex
	debug    bool
	baZi     *BaZi
	zodiac   *Zodiac
	handle   HandleOutputFunc
}

func (f *fateImpl) Debug() bool {
	return f.cfg.Debug
}

func (f *fateImpl) LoadRules(rules ...RuleOption) Maker {
	rule := DefaultRule()
	for _, option := range rules {
		option(rule)
	}
	return newMaker(f, rule)
}

func (f *fateImpl) DB() *model.Model {
	return f.db
}

func (f *fateImpl) Initialize(ctx context.Context) error {
	panic("implement me")
}

func (f *fateImpl) Make(ctx context.Context, rules ...RuleOption) error {
	panic("implement me")
}

func InitDayanLuckyTable(ctx context.Context, model *model.Model) error {
	err := model.Schema.Create(ctx)
	if err != nil {
		return err
	}
	lucky := make(chan *WuGeLucky)
	go initWuGe(lucky)
	for la := range lucky {
		//todo
		fmt.Println("la", la)
	}
	return nil
}

// RegisterHandle ...
func (f *fateImpl) RegisterHandle(outputFunc HandleOutputFunc) {
	f.handle = outputFunc
}

//
//// MakeName ...
//func (f *fateImpl) MakeName(ctx context.Context) (e error) {
//	log.Info(HelpContent)
//	e = f.out.Head(f.config.FileOutput.Heads...)
//	if e != nil {
//		return Wrap(e, "write head failed")
//	}
//	e = f.RunInit()
//	if e != nil {
//		return Wrap(e, "init failed")
//	}
//	n, e := f.db.CountWuGeLucky()
//	if e != nil || n == 0 {
//		return Wrap(e, "count total error")
//	}
//
//	e = f.getLastCharacter()
//	if e != nil {
//		return Wrap(e, "get char failed")
//	}
//	name := make(chan *Name)
//	go func() {
//		e := f.getWugeName(name)
//		if e != nil {
//			log.Error(e)
//		}
//	}()
//
//	var tmpChar []*Character
//	//supplyFilter := false
//	for n := range name {
//		select {
//		case <-ctx.Done():
//			log.Info("end")
//			return
//		default:
//		}
//
//		tmpChar = n.FirstName
//		tmpChar = append(tmpChar, n.LastName...)
//		//filter bazi
//		if f.config.SupplyFilter && !filterXiYong(f.XiYong().Shen(), tmpChar...) {
//			//log.Infow("supply", "name", n.String())
//			continue
//		}
//		//filter zodiac
//		if f.config.ZodiacFilter && !filterZodiac(f.born, n.FirstName...) {
//			//log.Infow("zodiac", "name", n.String())
//			continue
//		}
//		//filter bagua
//		if f.config.BaguaFilter && !filterYao(n.BaGua(), "凶") {
//			//log.Infow("bagua", "name", n.String())
//			continue
//		}
//		ben := n.BaGua().Get(yi.BenGua)
//		bian := n.BaGua().Get(yi.BianGua)
//		if f.debug {
//			log.Infow("bazi", "born", f.born.LunarDate(), "time", f.born.Lunar().EightCharacter())
//			log.Infow("xiyong", "wuxing", n.WuXing(), "god", f.XiYong().Shen(), "pinheng", f.XiYong())
//			log.Infow("ben",
//				"ming",
//				ben.GuaMing,
//				"chu",
//				ben.ChuYaoJiXiong,
//				"er",
//				ben.ErYaoJiXiong,
//				"san",
//				ben.SanYaoJiXiong,
//				"si",
//				ben.SiYaoJiXiong,
//				"wu",
//				ben.WuYaoJiXiong,
//				"liu",
//				ben.ShangYaoJiXiong)
//			log.Infow("bian",
//				"ming",
//				bian.GuaMing,
//				"chu",
//				bian.ChuYaoJiXiong,
//				"er",
//				bian.ErYaoJiXiong,
//				"san",
//				bian.SanYaoJiXiong,
//				"si",
//				bian.SiYaoJiXiong,
//				"wu",
//				bian.WuYaoJiXiong,
//				"liu",
//				bian.ShangYaoJiXiong)
//		}
//
//		if err := f.out.Write(*n); err != nil {
//			return err
//		}
//		if f.debug {
//			log.Infow(n.String(),
//				"笔画",
//				n.Strokes(),
//				"拼音",
//				n.PinYin(),
//				"八字",
//				f.born.Lunar().EightCharacter(),
//				"喜用神",
//				f.XiYong().Shen(),
//				"本卦",
//				ben.GuaMing,
//				"变卦",
//				bian.GuaMing)
//		}
//	}
//	return nil
//}

//// XiYong ...
//func (f *fateImpl) XiYong() *XiYong {
//	if f.baZi == nil {
//		f.baZi = NewBazi(f.born)
//	}
//	return f.baZi.XiYong()
//}

//func (f *fateImpl) init() {
//	if f.config == nil {
//		f.config = config.DefaultConfig()
//	}
//
//	if f.config.FileOutput.Heads == nil {
//		f.config.FileOutput.Heads = config.DefaultHeads
//	}
//
//	f.db = initDatabaseWithConfig(f.config.Database)
//	f.out = initOutputWithConfig(f.config.FileOutput)
//}
//
//func (f *fateImpl) getWugeName(name chan<- *Name) (e error) {
//	defer func() {
//		close(name)
//	}()
//	lucky := make(chan *WuGeLucky)
//	go func() {
//		e = f.db.FilterWuGe(f.lastChar, lucky)
//		if e != nil {
//			log.Error(e)
//			return
//		}
//	}()
//	var f1s []*Character
//	var f2s []*Character
//	fsa := map[int][]*Character{}
//	bazi := NewBazi(f.born)
//	for l := range lucky {
//		if f.config.FilterMode == config.FilterModeCustom {
//			//TODO
//		}
//
//		if bool(f.sex) && filterSex(l) {
//			continue
//		}
//
//		if f.config.DayanFilter && hardFilter(l) {
//			sc := NewSanCai(l.TianGe, l.RenGe, l.DiGe)
//			if !Check(f.db.Database().(*xorm.Engine), sc, 5) {
//				continue
//			}
//		}
//
//		if f.config.StrokeMin > l.FirstStroke1 || f.config.StrokeMin > l.FirstStroke2 || f.config.StrokeMax < l.FirstStroke1 || f.config.StrokeMax < l.FirstStroke2 {
//			continue
//		}
//
//		if f.debug {
//			log.Infow("lucky", "l1", l.LastStroke1, "l2", l.LastStroke2, "f1", l.FirstStroke1, "f2", l.FirstStroke2)
//		}
//		if fsa[l.FirstStroke1] == nil {
//			if f.config.Regular {
//				f1s, e = f.db.GetCharacters(Stoker(l.FirstStroke1, Regular()))
//			} else {
//				f1s, e = f.db.GetCharacters(Stoker(l.FirstStroke1))
//			}
//
//			if e != nil {
//				return Wrap(e, "first stroke1 error")
//			}
//
//			fsa[l.FirstStroke1] = f1s
//		} else {
//			f1s = fsa[l.FirstStroke1]
//		}
//
//		if fsa[l.FirstStroke2] == nil {
//			if f.config.Regular {
//				f2s, e = f.db.GetCharacters(Stoker(l.FirstStroke2, Regular()))
//			} else {
//				f2s, e = f.db.GetCharacters(Stoker(l.FirstStroke2))
//			}
//
//			if e != nil {
//				return Wrap(e, "first stoke2 error")
//			}
//
//			fsa[l.FirstStroke2] = f2s
//		} else {
//			f2s = fsa[l.FirstStroke2]
//		}
//
//		for _, f1 := range f1s {
//			if len(f1.PinYin) == 0 {
//				continue
//			}
//			for _, f2 := range f2s {
//				if len(f2.PinYin) == 0 {
//					continue
//				}
//				n := createName(f, f1, f2)
//				n.baZi = bazi
//				name <- n
//			}
//		}
//	}
//	return nil
//}

func filterSex(lucky *WuGeLucky) bool {
	return lucky.ZongSex == true
}

func isLucky(s string) bool {
	if strings.Compare(s, "吉") == 0 || strings.Compare(s, "半吉") == 0 {
		return true
	}
	return false
}

func hardFilter(lucky *WuGeLucky) bool {
	if !isLucky(GetDaYan(lucky.DiGe).Lucky) ||
		!isLucky(GetDaYan(lucky.RenGe).Lucky) ||
		!isLucky(GetDaYan(lucky.WaiGe).Lucky) ||
		!isLucky(GetDaYan(lucky.ZongGe).Lucky) {
		return true
	}
	return false
}

func New(cfg *config.Config) Fate {
	return &fateImpl{
		cfg: cfg,
	}
}

var _ Fate = (*fateImpl)(nil)
