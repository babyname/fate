package mongo

import "gopkg.in/mgo.v2/bson"

type DaYan struct {
	ID      bson.ObjectId `bson:"_id,omitempty"` //id
	Index   int           `bson:"index"`         //use array index
	Fortune string        `bson:"fortune"`       //吉凶
	TianJiu string        `bson:"tian_jiu"`      //天九(天九地十取天九)
	Comment string        `bson:"comment"`
}
