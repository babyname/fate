package mongo_test

import (
	"github.com/godcong/fate/model"
	"github.com/godcong/fate/mongo"
	"gopkg.in/mgo.v2"
	"log"
	"testing"
)

func TestDaYan(t *testing.T) {
	mongo.Dial("localhost", &mgo.Credential{
		Username: "root",
		Password: "v2RgzSuIaBlx",
	})

	var by []*model.BigYan
	err := model.DB().Find(&by)
	if err != nil {
		panic(err)
	}
	for idx := range by {
		b := by[idx]
		dy := mongo.DaYan{
			Index:   b.Index,
			Fortune: b.Goil,
			TianJiu: b.SkyNine,
			Comment: b.Comment,
		}
		log.Println(mongo.InsertIfNotExist(mongo.C("dayan"), &dy))
	}

}
