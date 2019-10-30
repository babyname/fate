package fate_test

import (
	"github.com/godcong/fate"
	"testing"
	"time"
)

func TestFate_FirstRunInit(t *testing.T) {
	eng := fate.InitMysql("192.168.1.161:3306", "root", "111111")
	f := fate.NewFate("蒋", time.Now(), fate.Database(eng))

	//f.SetDB(eng)
	e := f.MakeName()
	if e != nil {
		t.Fatal(e)
	}
}
