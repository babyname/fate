package mongo

import (
	"log"
	"strings"

	"github.com/globalsign/mgo/bson"
)

//WuXing 五行八字
type WuXing struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	WuXing  []string      `bson:"wu_xing"`
	Fortune string        `bson:"fortune"`
	Comment string        `bson:"comment"`
}

var wuxing map[string]*WuXing

func GetWuXing() map[string]*WuXing {
	return getWuXing()
}

func getWuXing() map[string]*WuXing {
	var wx WuXing
	wuxing = make(map[string]*WuXing)
	iter := C("wuxing").Find(nil).Iter()
	for iter.Next(&wx) {
		mwx := WuXing{
			ID:      wx.ID,
			WuXing:  wx.WuXing,
			Fortune: wx.Fortune,
			Comment: wx.Comment,
		}

		key := strings.Join(mwx.WuXing, "")
		log.Println(mwx)
		wuxing[key] = &mwx
	}
	return wuxing
}
