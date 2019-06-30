package fate

import (
	"log"
	"strconv"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/godcong/fate/mongo"
	"github.com/godcong/yi"
)

//Name 姓名
type Name struct {
	FirstName []string //名姓
	LastName  []string
	baGua     *yi.Yi //八卦
}

func MakeName(last string) *Name {
	return &Name{
		FirstName: nil,
		LastName:  strings.Split(last, ""),
		baGua:     nil,
	}
}

//MakeName input the lastname to make a name
func FilterName(names []*Name) *Name {
	//todo:make name generate list

	return &Name{}
}

//BaGua
func (n *Name) BaGua() *yi.Yi {
	return n.baGua
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
