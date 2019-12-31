package main

import (
	"encoding/json"
	"github.com/godcong/fate"
	"github.com/goextension/log"
	"github.com/xormsharp/xorm"
	"io/ioutil"
	"strconv"
)

type Dict struct {
	Jin  map[string][]string `json:"jin"`
	Mu   map[string][]string `json:"mu"`
	Huo  map[string][]string `json:"huo"`
	Shui map[string][]string `json:"shui"`
	Tu   map[string][]string `json:"tu"`
}

var dict Dict

type WuXingFunc func(s string) bool

func CheckLoader(s string) error {
	bytes, e := ioutil.ReadFile(s)
	if e != nil {
		return e
	}
	e = json.Unmarshal(bytes, &dict)
	if e != nil {
		return e
	}

	return nil
}

func CheckVerify(db *xorm.Engine) error {
	if err := verifySub(db, dict.Jin, "金"); err != nil {
		return err
	}
	if err := verifySub(db, dict.Mu, "木"); err != nil {
		return err
	}
	if err := verifySub(db, dict.Shui, "水"); err != nil {
		return err
	}
	if err := verifySub(db, dict.Huo, "火"); err != nil {
		return err
	}
	if err := verifySub(db, dict.Tu, "土"); err != nil {
		return err
	}

	return nil
}

func verifySub(engine *xorm.Engine, m map[string][]string, wx string) error {
	count := 0
	for k, v := range m {
		for _, vv := range v {
			count++
			character, e := fate.GetCharacter(engine, func(eng *xorm.Engine) *xorm.Session {
				return eng.Where("ch = ?", vv)
			})
			if e != nil {
				log.Errorw("get character error", "character", vv)
				continue
			}
			if character.WuXing != wx {
				if character.WuXing == "" {
					//fix wuxing
					character.WuXing = wx
				} else {
					log.Warnw("wrong wuxing", "character", vv, "charwuxing", character.WuXing, "dictwuxing", wx)
				}
			}
			i, _ := strconv.Atoi(k)
			if character.ScienceStroke != i {
				log.Warnw("check warning", "character", vv, "db", character.ScienceStroke, "need", k)
				//fix stroke
				character.ScienceStroke = i
			}
			update, e := engine.Where("hash = ?", character.Hash).Update(character)
			if e != nil {
				return e
			}
			if update != 1 {
				//log.Errorw("not updated", "update", update)
			}
		}
	}
	log.Infow("total", "count", count)
	return nil
}
