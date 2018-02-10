package fate_test

import (
	"testing"

	"github.com/godcong/fate"
	"github.com/godcong/fate/debug"
	"github.com/godcong/fate/model"
)

func TestNewThreeTalent(t *testing.T) {
	for i, j, k := 17, 1, 1; j < 33 && k < 33; k++ {

		fg := fate.MakeFiveGridFromStrokes(i, 0, j, k)

		tt := fate.NewThreeTalent(fg)
		fp := model.NewFivePhase(string(tt.SkyTalent.ThreeTalentAttribute), string(tt.PersonTalent.ThreeTalentAttribute), string(tt.LandTalent.ThreeTalentAttribute))
		f := fp.GetFortune()
		if f == "大吉" {
			if fg.PrintBigYan(true) {
				debug.Println("笔画:", i, j, k)
				tt.PrintThreeTalent()
				debug.Println(f)
			}

		}

		if k == 32 {
			k = 1
			j++
		}
	}
	//fg := fate.MakeFiveGridFromStrokes(17, 0, 8, 7)
	//fg.PrintBigYan(false)
	//tt := fate.NewThreeTalent(fg)
	//tt.PrintThreeTalent()
	//fp := model.NewFivePhase(string(tt.SkyTalent.ThreeTalentAttribute), string(tt.PersonTalent.ThreeTalentAttribute), string(tt.LandTalent.ThreeTalentAttribute))
	//f := fp.GetFortune()
	//debug.Println(f)
}

func TestNewName(t *testing.T) {
	fate.NewName("李")
}

func TestFilterBest(t *testing.T) {
	fate.FilterBest(fate.NewName("蒋"))
}
