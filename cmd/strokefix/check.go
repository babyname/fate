package main

import (
	"encoding/json"
	"github.com/godcong/fate"
	"github.com/goextension/log"
	"github.com/xormsharp/xorm"
	"golang.org/x/xerrors"
	"io/ioutil"
)

var checks map[int]string

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
		return xerrors.New("no data")
	}

	return nil
}

func CheckVerify(db *xorm.Engine) error {
	for idx, s := range checks {
		log.Infow("check", "character", s)
		character, e := fate.GetCharacter(db, func(eng *xorm.Engine) *xorm.Session {
			return eng.Where("ch = ?", s)
		})
		if e != nil {
			log.Errorw("check error", "character", s)
			continue
		}
		if character.ScienceStroke != idx {
			log.Warnw("check warning", "character", s, "db", character.ScienceStroke, "need", idx)
		}
	}

}
