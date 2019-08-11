package data

import (
	"encoding/json"
	_ "github.com/godcong/fate/statik"
	"github.com/godcong/go-trait"
)

var DaYanList [81]DaYan

func init() {
	bytes, e := trait.DecompressFromStatik("dayan.data")
	if e != nil {
		panic(e)
	}
	e = json.Unmarshal(bytes, &DaYanList)
	if e != nil {
		panic(e)
	}
}

type DaYan struct {
	Number  int
	Lucky   string
	Max     bool
	Sex     bool //male(false),female(true)
	SkyNine string
	Comment string
}
