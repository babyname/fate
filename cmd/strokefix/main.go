package main

import (
	"github.com/godcong/fate"
	"github.com/godcong/fate/config"
	"github.com/xormsharp/xorm"
)

func main() {
	var e error

	cfg := config.DefaultConfig()
	db := fate.InitDatabaseWithConfig(*cfg)

	e = db.Sync(fate.Character{})
	if e != nil {
		return
	}

	e = UpdateFix(db.Database().(*xorm.Engine))
	if e != nil {
		panic(e)
	}
	e = CheckLoader(`E:\project\fate\cmd\strokefix\dict.json`)
	if e != nil {
		panic(e)
	}
	e = CheckVerify(db)
	if e != nil {
		panic(e)
	}
}
