package source

import (
	"encoding/json"
	"os"
)

type PolyPhone struct {
	Index     int      `json:"index"`
	Char      string   `json:"char"`
	Pinyin    []string `json:"pinyin"`
	Frequency int      `json:"frequency"`
	Strokes   int      `json:"strokes"`
}

func Load(path string) []PolyPhone {
	//load polyphone from json file
	of, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer of.Close()
	decoder := json.NewDecoder(of)
	var polyphone []PolyPhone
	err = decoder.Decode(&polyphone)
	if err != nil {
		return nil
	}
	return polyphone
}
