package scripts

import (
	"encoding/json"
	"os"
)

type Characters []Character

type Character struct {
	Char           string          `json:"char"`
	Pronunciations []Pronunciation `json:"pronunciations"`
	Index          int64           `json:"index"`
}

type Pronunciation struct {
	Pinyin       string        `json:"pinyin"`
	Explanations []Explanation `json:"explanations"`
}

type Explanation struct {
	Content    *string  `json:"content,omitempty"`
	Detail     []Detail `json:"detail"`
	Words      any      `json:"words"`
	Same       *string  `json:"same,omitempty"`
	Example    *string  `json:"example,omitempty"`
	Refer      *string  `json:"refer,omitempty"`
	Modern     *string  `json:"modern,omitempty"`
	Speech     *Speech  `json:"speech,omitempty"`
	Pinyin     *string  `json:"pinyin,omitempty"`
	Simplified *string  `json:"simplified,omitempty"`
	Variant    *string  `json:"variant,omitempty"`
	Unknown    *bool    `json:"unknown,omitempty"`
	Typo       *string  `json:"typo,omitempty"`
}

type Detail struct {
	Text  string       `json:"text"`
	Book  *string      `json:"book,omitempty"`
	Words []DetailWord `json:"words"`
}

type DetailWord struct {
	Word string  `json:"word"`
	Text *string `json:"text,omitempty"`
}

type PurpleWord struct {
	Word    *string `json:"word,omitempty"`
	Text    *string `json:"text,omitempty"`
	Example *string `json:"example,omitempty"`
	Refer   *string `json:"refer,omitempty"`
	Pinyin  *string `json:"pinyin,omitempty"`
}

type Speech string

const (
// 代  Speech = "代"
// 前缀 Speech = "前缀"
// 副  Speech = "副"
// 动  Speech = "动"
// 助  Speech = "助"
// 叹  Speech = "叹"
// 名  Speech = "名"
// 名动 Speech = "名,动"
// 形  Speech = "形"
// 数量 Speech = "数量"
// 语气 Speech = "语气"
// 象  Speech = "象"
// 量  Speech = "量"
)

type WordUnion struct {
	DetailWordArray []DetailWord
	PurpleWord      *PurpleWord
}

func LoadCharDetailJSON(path string, hook func(ch Character) bool) error {
	//load polyphone from json file
	of, err := os.Open(path)
	if err != nil {
		return err
	}
	defer of.Close()
	decoder := json.NewDecoder(of)
	var chs []Character
	err = decoder.Decode(&chs)
	if err != nil {
		return err
	}
	for _, ch := range chs {
		if !hook(ch) {
			return nil
		}
	}
	return nil
}
