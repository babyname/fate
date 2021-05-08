package fate_test

import (
	"log"
	"testing"

	"github.com/godcong/chronos"
	"github.com/godcong/fate"
)

func TestPoint(t *testing.T) {
	t1 := chronos.New("1986/08/20 11:30")
	log.Println(t1.Lunar().EightCharacter())

	bz := fate.NewBazi(t1)
	t.Log(bz.XiYong())
	t.Log(bz.XiYongShen())
}
