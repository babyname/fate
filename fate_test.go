package fate_test

import (
	"github.com/godcong/fate"
	"testing"
	"time"
)

func TestFate_FirstRunInit(t *testing.T) {
	eng := fate.InitMysql("localhost:3306", "root", "111111")
	f := fate.NewFate("张周李", time.Now(), fate.Database(eng))

	//f.SetDB(eng)
	e := f.MakeName()
	if e != nil {
		t.Fatal(e)
	}
}
