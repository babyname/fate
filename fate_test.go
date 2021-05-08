package fate_test

import (
	"context"
	"testing"

	"github.com/godcong/chronos"
	"github.com/godcong/fate"
	"github.com/godcong/fate/config"
	"github.com/godcong/yi"
)

func init() {
	//trait.NewZapFileSugar("fate.log")
}

func TestFate_RunMakeName(t *testing.T) {

	//出生日期
	born := chronos.New("2020/02/06 15:04").Solar().Time()
	//姓氏
	last := "张"
	xiyong := "火火火"
	cfg := config.LoadConfig()
	cfg.Database.Driver = "sqlite3"
	cfg.Database.File = "fate.db"
	cfg.BaguaFilter = true
	cfg.ZodiacFilter = true
	cfg.SupplyFilter = true
	cfg.HardFilter = true
	cfg.StrokeMin = 3
	cfg.StrokeMax = 24
	cfg.Regular = true
	cfg.RunInit = false
	cfg.FileOutput = config.FileOutput{
		OutputMode: config.OutputModeLog,
		Path:       "name.log",
	}
	f := fate.NewFate(last, born, fate.ConfigOption(cfg), fate.SexOption(yi.SexBoy), fate.XiYongOption(xiyong))

	//f.SetDB(eng)
	e := f.MakeName(context.Background())
	if e != nil {
		t.Fatal(e)
	}
}
