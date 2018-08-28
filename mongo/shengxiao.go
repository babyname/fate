package mongo

import (
	"github.com/globalsign/mgo/bson"
)

//ShengXiao 生肖喜忌
type ShengXiao struct {
	ID        bson.ObjectId `bson:"_id,omitempty"` //id
	Character string        `bson:"character"`
	XiShu     bool          `bson:"xi_shu"`
	XiNiu     bool          `bson:"xi_niu"`
	XiHu      bool          `bson:"xi_hu"`
	XiTu      bool          `bson:"xi_tu"`
	XiLong    bool          `bson:"xi_long"`
	XiShe     bool          `bson:"xi_she"`
	XiMa      bool          `bson:"xi_ma"`
	XiYang    bool          `bson:"xi_yang"`
	XiHou     bool          `bson:"xi_hou"`
	XiJi      bool          `bson:"xi_ji"`
	XiGou     bool          `bson:"xi_gou"`
	XiZhu     bool          `bson:"xi_zhu"`
	JiShu     bool          `bson:"ji_shu"`
	JiNiu     bool          `bson:"ji_niu"`
	JiHu      bool          `bson:"ji_hu"`
	JiTu      bool          `bson:"ji_tu"`
	JiLong    bool          `bson:"ji_long"`
	JiShe     bool          `bson:"ji_she"`
	JiMa      bool          `bson:"ji_ma"`
	JiYang    bool          `bson:"ji_yang"`
	JiHou     bool          `bson:"ji_hou"`
	JiJi      bool          `bson:"ji_ji"`
	JiGou     bool          `bson:"ji_gou"`
	JiZhu     bool          `bson:"ji_zhu"`
}
