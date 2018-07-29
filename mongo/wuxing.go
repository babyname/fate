package mongo

import "github.com/globalsign/mgo/bson"

//WuXing 五行八字
type WuXing struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	WuXing  []string      `bson:"wu_xing"`
	Fortune string        `bson:"fortune"`
	Comment string        `bson:"comment"`
}
