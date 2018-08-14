package mongo

import (
	"strings"

	"github.com/globalsign/mgo/bson"
)

type ZhouYi struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	ShangGua string        `bson:"shang_gua"`
	XiaGua   string        `bson:"xia_gua"`
	GuaXiang string        `bson:"gua_xiang"`
	GuaMing  string        `bson:"gua_ming"`
	GuaYi    string        `bson:"gua_yi"`
	GuaYun   string        `bson:"gua_yun"`
	XiangYue string        `bson:"xiang_yue"`
}

var zhouyi map[string]*ZhouYi

func GetZhouYi() map[string]*ZhouYi {
	return getZhouYi()
}
func getZhouYi() map[string]*ZhouYi {
	if zhouyi == nil {
		zhouyi = make(map[string]*ZhouYi)
		var zy ZhouYi
		iter := C("zhouyi").Find(nil).Iter()
		for iter.Next(&zy) {
			mzy := ZhouYi{
				ID:       zy.ID,
				ShangGua: zy.ShangGua,
				XiaGua:   zy.XiaGua,
				GuaXiang: zy.GuaXiang,
				GuaMing:  zy.GuaMing,
				GuaYi:    zy.GuaYi,
				GuaYun:   zy.GuaYun,
				XiangYue: zy.XiangYue,
			}

			key := strings.Join([]string{zy.ShangGua, zy.XiaGua}, "")
			zhouyi[key] = &mzy
		}
	}
	return zhouyi
}
