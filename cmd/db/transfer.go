package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/xormsharp/xorm"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/model"
	"github.com/babyname/fate/olddata"
)

func runTransfer(p string) error {
	db, err := readConfig(p)
	if err != nil {
		return err
	}

	olddb, err := openOldDatabase(dsn, db.From)
	if err != nil {
		return err
	}
	defer olddb.Close()
	newdb, err := openNewDatabase(db.To.DSN, db.To)
	if err != nil {
		return err
	}
	defer newdb.Close()

	var total int64
	if total, err = transferCharacter(olddb, newdb); err != nil {
		return err
	}

	fmt.Println("Finished reading total character:", total)
	return nil
}

func transferCharacter(olddb *xorm.Engine, newdb *model.Model) (total int64, err error) {
	total, err = olddata.RangeCharacters(olddb, func(c *olddata.Character) bool {
		ch := &ent.Character{
			ID:                       model.ID(c.Ch),
			PinYin:                   c.PinYin,
			Ch:                       c.Ch,
			ScienceStroke:            c.ScienceStroke,
			Radical:                  c.Radical,
			RadicalStroke:            c.RadicalStroke,
			Stroke:                   c.Stroke,
			IsKangxi:                 c.IsKangXi,
			Kangxi:                   c.KangXi,
			KangxiStroke:             c.KangXiStroke,
			SimpleRadical:            c.SimpleRadical,
			SimpleRadicalStroke:      c.SimpleRadicalStroke,
			SimpleTotalStroke:        c.SimpleTotalStroke,
			TraditionalRadical:       c.TraditionalRadical,
			TraditionalRadicalStroke: c.TraditionalRadicalStroke,
			TraditionalTotalStroke:   c.TraditionalTotalStroke,
			IsNameScience:            c.NameScience,
			WuXing:                   c.WuXing,
			Lucky:                    c.Lucky,
			IsRegular:                c.Regular,
			TraditionalCharacter:     c.TraditionalCharacter,
			VariantCharacter:         c.VariantCharacter,
			//Comment:                  c.Comment,
		}
		character, e := newdb.InsertOrUpdateCharacter(context.TODO(), ch)
		if e != nil {
			fmt.Println("ERROR:", e)
			return false
		}
		fmt.Println("character:", character.Ch, " id:", character.ID)
		return e == nil
	})
	return
}

func openOldDatabase(dsn string, f From) (*xorm.Engine, error) {
	dbURL := fmt.Sprintf(dsn, f.User, f.Pwd, f.Host+":"+f.Port, f.DBName, url.QueryEscape("Asia/Shanghai"))
	fmt.Println("old dsn:", dbURL)
	dbEngine, err := xorm.NewEngine(f.Driver, dbURL)
	if err != nil {
		return nil, err
	}
	return dbEngine, nil
}

func openNewDatabase(dsn string, t To) (*model.Model, error) {
	open, err := model.Open(config.DBConfig{
		Driver: "sqlite3",
		DSN:    dsn,
	}, false)
	if err != nil {
		return nil, err
	}
	err = open.Schema.Create(context.TODO())
	if err != nil {
		return nil, err
	}
	return open, nil
}
