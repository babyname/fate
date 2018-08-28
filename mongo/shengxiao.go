package mongo

import (
	"github.com/globalsign/mgo/bson"
)

//ShengXiao 生肖喜忌
type ShengXiao struct {
	ID        bson.ObjectId `bson:"_id,omitempty"` //id
	Character string
	XiShu     bool
	XiNiu     bool
	XiHu      bool
	XiTu      bool
	XiLong    bool
	XiShe     bool
	XiMa      bool
	XiYang    bool
	XiHou     bool
	XiJi      bool
	XiGou     bool
	XiZhu     bool
	JiShu     bool
	JiNiu     bool
	JiHu      bool
	JiTu      bool
	JiLong    bool
	JiShe     bool
	JiMa      bool
	JiYang    bool
	JiHou     bool
	JiJi      bool
	JiGou     bool
	JiZhu     bool
}
