package fate

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	_ "time/tzdata"

	xt "github.com/free-utils-go/xorm_type_assist"
	"github.com/godcong/chronos"
	"github.com/godcong/fate/config"
	"github.com/godcong/yi"
	"github.com/goextension/log"
)

// HandleOutputFunc ...
type HandleOutputFunc func(name Name)

// HelpContent ...
const HelpContent = "正在使用Fate生成姓名列表，如遇到问题请访问项目地址：https://github.com/godcong/fate获取帮助!"

// Fate ...
type Fate interface {
	MakeName(ctx context.Context) (e error)
	XiYong() *XiYong
	RunInit() (e error)
	RegisterHandle(outputFunc HandleOutputFunc)
}

type fateImpl struct {
	config    *config.Config
	db        Database
	out       Information
	born      chronos.Calendar
	last      []string
	lastChar  []*Character
	names     []*Name
	nameType  int
	sex       yi.Sex
	debug     bool
	baZi      *BaZi
	xiYongStr string
	zodiac    *Zodiac
	handle    HandleOutputFunc
}

// RunInit ...
func (f *fateImpl) RunInit() (e error) {
	if f.config.RunInit {
		return nil
	}
	return nil
}

// Options ...
type Options func(f *fateImpl)

// ConfigOption ...
func ConfigOption(cfg *config.Config) Options {
	return func(f *fateImpl) {
		f.config = cfg
	}
}

// SexOption ...
func SexOption(sex yi.Sex) Options {
	return func(f *fateImpl) {
		f.sex = sex
	}
}

func XiYongOption(xiyong string) Options {
	return func(f *fateImpl) {
		f.xiYongStr = xiyong
	}
}

// Debug ...
func Debug() Options {
	return func(f *fateImpl) {
		f.debug = true
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
		panic("last name was bigger than 2 characters")
	}

	for _, op := range options {
		op(f)
	}

	f.init()

	return f
}

// RegisterHandle ...
func (f *fateImpl) RegisterHandle(outputFunc HandleOutputFunc) {
	f.handle = outputFunc
}

//
func (f *fateImpl) getLastCharacter() error {
	size := len(f.last)
	if size == 0 {
		return errors.New("last name was not inputted")
	} else if size > 2 {
		return fmt.Errorf("%d characters last name was not supported", size)
	}

	for i, c := range f.last {
		character, e := f.db.GetCharacter(Char(c))
		if e != nil {
			return e
		}
		f.lastChar[i] = character
	}

	f.setStrokeCache()
	return nil
}

// MakeName ...
func (f *fateImpl) MakeName(ctx context.Context) (e error) {
	log.Info(HelpContent)
	e = f.out.Head(f.config.FileOutput.Heads...)
	if e != nil {
		return Wrap(e, "write head failed")
	}
	e = f.RunInit()
	if e != nil {
		return Wrap(e, "init failed")
	}

	e = f.getLastCharacter()
	if e != nil {
		return Wrap(e, "get char failed")
	}

	n := CountNameStrokesLucky()
	if n == 0 {
		return Wrap(e, "count total error")
	}

	name := make(chan *Name)
	go func() {
		e := f.getLuckyName(name)
		if e != nil {
			log.Error(e)
		}
	}()

	var tmpChar []*Character

	for n := range name {
		select {
		case <-ctx.Done():
			log.Info("end")
			return
		default:
		}

		tmpChar = n.FirstName
		tmpChar = append(tmpChar, n.LastName...)
		//filter bazi
		//if f.config.SupplyFilter && !filterXiYong(f.XiYong().Shen(), tmpChar...) {
		if f.config.SupplyFilter && !fixedFilterXiYong(f.xiYongStr, f.config.HardFilter, tmpChar...) {
			//log.Infow("supply", "name", n.String())
			continue
		}
		//filter zodiac
		if f.config.ZodiacFilter && !filterZodiac(f.born, n.FirstName...) {
			//log.Infow("zodiac", "name", n.String())
			continue
		}
		//filter bagua
		if f.config.BaguaFilter && !n.IsLucky(f.sex, true, f.config.HardFilter) {
			//log.Infow("bagua", "name", n.String())
			continue
		}

		if !n.IsLucky(f.sex, f.config.BaguaFilter, f.config.HardFilter) {
			panic("bad name")
		}

		ben := n.nameScienceStroke.BaGua().Get(yi.BenGua)
		bian := n.nameScienceStroke.BaGua().Get(yi.BianGua)

		if f.debug {
			log.Infow("bazi", "born", f.born.LunarDate(), "time", f.born.Lunar().EightCharacter())
			log.Infow("xiyong", "wuxing", n.WuXing(), "god", f.XiYong().Shen(), "pinheng", f.XiYong())
			log.Infow("ben", "ming", ben.GuaMing, "chu", ben.GuaYaos[0].JiXiong, "er", ben.GuaYaos[1].JiXiong, "san", ben.GuaYaos[2].JiXiong, "si", ben.GuaYaos[3].JiXiong, "wu", ben.GuaYaos[4].JiXiong, "liu", ben.GuaYaos[5].JiXiong)
			log.Infow("bian", "ming", bian.GuaMing, "chu", bian.GuaYaos[0].JiXiong, "er", bian.GuaYaos[1].JiXiong, "san", bian.GuaYaos[2].JiXiong, "si", bian.GuaYaos[3].JiXiong, "wu", bian.GuaYaos[4].JiXiong, "liu", bian.GuaYaos[5].JiXiong)
		}

		if err := f.out.Write(*n); err != nil {
			return err
		}
		if f.debug {
			log.Infow(n.String(), "笔画", n.Strokes(), "拼音", n.PinYin(), "八字", f.born.Lunar().EightCharacter(), "喜用神", f.XiYong().Shen(), "本卦", ben.GuaMing, "变卦", bian.GuaMing)
		}
	}
	return nil
}

