package mongo

import "github.com/globalsign/mgo/bson"

//DaYan 大衍之数
type DaYan struct {
	ID      bson.ObjectId `bson:"_id,omitempty"` //id
	Index   int           `bson:"index"`         //use array index
	Fortune string        `bson:"fortune"`       //吉凶
	TianJiu string        `bson:"tian_jiu"`      //天九(天九地十取天九)
	Comment string        `bson:"comment"`
}

var dayan []*DaYan

func init() {
	dayan = getDaYan()
}

func GetDaYan() []*DaYan {
	return getDaYan()
}

func getDaYan() []*DaYan {
	var dy []*DaYan
	if dayan == nil {
		err := C("dayan").Find(nil).Sort("index").All(&dy)
		if err != nil {
			panic(err)
		}
	}
	return dy
}
