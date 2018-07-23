//天运纳音
package mongo

import "github.com/globalsign/mgo/bson"

type TianYunNaYin struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	GanZhi         string        `bson:"gan_zhi"`
	ShengXiao      string        `bson:"sheng_xiao"`
	TianYun        string        `bson:"tian_yun"`
	TianYunYinYang string        `bson:"tian_yun_yin_yang"`
}
