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
	eng := fate.InitDatabaseFromConfig(config.Config{
		HardFilter:   false,
		StrokeMax:    0,
		StrokeMin:    0,
		FixBazi:      false,
		SupplyFilter: false,
		ZodiacFilter: false,
		BaguaFilter:  false,
		Database: config.Database{
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
		},
	})
	c := chronos.New("2020/01/14 02:45")
	//t.Log(c.Solar().Time())
	cfg := config.DefaultConfig()
	cfg.BaguaFilter = true
	cfg.ZodiacFilter = true
	cfg.SupplyFilter = true
	cfg.HardFilter = true
	cfg.StrokeMin = 3
	cfg.StrokeMax = 24
	//cfg.FileOutput = "output.csv"
	f := fate.NewFate("刘", c.Solar().Time(), fate.DBOption(eng), fate.ConfigOption(*cfg))

	//f.SetDB(eng)
	e := f.MakeName(context.Background())
	if e != nil {
		t.Fatal(e)
	}
}
