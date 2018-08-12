package mongo_test

import (
	"github.com/globalsign/mgo"
	"github.com/godcong/fate/model"
	"github.com/godcong/fate/mongo"
	"log"
	"testing"
)

func TestWuXing(t *testing.T) {
	mongo.Dial("localhost", &mgo.Credential{
		Username: "root",
		Password: "v2RgzSuIaBlx",
	})

	var by []*model.FivePhase
	err := model.DB().Find(&by)
	if err != nil {
		panic(err)
	}
	for idx := range by {
		b := by[idx]
		dy1 := mongo.WuXing{
			WuXing:  []string{b.First, b.Second, b.Third},
			Fortune: b.Fortune,
			Comment: b.Comment,
		}
		log.Println(mongo.InsertIfNotExist(mongo.C("wuxing"), &dy1))
	}
}
