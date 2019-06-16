package mongo

import (
	"strings"

	"github.com/globalsign/mgo/bson"
)

//WuXing 五行：five elements of metal,wood,water,fire and earth
type WuXing struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	WuXing  string        `bson:"wu_xing"`
	Fortune string        `bson:"fortune"`
	Comment string        `bson:"comment"`
}

var wuxing map[string]*WuXing

func FindWuXingBy(s ...string) *WuXing {
	key := strings.Join(s, "")
	wx := GetWuXing()
	if v, b := wx[key]; b {
		return v
	}
	return nil
}

func GetWuXing() map[string]*WuXing {
	return getWuXing()
}

func getWuXing() map[string]*WuXing {
	if wuxing == nil {
		wuxing = make(map[string]*WuXing)
		var wx WuXing
		iter := C("wuxing").Find(nil).Iter()
		for iter.Next(&wx) {
			mwx := WuXing{
				ID:      wx.ID,
				WuXing:  wx.WuXing,
				Fortune: wx.Fortune,
				Comment: wx.Comment,
			}

			key := strings.Join(mwx.WuXing, "")
			wuxing[key] = &mwx
		}
	}
	return wuxing
}
