package fate

import (
	"sync"

	"github.com/godcong/yi"
)

// NameStroke ...
type NameStroke struct {
	//ID     bson.ObjectId `bson:"_id,omitempty"`
	LastStroke1  int
	LastStroke2  int
	FirstStroke1 int
	FirstStroke2 int
}

//算卦象时有天人地，但是这个三才是不能用的，有很多姓氏依照五格三才是算不出吉祥的
func (nk *NameStroke) getGuaSanCai() SanCai {
	tian := nk.FirstStroke1 + nk.FirstStroke2
	di := nk.LastStroke1 + nk.LastStroke2
	ren := tian + di
	return GetSanCai(tian, ren, di)
}

func (nk *NameStroke) getWuGes(filter_hard bool) []*WuGe {
	return GetWuGe(*nk, filter_hard)
}

func (nk *NameStroke) getWuGeDaYans(sex yi.Sex, filter_hard bool) []*WuGeDaYan {
	return GetWuGeDaYan(*nk, sex, filter_hard)
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
	var sanCais []SanCai = []SanCai{}

	wuGeDaYans := nk.getWuGeDaYans(sex, filter_hard)

	wuGes := nk.getWuGes(filter_hard)

	for _, wuGe := range wuGes {
		sanCais = append(sanCais, wuGe.getSanCai())
	}

	for _, daYan := range wuGeDaYans {
		if !daYan.IsLucky(filter_hard) {
			return false
		}
	}

	for _, sanCai := range sanCais {
		if !sanCai.getSanCaiFortune().IsLucky() {
			return false
		}
	}

	return true
}

//卦象吉祥
func (nk *NameStroke) IsGuaLucky(sex yi.Sex, filter_hard bool) bool {
	var guas []*yi.Yi = []*yi.Yi{}

	guas = append(guas, nk.BaGua())

	if filter_hard {
		guas = append(guas, nk.BaGuaSs()...)
		guas = append(guas, nk.BaGuaM1())
		guas = append(guas, nk.BaGuaM2())
	}

	for _, gua := range guas {
		if !gua.IsLucky(sex) {
			return false
		}
	}

	return true
}

func (nk *NameStroke) Clear(sex yi.Sex) {
	hash := nk.hash()

	ClearWuGeDaYan(*nk, sex)
	ClearWuGe(*nk)
	delete(guaList, hash)
}

var guaList map[int]*yi.Yi = map[int]*yi.Yi{}

var guaListLock sync.Mutex

// 周易八卦
// 姓总数取上卦，名总数取下卦，姓名总数取动爻
func (nk *NameStroke) BaGua() *yi.Yi {
	hash := nk.hash()

	sanCai := nk.getGuaSanCai()

	guaListLock.Lock()
	defer guaListLock.Unlock()

	gua := guaList[hash]

	if gua == nil {
		guaList[hash] = yi.NumberQiGua(sanCai.tian, sanCai.di, sanCai.ren)
	}

	return guaList[hash]
}

// 民国卦象1
// 名总数取上卦，姓名总数取下卦，姓名总数取动爻
func (nk *NameStroke) BaGuaM1() *yi.Yi {
	hash := nk.hash()

	sanCai := nk.getGuaSanCai()

	guaListLock.Lock()
	defer guaListLock.Unlock()

	gua := guaList[hash]

	if gua == nil {
		guaList[hash] = yi.NumberQiGua(sanCai.ren, sanCai.tian, sanCai.ren)
	}

	return guaList[hash]
}

// 民国卦象2
// 姓总数取上卦，姓名总数取下卦，姓名总数取动爻
func (nk *NameStroke) BaGuaM2() *yi.Yi {
	hash := nk.hash()

	sanCai := nk.getGuaSanCai()

	guaListLock.Lock()
	defer guaListLock.Unlock()

	gua := guaList[hash]

	if gua == nil {
		guaList[hash] = yi.NumberQiGua(sanCai.ren, sanCai.di, sanCai.ren)
	}

	return guaList[hash]
}

// 周神松卦象
// 姓名自下而上得三爻取下卦，自上而下得天人地三格，天格取上卦，天格取动爻
// 三字名自下而上由银行取卦，较为特殊
func (nk *NameStroke) BaGuaSs() []*yi.Yi {
	hash := nk.hash()

	sanCai := nk.getGuaSanCai()

	guas := []*yi.Yi{}

	guaListLock.Lock()
	defer guaListLock.Unlock()

	gua := guaList[hash]

	if gua == nil {

		guaList[hash] = yi.NumberQiGua(sanCai.ren, sanCai.tian, sanCai.tian)
		if nk.FirstStroke2 != 0 && nk.LastStroke2 == 0 {
			guas = append(guas, yi.NumberQiGua(yi.GetGua3Num(nk.FirstStroke2, nk.FirstStroke1, nk.LastStroke1), sanCai.tian, sanCai.ren))
		}
		guas = append(guas, yi.NumberQiGua(sanCai.ren, sanCai.tian, sanCai.tian))
	}

	return guas
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

func clearNameStroke(nk NameStroke) {
	hash := nk.hash()
	nameStroke := strokeCache[hash]

	if nameStroke != nil {
		delete(strokeCache, hash)
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
		l1 = f.lastChar[0].getStrokeScience(true)
		l2 = 0
	} else if lastCharLen == 2 {
		l1 = f.lastChar[0].getStrokeScience(true)
		l2 = f.lastChar[1].getStrokeScience(true)
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

			clearNameStroke(*nk)
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
		l1 = f.lastChar[0].getStrokeScience(true)
		l2 = 0
	} else if lastCharLen == 2 {
		l1 = f.lastChar[0].getStrokeScience(true)
		l2 = f.lastChar[1].getStrokeScience(true)
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
