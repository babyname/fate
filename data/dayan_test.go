package data

import (
	"encoding/json"
	"github.com/godcong/go-trait"
	"testing"
)

func TestDayan(t *testing.T) {
	bytes, e := json.Marshal(DaYanList)
	if e != nil {
		panic(e)
	}
	e = trait.CompressToFile("dayan.data", bytes)
	if e != nil {
		panic(e)
	}
}

func TestDecompress(t *testing.T) {
	for _, value := range DaYanList {
		if value.Number == 0 {
			t.Error(value)
		}
	}
}
