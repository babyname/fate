package mongo

import (
	"strings"

	"github.com/globalsign/mgo/bson"
)

type GuaXiang struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	ShangGua string        `bson:"shang_gua"`
	ShangShu int           `bson:"shang_shu"`
	XiaGua   string        `bson:"xia_gua"`
	XiaShu   int           `bson:"xia_shu"`
	GuaXiang string        `bson:"gua_xiang"`
	GuaMing  string        `bson:"gua_ming"`
	GuaYi    string        `bson:"gua_yi"`
	GuaYun   string        `bson:"gua_yun"`
	XiangYue string        `bson:"xiang_yue"`
	FuHao    string        `bson:"fu_hao"`
}

var gua map[string]*GuaXiang

func GetGuaXiang() map[string]*GuaXiang {
	return getGuaXiang()
}
func getGuaXiang() map[string]*GuaXiang {
	if gua == nil {
		gua = make(map[string]*GuaXiang)
		var zy GuaXiang
		iter := C("zhouyi").Find(nil).Iter()
		for iter.Next(&zy) {
			mzy := GuaXiang{
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
			gua[key] = &mzy
		}
	}
	return gua
}
