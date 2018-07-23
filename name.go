package fate

import (
	"log"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/godcong/fate/mongo"
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
		name.LastName = []string{last}
		for _, v := range name.LastName {
			name.lastChar = []*mongo.Character{nameCharacter(v)}
		}
	}

	return &name
}

func nameCharacter(s string) *mongo.Character {
	c := mongo.Character{}
	err := mongo.C("character").Find(bson.M{
		"character": s,
	}).One(&c)
	log.Printf("%+v", c)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &c
}
