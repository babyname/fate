package fate

//三才五格吉凶

import (
	nw "github.com/godcong/name_wuge"
	"github.com/godcong/yi"
)

func (nk *NameStroke) getWuGeDaYans(sex yi.Sex, filter_hard bool) []*nw.WuGeDaYan {
	wuGeDaYans_new := []*nw.WuGeDaYan{}

	for _, wuGe := range nk.GetWuGes(filter_hard) {
		wuGeDaYans_new = append(wuGeDaYans_new, wuGe.GetWuGeDaYan(sex))
	}

	return wuGeDaYans_new
}
