package source

import (
	"encoding/json"
	"os"
)

type Word struct {
	Word        string `json:"word"`
	Oldword     string `json:"oldword"`
	Strokes     string `json:"strokes"`
	Pinyin      string `json:"pinyin"`
	Radicals    string `json:"radicals"`
	Explanation string `json:"explanation"`
	More        string `json:"more"`
}

func LoadWord(path string, hook func(w Word) bool) error {
	of, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer of.Close()

	var words []Word
	dec := json.NewDecoder(of)
	err = dec.Decode(&words)
	if err != nil {
		return err
	}
	for _, word := range words {
		if !hook(word) {
			return nil
		}
	}
	return nil
}
