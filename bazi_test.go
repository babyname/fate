package fate_test

import (
	"log"
	"testing"

	"github.com/babyname/fate"
	"github.com/godcong/chronos/v2"
)

func TestPoint(t *testing.T) {
	t1 := chronos.ParseSolarString("2020/01/24 15:30")
	log.Println(t1.Lunar().GetEightChar())

	bz := fate.NewBazi(t1)
	t.Log(bz.XiYong())
	t.Log(bz.XiYongShen())
}