// XiYong ...
func (f *fateImpl) XiYong() *XiYong {
	if f.baZi == nil {
		f.baZi = NewBazi(f.born)
	}
	return f.baZi.XiYong()
}

func (f *fateImpl) init() {
	if f.config == nil {
		f.config = config.LoadConfig()
	}

	if f.config.FileOutput.Heads == nil {
		f.config.FileOutput.Heads = config.DefaultHeads
	}

	f.db = InitDatabaseWithConfig(*f.config)
	f.out = initOutputWithConfig(f.config.FileOutput)
}

//SetBornData 设定生日
func (f *fateImpl) SetBornData(t time.Time) {
	f.born = chronos.New(t)
}

func (f *fateImpl) getLuckyName(name chan<- *Name) (e error) {
	defer func() {
		close(name)
	}()
	lucky := make(chan *NameStroke)
	go func() {
		e = f.FilterNameStrokes(lucky)
		if e != nil {
			log.Error(e)
			return
		}
	}()
	var f1s []*Character
	var f2s []*Character
	fsa := map[int]map[string]*Character{}
	bazi := NewBazi(f.born)
	for l := range lucky {
		f1s_u := map[string]*Character{}
		f2s_u := map[string]*Character{}
		if f.config.FilterMode == config.FilterModeCustom {
			//TODO
		}

		if f.debug {
			log.Infow("lucky", "l1", l.LastStroke1, "l2", l.LastStroke2, "f1", l.FirstStroke1, "f2", l.FirstStroke2)
		}
		if fsa[l.FirstStroke1] == nil {
			if f.config.Regular {
				f1s, e = f.db.GetCharacters(StokerX(l.FirstStroke1, Regular()))
			} else {
				f1s, e = f.db.GetCharacters(StokerX(l.FirstStroke1))
			}

			if e != nil {
				return Wrap(e, "first stroke1 error")
			}

			for _, f1 := range f1s {
				f1s_u[f1.Ch] = f1
			}

			fsa[l.FirstStroke1] = f1s_u
		} else {
			f1s_u = fsa[l.FirstStroke1]
		}

		if fsa[l.FirstStroke2] == nil {
			if f.config.Regular {
				f2s, e = f.db.GetCharacters(StokerX(l.FirstStroke2, Regular()))
			} else {
				f2s, e = f.db.GetCharacters(StokerX(l.FirstStroke2))
			}

			if e != nil {
				return Wrap(e, "first stoke2 error")
			}

			for _, f2 := range f2s {
				f2s_u[f2.Ch] = f2
			}

			fsa[l.FirstStroke2] = f2s_u
		} else {
			f2s_u = fsa[l.FirstStroke2]
		}

		for _, f1 := range f1s_u {
			if len(f1.PinYin) == 0 {
				continue
			}
			if f1.getStrokeScience() != l.FirstStroke1 {
				fmt.Println(f1)
				panic(fmt.Sprintf("%d,%d", l.FirstStroke1, f1.getStrokeScience()))
			}
			for _, f2 := range f2s_u {
				if len(f2.PinYin) == 0 {
					continue
				}

				n := f.createName([]*Character{f1, f2}, *l)

				if !n.IsLucky(f.sex, f.config.BaguaFilter, f.config.HardFilter) {
					continue
				}

				if f.config.StrokeMin > n.nameStroke.FirstStroke1 || f.config.StrokeMax < n.nameStroke.FirstStroke1 {
					continue
				} else if (n.nameStroke.FirstStroke2 != 0 && f.config.StrokeMin > n.nameStroke.FirstStroke2) || (n.nameStroke.FirstStroke2 != 0 && f.config.StrokeMax < n.nameStroke.FirstStroke2) {
					continue
				}

				if f.config.HardFilter {
					if (n.nameStroke.FirstStroke2 != 0 && n.FirstName[1].IsDuoYin == xt.TRUE) || n.FirstName[0].IsDuoYin == xt.TRUE {
						continue
					}
				}

				n.baZi = bazi
				name <- n
			}
		}
	}
	return nil
}
