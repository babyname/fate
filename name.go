package fate

import (
	"strings"

	"github.com/godcong/fate/mongo"
	"gopkg.in/mgo.v2/bson"
)

const LenMax = 32

type Name struct {
	FirstName []string
	cFirst    []*mongo.Character
	LastName  []string
	cLast     []*mongo.Character
}

func newName(last string) *Name {
	name := Name{}
	if len(last) > 1 {
		name.LastName = strings.Split(last, "")
		for _, v := range name.LastName {
			name.cLast = append(name.cLast, characterFromName(v))
		}
	} else {
		name.LastName[0] = last
		for i, v := range name.LastName {
			name.cLast[i] = characterFromName(v)
		}
	}

	return &name
}

func characterFromName(s string) *mongo.Character {
	var c mongo.Character
	err := mongo.C("character").Find(bson.M{
		"character": s,
	}).One(&c)
	if err != nil {
		return nil
	}
	return &c
}
