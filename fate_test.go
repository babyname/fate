package fate_test

import (
	"github.com/godcong/fate"
	"testing"
	"time"
)

func TestFate_FirstRunInit(t *testing.T) {
	f := fate.NewFate("李", time.Now())
	eng, e := fate.NewSQLite3("./data/character.db")
	if e != nil {
		t.Fatal(e)
	}
	f.SetCharDB(eng)
	e = f.GetLastCharacter()
	if e != nil {
		t.Fatal(e)
	}
}
