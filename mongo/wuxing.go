package mongo

import "gopkg.in/mgo.v2/bson"

type WuXing struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	WuXing  []string      `bson:"wu_xing"`
	Fortune string        `bson:"fortune"`
	Comment string        `bson:"comment"`
}
