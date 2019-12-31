package main

import (
	"encoding/json"
	"errors"
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

func CheckLoader(s string) error {
	bytes, e := ioutil.ReadFile(s)
	if e != nil {
		return e
	}

	e = json.Unmarshal(bytes, &checks)
	if e != nil {
		return e
	}

	if len(checks) == 0 {
		return errors.New("no data")
	}

	return nil
}

func CheckVerify(db *xorm.Engine) error {
	for _, s := range checks {
		for five, val := range s {
			for idx, vv := range val {
				log.Infof("%+v,val:%+v", idx, val)
				continue
				log.Infow("check", "character", s)
				character, e := fate.GetCharacter(db, func(eng *xorm.Engine) *xorm.Session {
					return eng.Where("ch = ?", s)
				})
				if e != nil {
					log.Errorw("check error", "character", s)
					return e
				}
				i, _ := strconv.Atoi(idx)
				if character.ScienceStroke != i {
					log.Warnw("check warning", "character", s, "db", character.ScienceStroke, "need", idx)
				}
			}

		}
	}
	return nil
}
