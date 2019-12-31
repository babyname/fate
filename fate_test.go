package fate_test

import (
	"context"
	"github.com/godcong/chronos"
	"github.com/godcong/fate"
	"testing"
)

func init() {
	//trait.NewZapFileSugar("fate.log")
}
func TestFate_FirstRunInit(t *testing.T) {
	eng := fate.InitMysql("127.0.0.1", "root", "111111")

	c := chronos.New("2020/01/23 11:31")
	//t.Log(c.Solar().Time())
	fate.DefaultStrokeMin = 3
	fate.DefaultStrokeMax = 18
	fate.HardMode = true
	f := fate.NewFate("王", c.Solar().Time(), fate.Database(eng), fate.BaGuaFilter(), fate.ZodiacFilter(), fate.SupplyFilter())

	//f.SetDB(eng)
	e := f.MakeName(context.Background())
	if e != nil {
		t.Fatal(e)
	}
}
