package bazi_test

import (
	"log"
	"testing"

	"github.com/babyname/fate/internal/bazi"
	"github.com/babyname/chronos/v2"
)

func TestPoint(t *testing.T) {
	t1 := chronos.ParseSolarString("2020/01/24 15:30")
	log.Println(t1.Lunar().GetEightChar())

	bz := bazi.NewBazi(t1)
	t.Log(bz.XiYong())
	t.Log(bz.XiYongShen())
}
