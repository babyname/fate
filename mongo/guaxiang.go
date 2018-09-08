package mongo

import (
	"strings"

	"github.com/globalsign/mgo/bson"
)

type GuaXiang struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	ShangGua string        `bson:"shang_gua"` //上卦
	ShangShu int           `bson:"shang_shu"` //上卦数
	XiaGua   string        `bson:"xia_gua"`   //下卦
	XiaShu   int           `bson:"xia_shu"`   //下卦数
	GuaXiang string        `bson:"gua_xiang"` //卦象
	GuaMing  string        `bson:"gua_ming"`  //卦名
	GuaYi    string        `bson:"gua_yi"`    //卦意
	GuaYun   string        `bson:"gua_yun"`   //卦云
	XiangYue string        `bson:"xiang_yue"` //象曰
	FuHao    string        `bson:"fu_hao"`    //符号
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
