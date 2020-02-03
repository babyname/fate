package fate_test

import (
	"context"
	"github.com/godcong/chronos"
	"github.com/godcong/fate"
	"github.com/godcong/fate/config"
	"testing"
)

func init() {
	//trait.NewZapFileSugar("fate.log")
}

func TestFate_RunMakeName(t *testing.T) {
	born := chronos.New("2020/01/14 02:45").Solar().Time()
	last := "张"
	cfg := config.DefaultConfig()
	cfg.BaguaFilter = true
	cfg.ZodiacFilter = true
	cfg.SupplyFilter = true
	cfg.HardFilter = true
	cfg.StrokeMin = 3
	cfg.StrokeMax = 24
	cfg.RunInit = false
	cfg.FileOutput = config.FileOutput{
		OutputMode: config.OutputModeCSV,
		Path:       "name2.csv",
	}
	cfg.Database = config.Database{
		Host:         "localhost",
		Port:         "3306",
		User:         "root",
		Pwd:          "111111",
		Name:         "fate",
		MaxIdleCon:   0,
		MaxOpenCon:   0,
		Driver:       "mysql",
		File:         "",
		Dsn:          "",
		ShowSQL:      false,
		ShowExecTime: false,
	}
	f := fate.NewFate(last, born, fate.ConfigOption(cfg), fate.SexOption(fate.SexGirl))

	//f.SetDB(eng)
	e := f.MakeName(context.Background())
	if e != nil {
		t.Fatal(e)
	}
}
