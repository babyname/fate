package fate

import (
	"strings"

	"github.com/godcong/fate/model"
)

type Name struct {
	FirstName []string
	cFirst    []*model.Character
	LastName  []string
	cLast     []*model.Character
}

func NewName(last string) Name {
	name := Name{}
	if len(last) > 1 {
		name.LastName = strings.Split(last, "")
		for i, v := range name.LastName {
			name.cFirst[i] = CharacterFromName(v)
		}
		return name
	}
	name.LastName[0] = last
	for i, v := range name.LastName {
		name.cFirst[i] = CharacterFromName(v)
	}
	return name
}

func CharacterFromName(s string) *model.Character {
	c := &model.Character{
		IsSur:          false,
		SimpleChar:     s,
		SimpleStrokes:  0,
		TradChar:       "",
		TradStrokes:    0,
		NameType:       "",
		NameRoot:       "",
		Radical:        "",
		ScienceStrokes: 0,
		Pinyin:         "",
		Comment:        "",
	}
	return c.Get()
}
