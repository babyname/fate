package main

import (
	"github.com/godcong/fate"
	"github.com/goextension/log"
	"github.com/xormsharp/xorm"
)

func UpdateFix(engine *xorm.Engine) error {
	rows, e := engine.Rows(&fate.Character{})
	if e != nil {
		return e
	}
	var ch *fate.Character
	for rows.Next() {
		ch = &fate.Character{}
		e := rows.Scan(ch)
		if e != nil {
			log.Errorw("fix", "charater", ch.Ch, "error", e)
			continue
		}
		if fixChar(ch) {
			_, e := engine.Where("hash = ?", ch.Hash).Update(&ch)
			if e != nil {
				log.Errorw("update", "charater", ch.Ch, "error", e)
				continue
			}
			log.Infow("updated", "charater", ch.Ch, "stroke", ch.ScienceStroke)
		}
	}
	return nil
}
