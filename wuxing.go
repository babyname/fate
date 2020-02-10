package fate

import (
	"errors"
	"github.com/xormsharp/xorm"
)

// Luck ...
type Luck int

var luckPoint = []string{"大凶", "凶", "凶多于吉", "吉凶参半", "吉多于凶", "吉", "大吉"}

// Point ...
func (l *Luck) Point() int {
	return int(*l) + 1
}

// ToLuck ...
func ToLuck(s string) (l Luck, e error) {
	for i, luck := range luckPoint {
		if luck == s {
			return Luck(i), nil
		}
	}
	return Luck(0), errors.New("parse error")
}

//WuXing 五行：five elements of metal,wood,water,fire and earth
type WuXing struct {
	WuXing  string `json:"wu_xing"`
	Luck    Luck   `json:"luck"`
	Comment string `json:"comment"`
}

//FindWuXing find a wuxing
func FindWuXing(engine *xorm.Engine, s ...string) *WuXing {
	var wx WuXing
	_, e := engine.Where("first = ?", s[0]).And("second = ?", s[1]).And("third = ?", s[2]).Get(&wx)
	if e != nil {
		return nil
	}
	return &wx
}
