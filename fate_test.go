package fate_test

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/godcong/fate"
)

func TestFate_EightCharacter(t *testing.T) {
	log.SetFlags(log.Llongfile)
	fate := fate.NewFate("毛")
	fate.SetLunarData(time.Now())
	log.Println(fate.EightCharacter())
}

func TestTime(t *testing.T) {
	tt, err := time.Parse("2006-01-02T15:04", "2017-11-14T08:17")
	if err != nil {
		panic(err)
	}
	sour := rand.NewSource(tt.UnixNano())
	fmt.Println(rand.New(sour).Uint64())
	fmt.Println(rand.New(sour).Uint64())
	//10547234835636046708
}

func TestFate_BestStrokes(t *testing.T) {
	f := fate.NewFate("毛")
	//fate.SetMartial()
	m := &fate.Martial{
		BiHua:     true,
		SanCai:    true,
		BaZi:      false,
		GuaXiang:  false,
		TianYun:   false,
		ShengXiao: false,
	}
	f.SetMartial(m)
	f.SetLunarData(time.Now())
	strokes := f.BestCharacters()
	log.Printf("%+v", strokes)
}
