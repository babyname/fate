package fate

import (
	"sync"

	ng "github.com/godcong/name_gua"
	"github.com/godcong/yi"
)

type GuaType = int

const (
	DefaultBaGua GuaType = iota
	Ss1BaGua
	Ss2BaGua
	M1BaGua
	M2BaGua
)

var guaList map[int]*yi.Yi = map[int]*yi.Yi{}

var guaListLock sync.Mutex

func (nk *NameStroke) BaGuaHash(gt GuaType) int {
	return nk.hash()<<3 + gt
}

// 周易八卦
// 姓总数取上卦，名总数取下卦，姓名总数取动爻
func (nk *NameStroke) BaGua() *yi.Yi {
	hash := nk.BaGuaHash(DefaultBaGua)

	guaSanCai := nk.getGuaSanCai()

	guaListLock.Lock()
	defer guaListLock.Unlock()

	gua := guaList[hash]

	if gua == nil {
		guaList[hash] = guaSanCai.BaGuaS1()
	}

	return guaList[hash]
}

// 民国卦象1
// 名总数取上卦，姓名总数取下卦，姓名总数取动爻
func (nk *NameStroke) BaGuaM1() *yi.Yi {
	hash := nk.BaGuaHash(M1BaGua)

	guaSanCai := nk.getGuaSanCai()

	guaListLock.Lock()
	defer guaListLock.Unlock()

	gua := guaList[hash]

	if gua == nil {
		guaList[hash] = guaSanCai.BaGuaM1()
	}

	return guaList[hash]
}

// 民国卦象2
// 姓总数取上卦，姓名总数取下卦，姓名总数取动爻
func (nk *NameStroke) BaGuaM2() *yi.Yi {
	hash := nk.BaGuaHash(M2BaGua)

	guaSanCai := nk.getGuaSanCai()

	guaListLock.Lock()
	defer guaListLock.Unlock()

	gua := guaList[hash]

	if gua == nil {
		guaList[hash] = guaSanCai.BaGuaM2()
	}

	return guaList[hash]
}

// 周神松卦象1
// 姓名自下而上得三爻取下卦，自上而下得天人地三格，天格取上卦，天格取动爻
// 三四字名自下而上取卦，较为特殊
func (nk *NameStroke) BaGuaSs1() []*yi.Yi {
	hash := nk.BaGuaHash(Ss1BaGua)

	guaSanCai := nk.getGuaSanCai()

	guas := []*yi.Yi{}

	guaListLock.Lock()
	defer guaListLock.Unlock()

	gua := guaList[hash]

	if gua == nil {

		if nk.FirstStroke2 != 0 {
			guaList[hash] = guaSanCai.BaGuaS2(nk.FirstStroke2, nk.FirstStroke1)
			guas = append(guas, guaList[hash])
		}

	}

	return guas
}

// 周神松卦象2
// 人格数取下卦，天格取上卦，天格取动爻
func (nk *NameStroke) BaGuaSs2() *yi.Yi {
	hash := nk.BaGuaHash(Ss2BaGua)

	guaSanCai := nk.getGuaSanCai()

	guaListLock.Lock()
	defer guaListLock.Unlock()

	gua := guaList[hash]

	if gua == nil {

		guaList[hash] = guaSanCai.BaGuaS3()
	}

	return guaList[hash]
}

//算卦象时有天人地，但是这个三才是不能用的，有很多姓氏依照五格三才是算不出吉祥的
func (nk *NameStroke) getGuaSanCai() *ng.GuaSanCai {
	tian := nk.FirstStroke1 + nk.FirstStroke2
	di := nk.LastStroke1 + nk.LastStroke2
	ren := tian + di
	return ng.GetGuaSanCai(tian, ren, di)
}

//卦象吉祥
func (nk *NameStroke) IsGuaLucky(sex yi.Sex, filter_hard bool) bool {
	var guas []*yi.Yi = []*yi.Yi{}

	guas = append(guas, nk.BaGua())

	if filter_hard {
		guas = append(guas, nk.BaGuaSs1()...)
		guas = append(guas, nk.BaGuaSs2())
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
