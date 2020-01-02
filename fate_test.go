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
func TestFate_FirstRunInit(t *testing.T) {
	eng := fate.InitDatabaseFromConfig(config.Config{
		HardMode:     false,
		StrokeMax:    0,
		StrokeMin:    0,
		FixBazi:      false,
		SupplyFilter: false,
		ZodiacFilter: false,
		BaguaFilter:  false,
		Database: config.Database{
			Host:         "",
			Port:         "",
			User:         "",
			Pwd:          "",
			Name:         "",
			MaxIdleCon:   0,
			MaxOpenCon:   0,
			Driver:       "",
			File:         "",
			Dsn:          "",
			ShowSQL:      false,
			ShowExecTime: false,
		},
	})
	c := chronos.New("2020/01/23 11:31")
	//t.Log(c.Solar().Time())
	fate.DefaultStrokeMin = 3
	fate.DefaultStrokeMax = 18
	fate.HardMode = true
	f := fate.NewFate("王", c.Solar().Time(), fate.DBOption(eng), fate.BaGuaFilter(), fate.ZodiacFilter(), fate.SupplyFilter())

	//f.SetDB(eng)
	e := f.MakeName(context.Background())
	if e != nil {
		t.Fatal(e)
	}
}
