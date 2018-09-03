package fate

import (
	"log"
	"strconv"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/godcong/fate/mongo"
)

//Name 姓名
type Name struct {
	FirstName []string //名
	firstChar []*mongo.Character
	LastName  []string //姓
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

func CountStroke(chars ...*mongo.Character) int {
	i := 0
	if chars == nil {
		return i
	}
	for k := range chars {
		t, _ := strconv.Atoi(chars[k].KangxiStrokes)
		i += t
	}
	return i
}
