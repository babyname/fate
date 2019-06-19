package fate

import (
	"log"
	"strconv"

	"github.com/globalsign/mgo/bson"
	"github.com/godcong/fate/mongo"
)

//Name 姓名
type Name struct {
	FirstName string //名
	LastName  string //姓
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

//CountStroke 统计笔画
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
