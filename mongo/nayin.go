package mongo

import "github.com/globalsign/mgo/bson"

//NaYin
type NaYin struct {
	ID       bson.ObjectId `bson:"_id,omitempty"` //id
	GanZhi   []string
	WuXing   string
	ZhiLiang string
}
