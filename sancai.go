package fate

import (
	"strings"

	"github.com/godcong/yi"
)

type SanCai struct {
	tian      int    //天格
	ren       int    //人格
	di        int    //地格
	tianRenDi string //三才
}

//SanCaiWuGe 三才五格
func GetSanCai(tian, ren, di int) SanCai {
	sanCaiNum := SanCai{
		tian: tian,
		ren:  ren,
		di:   di,
	}
	sanCaiNum.tianRenDi = strings.Join([]string{yi.NumberWuXing(tian), yi.NumberWuXing(ren), yi.NumberWuXing(di)}, "")

	return sanCaiNum
}

// SanCaiFortune ...
type SanCaiFortune struct {
	fortune string //吉凶
	comment string //说明
}

//key为tianRenDi
var sanCaiList map[string]*SanCaiFortune = make(map[string]*SanCaiFortune)

func init() {
	file_3wuxing, err := DataFiles.Open("data/3wuxing.csv")
	if err != nil {
		panic(err)
	}

	records, err := readData(file_3wuxing)

	if err != nil {
		panic(err)
	}

	for _, record := range records {
		sancai := SanCaiFortune{
			fortune: record[1],
			comment: record[2],
		}

		sanCaiList[record[0]] = &sancai
	}
}

// 由三才数获得描述
func (sc *SanCai) getSanCaiFortune() *SanCaiFortune {
	if sanCaiList[sc.tianRenDi] == nil {
		panic(sc.tianRenDi)
	}

	return sanCaiList[sc.tianRenDi]
}

//Check 检查三才吉凶
func (sc *SanCaiFortune) IsLucky() bool {
	return sc.fortune == "吉"
}
