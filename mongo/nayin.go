package mongo

import "github.com/globalsign/mgo/bson"

//NaYin
type NaYin struct {
	ID       bson.ObjectId `bson:"_id,omitempty"` //id
	GanZhi   []string                             //干支
	WuXing   string                               //五行
	ZhiLiang string                               //质量
	Comment  string                               //说明
}
