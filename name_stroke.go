package fate

import (
	nw "github.com/godcong/name_wuge"
	"github.com/godcong/yi"
)

// NameStroke ...
type NameStroke struct {
	LastStroke1  int
	LastStroke2  int
	FirstStroke1 int
	FirstStroke2 int
}

func (nk *NameStroke) IsLucky(sex yi.Sex, filter_hard bool) bool {
	if !nk.IsWuGeLucky(sex, filter_hard) {
		return false
	}

	if !nk.IsGuaLucky(sex, filter_hard) {
		return false
	}

	return true
}

//三才五格吉祥
func (nk *NameStroke) IsWuGeLucky(sex yi.Sex, filter_hard bool) bool {
	var sanCais []*nw.WuGeSanCai = []*nw.WuGeSanCai{}

	wuGeDaYans := nk.getWuGeDaYans(sex, filter_hard)

	wuGes := nk.GetWuGes(filter_hard)

	for _, wuGe := range wuGes {
		sanCais = append(sanCais, wuGe.GetSanCai())
	}

	for _, daYan := range wuGeDaYans {
		if !daYan.IsLucky(filter_hard) {
			return false
		}
	}

	for _, sanCai := range sanCais {
		if !sanCai.GetFortune().IsLucky() {
			return false
		}
	}

	return true
}

// treat int as 32 bits
// 31bit-30bit-...1bit-0bit
// 5bit~0bit nk.LastStroke2
// 11bit~6bit nk.LastStroke1
// 17bit~12it nk.FirstStroke2
// 23bit~18it nk.FirstStroke1
func GetNameStrokeHash(last1, last2, first1, first2 int) int {
	if last1 < 0 || last1 > BiHuaMax ||
		last2 < 0 || last2 > BiHuaMax ||
		first1 < 0 || first1 > BiHuaMax ||
		first2 < 0 || first2 > BiHuaMax {
		panic("bad field")
	}

	result := first1
	result += first2 + result<<6
	result += last1 + result<<6
	result += last2 + result<<6

	return result
}

func (nk *NameStroke) hash() int {
	return GetNameStrokeHash(nk.LastStroke1, nk.LastStroke2, nk.FirstStroke1, nk.FirstStroke2)
}

func addNameStroke(last1, last2, first1, first2 int, sex yi.Sex) *NameStroke {
	nameStroke := NewNameStroke(last1, last2, first1, first2)
	hash := nameStroke.hash()

	if strokeCache[hash] == nil {
		strokeCache[hash] = nameStroke
	}

	return nameStroke
}

func NewNameStroke(last1, last2, first1, first2 int) *NameStroke {
	nameStroke := NameStroke{
		LastStroke1:  last1,
		LastStroke2:  last2,
		FirstStroke1: first1,
		FirstStroke2: first2,
	}

	return &nameStroke
}

func GetNameStroke(last1, last2, first1, first2 int) *NameStroke {
	hash := GetNameStrokeHash(last1, last2, first1, first2)
	return strokeCache[hash]
}

func clearNameStroke(nk NameStroke, filter_hard bool) {
	hash := nk.hash()
	nameStroke := strokeCache[hash]

	if nameStroke != nil {
		delete(strokeCache, hash)
		ClearWuGe(*nameStroke, filter_hard)
	}
}

func CountNameStrokesLucky() int {
	return len(strokeCache)
}

// BiHuaMax ...
const BiHuaMax = 32

var strokeCache map[int]*NameStroke = map[int]*NameStroke{}

//缓存，替代之前的数据库方案
func (f *fateImpl) setStrokeCache() {
	cfg := f.config

	var lastCharLen = len(f.lastChar)

	var l1, l2 int
	if lastCharLen == 1 {
		l1 = f.lastChar[0].getStrokeScience()
		l2 = 0
	} else if lastCharLen == 2 {
		l1 = f.lastChar[0].getStrokeScience()
		l2 = f.lastChar[1].getStrokeScience()
	} else {
		panic("姓的长度不对")
	}

	for f1 := 1; f1 <= BiHuaMax; f1++ {
		for f2 := 0; f2 <= BiHuaMax; f2++ {
			nk := addNameStroke(l1, l2, f1, f2, f.sex)
			if cfg.BaguaFilter {
				if cfg.HardFilter {
					if nk.IsLucky(f.sex, true) {
						continue
					}
				} else {
					if nk.IsLucky(f.sex, false) {
						continue
					}
				}
			} else {
				if nk.IsWuGeLucky(f.sex, cfg.HardFilter) {
					continue
				}
			}

			clearNameStroke(*nk, f.config.HardFilter)
		}
	}
}

//替代之前的filterWuGe方法
func (f *fateImpl) FilterNameStrokes(wg chan<- *NameStroke) error {
	defer func() {
		close(wg)
	}()
	cfg := f.config
	var lastCharLen = len(f.lastChar)
	var l1, l2 int
	if lastCharLen == 1 {
		l1 = f.lastChar[0].getStrokeScience()
		l2 = 0
	} else if lastCharLen == 2 {
		l1 = f.lastChar[0].getStrokeScience()
		l2 = f.lastChar[1].getStrokeScience()
	} else {
		panic("姓的长度不对")
	}

	for _, v := range strokeCache {
		if v.LastStroke1 == l1 && v.LastStroke2 == l2 {
			if cfg.BaguaFilter {
				if cfg.HardFilter {
					if !v.IsLucky(f.sex, true) {
						continue
					}
				} else {
					if !v.IsLucky(f.sex, false) {
						continue
					}
				}
			} else {
				if !v.IsWuGeLucky(f.sex, cfg.HardFilter) {
					continue
				}
			}

			wg <- v
		} else {
			continue
		}
	}

	return nil
}
