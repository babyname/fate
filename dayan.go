package fate

//三才五格吉凶

import (
	nw "github.com/godcong/name_wuge"
	"github.com/godcong/yi"
)

var wuGeDayans map[int][]*nw.WuGeDaYan = map[int][]*nw.WuGeDaYan{}

func ClearWuGeDaYan(nk NameStroke, sex yi.Sex) {
	hash := GetWuGeDaYanHash(nk, sex)
	if wuGeDayans[hash] != nil {
		delete(wuGeDayans, hash)
	}
}

// 1bit~0bit 男为0b01，女为0b10，详情见yi中定义
func GetWuGeDaYanHash(nk NameStroke, sex yi.Sex) int {
	hash := nk.hash()
	hash = int(sex) + hash<<2
	return hash
}

func GetWuGeDaYan(nk NameStroke, sex yi.Sex, filter_hard bool) []*nw.WuGeDaYan {
	hash := GetWuGeDaYanHash(nk, sex)

	wuGeDaYan := wuGeDayans[hash]

	if wuGeDaYan != nil {
		return wuGeDaYan
	}

	var wuGeDaYans_new []*nw.WuGeDaYan = []*nw.WuGeDaYan{}

	for _, wuGe := range GetWuGe(nk, filter_hard) {
		wuGeDaYans_new = append(wuGeDaYans_new, wuGe.GetWuGeDaYan(sex))
	}

	wuGeDayans[hash] = wuGeDaYans_new

	return wuGeDaYans_new
}
