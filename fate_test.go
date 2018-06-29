package fate_test

import (
	"log"
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
