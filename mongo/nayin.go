package mongo

import "github.com/globalsign/mgo/bson"

//NaYin
type NaYin struct {
	ID       bson.ObjectId `bson:"_id,omitempty"` //id
	GanZhi   []string      `json:"gan_zhi"`       //干支
	WuXing   string        `json:"wu_xing"`       //五行
	ZhiLiang string        `json:"zhi_liang"`     //质量
	Comment  string        `json:"comment"`       //说明
}
