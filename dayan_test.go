package fate

import (
	"encoding/json"
	"github.com/godcong/go-trait"
	"testing"
)

func TestCompress(t *testing.T) {
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

func TestGetDayan(t *testing.T) {
	for i := 1; i < 100; i++ {
		dy := GetDaYan(i)
		//output i == number
		t.Logf("i:%d,%+v", i, dy)
	}
}
