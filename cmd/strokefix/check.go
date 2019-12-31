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
	verifySub(db, dict.Jin, func(s string) bool {
		return s == "金"
	})
	verifySub(db, dict.Mu, func(s string) bool {
		return s == "木"
	})
	verifySub(db, dict.Shui, func(s string) bool {
		return s == "水"
	})
	verifySub(db, dict.Huo, func(s string) bool {
		return s == "火"
	})
	verifySub(db, dict.Tu, func(s string) bool {
		return s == "土"
	})

	return nil
}

func verifySub(engine *xorm.Engine, m map[string][]string, fn WuXingFunc) error {
	for k, v := range m {
		for _, vv := range v {
			character, e := fate.GetCharacter(engine, func(eng *xorm.Engine) *xorm.Session {
				return eng.Where("ch = ?", vv)
			})
			if e != nil {
				log.Errorw("check error", "character", vv)
				continue
			}
			if !fn(character.WuXing) {
				log.Warnw("wrong wuxing", "character", vv, "wuxing", character.WuXing)
			}
			i, _ := strconv.Atoi(k)
			if character.ScienceStroke != i {
				log.Warnw("check warning", "character", vv, "db", character.ScienceStroke, "need", k)
			}
		}

	}
	return nil
}
