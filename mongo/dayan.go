package mongo

import "github.com/globalsign/mgo/bson"

//DaYan 大衍之数
type DaYan struct {
	ID      bson.ObjectId `bson:"_id,omitempty"` //id
	Index   int           `bson:"index"`         //use array index
	Fortune string        `bson:"fortune"`       //吉凶
	TianJiu string        `bson:"tian_jiu"`      //天九(天九地十取天九)
	Comment string        `bson:"comment"`
}
