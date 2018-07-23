package mongo

import "github.com/globalsign/mgo/bson"

type WuXing struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	WuXing  []string      `bson:"wu_xing"`
	Fortune string        `bson:"fortune"`
	Comment string        `bson:"comment"`
}
