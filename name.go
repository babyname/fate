package fate

import (
	"strings"

	"github.com/godcong/fate/debug"
	"github.com/godcong/fate/model"
)

type Name struct {
	FirstName []string
	cFirst    []*model.Character
	LastName  []string
	cLast     []*model.Character
}

func NewName(last string) Name {
	name := Name{}
	if len(last) > 1 {
		name.LastName = strings.Split(last, "")
		for _, v := range name.LastName {
			name.cFirst = append(name.cFirst, CharacterFromName(v))
		}
		return name
	}
	name.LastName[0] = last
	for i, v := range name.LastName {
		name.cFirst[i] = CharacterFromName(v)
	}
	return name
}

func CharacterFromName(s string) *model.Character {
	c := &model.Character{
		IsSur:          false,
		SimpleChar:     s,
		SimpleStrokes:  0,
		TradChar:       "",
		TradStrokes:    0,
		NameType:       "",
		NameRoot:       "",
		Radical:        "",
		ScienceStrokes: 0,
		Pinyin:         "",
		Comment:        "",
	}
	return c.Get()
}

func FilterBest(name Name) {
	var fg FiveGrid
	for fmax, smax := 1, 1; fmax < 33; smax++ {
		if len(name.cLast) > 1 {
			fg = MakeFiveGridFromStrokes(name.cLast[0].ScienceStrokes, name.cLast[1].ScienceStrokes, fmax, smax)
		} else {
			fg = MakeFiveGridFromStrokes(name.cLast[0].ScienceStrokes, 0, fmax, smax)
		}
		tt := NewThreeTalent(fg)

		fp := model.NewFivePhase(string(tt.SkyTalent.ThreeTalentAttribute), string(tt.PersonTalent.ThreeTalentAttribute), string(tt.LandTalent.ThreeTalentAttribute))
		f := fp.GetFortune()
		if f == "大吉" {
			if fg.ContainBest("大吉", "半吉") {

				tt.PrintThreeTalent()
				debug.Println(f)
			}

		}

		if smax >= 33 {
			smax = 0
			fmax++
		}

	}

	//for i, j, k := 17, 1, 1; j < 33 && k < 33; k++ {
	//
	//	fg := fate.MakeFiveGridFromStrokes(i, 0, j, k)
	//
	//	tt := fate.NewThreeTalent(fg)
	//	fp := model.NewFivePhase(string(tt.SkyTalent.ThreeTalentAttribute), string(tt.PersonTalent.ThreeTalentAttribute), string(tt.LandTalent.ThreeTalentAttribute))
	//	f := fp.GetFortune()
	//	if f == "大吉" {
	//		if fg.PrintBigYan(true) {
	//			debug.Println("笔画:", i, j, k)
	//			tt.PrintThreeTalent()
	//			debug.Println(f)
	//		}
	//
	//	}
	//
	//	if k == 32 {
	//		k = 1
	//		j++
	//	}
	//}
}
