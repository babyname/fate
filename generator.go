package fate

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/godcong/chronos"
	//"github.com/godcong/chronos"
	"github.com/goextension/log"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/model"
)

var ErrLastNameNotInput = errors.New("last name was not inputted")

type Generator interface {
	Run(ctx context.Context, lastName string, born time.Time) error
	Wait() <-chan Name
}

type generateImpl struct {
	Fate
	strokeChars struct {
		sync.RWMutex
		c map[int][]*ent.Character
	}
	last []string
	born *Calendar
	rule *Rule
}

func (i *generateImpl) Run(ctx context.Context, lastName string, born time.Time) error {
	i.born = NewCalendar(born)
	return nil
}

func (i *generateImpl) Wait() <-chan Name {
	//TODO implement me
	panic("implement me")
}

func newGenerator(fate Fate, rule *Rule) Generator {
	return &generateImpl{
		Fate: fate,
		rule: rule,
	}
}

func (i *generateImpl) getLastCharacter(ctx context.Context, lastName string) ([]*ent.Character, error) {
	last := strings.Split(lastName, "")
	size := len(last)
	switch size {
	case 0:
		return nil, ErrLastNameNotInput
	case 1, 2:
		break
	default:
		return nil, fmt.Errorf("last name over(%d) was not supported", size)
	}

	lastChar := make([]*ent.Character, len(last))
	for idx, c := range i.last {
		character, e := i.Query().GetCharacter(ctx, model.Char(c))
		if e != nil {
			return nil, e
		}
		lastChar[idx] = character
	}
	return lastChar, nil
}

func (i *generateImpl) Generate(ctx context.Context, lastName string, born time.Time) error {
	last, err := i.getLastCharacter(ctx, lastName)
	if err != nil {
		return err
	}
	if i.Debug() {
		log.Debugw("get last character", "total", len(last), "last", last)
	}

	bornTime := chronos.New(born)
	if i.Debug() {
		log.Debugw("get born date", "time", bornTime.LunarDate())
	}

	log.Info(HelpContent)

	//name := generate(chan *Name)
	//go func() {
	//	e := f.getWugeName(name)
	//	if e != nil {
	//		log.Error(e)
	//	}
	//}()
	//
	//var tmpChar []*Character
	////supplyFilter := false
	//for n := range name {
	//	select {
	//	case <-ctx.Done():
	//		log.Info("end")
	//		return nil
	//	default:
	//	}
	//
	//	tmpChar = n.FirstName
	//	tmpChar = append(tmpChar, n.LastName...)
	//	//filter bazi
	//	if f.config.SupplyFilter && !filterXiYong(f.XiYong().Shen(), tmpChar...) {
	//		//log.Infow("supply", "name", n.String())
	//		continue
	//	}
	//	//filter zodiac
	//	if f.config.ZodiacFilter && !filterZodiac(f.born, n.FirstName...) {
	//		//log.Infow("zodiac", "name", n.String())
	//		continue
	//	}
	//	//filter bagua
	//	if f.config.BaguaFilter && !filterYao(n.BaGua(), "凶") {
	//		//log.Infow("bagua", "name", n.String())
	//		continue
	//	}
	//	ben := n.BaGua().Get(yi.BenGua)
	//	bian := n.BaGua().Get(yi.BianGua)
	//	if f.debug {
	//		log.Infow("bazi", "born", f.born.LunarDate(), "time", f.born.Lunar().EightCharacter())
	//		log.Infow("xiyong", "wuxing", n.WuXing(), "god", f.XiYong().Shen(), "pinheng", f.XiYong())
	//		log.Infow("ben",
	//			"ming",
	//			ben.GuaMing,
	//			"chu",
	//			ben.ChuYaoJiXiong,
	//			"er",
	//			ben.ErYaoJiXiong,
	//			"san",
	//			ben.SanYaoJiXiong,
	//			"si",
	//			ben.SiYaoJiXiong,
	//			"wu",
	//			ben.WuYaoJiXiong,
	//			"liu",
	//			ben.ShangYaoJiXiong)
	//		log.Infow("bian",
	//			"ming",
	//			bian.GuaMing,
	//			"chu",
	//			bian.ChuYaoJiXiong,
	//			"er",
	//			bian.ErYaoJiXiong,
	//			"san",
	//			bian.SanYaoJiXiong,
	//			"si",
	//			bian.SiYaoJiXiong,
	//			"wu",
	//			bian.WuYaoJiXiong,
	//			"liu",
	//			bian.ShangYaoJiXiong)
	//	}
	//
	//	if err := f.out.Write(*n); err != nil {
	//		return err
	//	}
	//	if f.debug {
	//		log.Infow(n.String(),
	//			"笔画",
	//			n.Strokes(),
	//			"拼音",
	//			n.PinYin(),
	//			"八字",
	//			f.born.Lunar().EightCharacter(),
	//			"喜用神",
	//			f.XiYong().Shen(),
	//			"本卦",
	//			ben.GuaMing,
	//			"变卦",
	//			bian.GuaMing)
	//	}
	//}
	return nil
}

var _ generater = (*generateImpl)(nil)
