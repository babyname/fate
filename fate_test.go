package fate_test

import (
	"github.com/godcong/chronos"
	"github.com/godcong/fate"
	"testing"
)

func TestFate_FirstRunInit(t *testing.T) {
	eng := fate.InitMysql("192.168.1.161:3306", "root", "111111")
	c := chronos.New("2020/01/23 11:31")
	//t.Log(c.Solar().Time())
	f := fate.NewFate("王", c.Solar().Time(), fate.Database(eng), fate.BaGuaFilter(), fate.ZodiacFilter(), fate.SupplyFilter())

	//f.SetDB(eng)
	e := f.MakeName()
	if e != nil {
		t.Fatal(e)
	}
}
