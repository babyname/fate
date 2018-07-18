package fate

import (
	"strings"

	"github.com/godcong/fate/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Name struct {
	FirstName []string
	firstChar []*mongo.Character
	LastName  []string
	lastChar  []*mongo.Character
}

func newName(last string) *Name {
	name := Name{}
	if len(last) == 2 {
		name.LastName = strings.Split(last, "")
		for _, v := range name.LastName {
			name.lastChar = append(name.lastChar, nameCharacter(v))
		}
	} else {
		name.LastName[0] = last
		for i, v := range name.LastName {
			name.lastChar[i] = nameCharacter(v)
		}
	}

	return &name
}

func nameCharacter(s string) *mongo.Character {
	c := mongo.Character{}
	err := mongo.C("character").Find(bson.M{
		"character": s,
	}).One(&c)
	if err != nil {
		return nil
	}
	return &c
}
