package fate

//三才五格吉凶

import (
	"strings"

	"github.com/godcong/yi"
)

type WuGeDaYan struct {
	TianDaYan string
	RenDaYan  string
	DiDaYan   string
	WaiDaYan  string
	ZongDaYan string
}

func (gd *WuGeDaYan) IsLucky(filter_hard bool) bool {
	var islucky bool = false

	// ignore gd.TianDaYan

	if strings.Contains(gd.RenDaYan, "吉") && !strings.Contains(gd.RenDaYan, "凶") &&
		strings.Contains(gd.DiDaYan, "吉") && !strings.Contains(gd.DiDaYan, "凶") &&
		strings.Contains(gd.ZongDaYan, "吉") && !strings.Contains(gd.ZongDaYan, "凶") {

		if filter_hard {
			if strings.Contains(gd.RenDaYan, "吉") && !strings.Contains(gd.RenDaYan, "凶") &&
				strings.Contains(gd.DiDaYan, "吉") && !strings.Contains(gd.DiDaYan, "凶") &&
				strings.Contains(gd.ZongDaYan, "吉") && !strings.Contains(gd.ZongDaYan, "凶") &&
				!strings.Contains(gd.ZongDaYan, "半") && !strings.Contains(gd.ZongDaYan, "平") &&
				strings.Contains(gd.WaiDaYan, "吉") && !strings.Contains(gd.WaiDaYan, "凶") {
				islucky = true
			}
		} else {
			islucky = true
		}

	}

	return islucky
}

var wuGeDayans map[int][]*WuGeDaYan = map[int][]*WuGeDaYan{}

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

func GetWuGeDaYan(nk NameStroke, sex yi.Sex, filter_hard bool) []*WuGeDaYan {
	hash := GetWuGeDaYanHash(nk, sex)

	wuGeDaYan := wuGeDayans[hash]

	if wuGeDaYan != nil {
		return wuGeDaYan
	}

	var wuGeDaYans_new []*WuGeDaYan = []*WuGeDaYan{}

	for _, wuGe := range GetWuGe(nk, filter_hard) {
		wuGeDaYans_new = append(wuGeDaYans_new, wuGe.getWuGeDaYan(sex))
	}

	wuGeDayans[hash] = wuGeDaYans_new

	return wuGeDaYans_new
}

func (ge *WuGe) getWuGeDaYan(sex yi.Sex) *WuGeDaYan {
	wuGeDaYan := WuGeDaYan{}

	if daYan := yi.GetDaYan(ge.tianGe); sex == yi.SexGirl && strings.Contains(daYan.NvMing, "凶") {
		wuGeDaYan.TianDaYan = daYan.NvMing
	} else {
		wuGeDaYan.TianDaYan = daYan.Lucky
	}

	if daYan := yi.GetDaYan(ge.renGe); sex == yi.SexGirl && strings.Contains(daYan.NvMing, "凶") {
		wuGeDaYan.RenDaYan = daYan.NvMing
	} else {
		wuGeDaYan.RenDaYan = daYan.Lucky
	}

	if daYan := yi.GetDaYan(ge.diGe); sex == yi.SexGirl && strings.Contains(daYan.NvMing, "凶") {
		wuGeDaYan.DiDaYan = daYan.NvMing
	} else {
		wuGeDaYan.DiDaYan = daYan.Lucky
	}

	if daYan := yi.GetDaYan(ge.waiGe); sex == yi.SexGirl && strings.Contains(daYan.NvMing, "凶") {
		wuGeDaYan.WaiDaYan = daYan.NvMing
	} else {
		wuGeDaYan.WaiDaYan = daYan.Lucky
	}

	if daYan := yi.GetDaYan(ge.zongGe); sex == yi.SexGirl && strings.Contains(daYan.NvMing, "凶") {
		wuGeDaYan.ZongDaYan = daYan.NvMing
	} else {
		wuGeDaYan.ZongDaYan = daYan.Lucky
	}

	return &wuGeDaYan
}
