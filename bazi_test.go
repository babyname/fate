package fate_test

import (
	"github.com/godcong/chronos"
	"github.com/babyname/fate"
	"log"
	"testing"
)

func TestPoint(t *testing.T) {
	t1 := chronos.New("2020/01/24 15:30")
	log.Println(t1.Lunar().EightCharacter())

	bz := fate.NewBazi(t1)
	t.Log(bz.XiYong())
	t.Log(bz.XiYongShen())
}
