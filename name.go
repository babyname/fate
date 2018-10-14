package fate

import (
	"log"
	"os"
	"strings"

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
			name.cLast = append(name.cLast, CharacterFromName(v))
		}
		return name
	}
	name.LastName[0] = last
	for i, v := range name.LastName {
		name.cLast[i] = CharacterFromName(v)
	}

	return name
}

func CharacterFromName(s string) *model.Character {
	c := &model.Character{
		Base:          model.Base{},
		CharacterInfo: model.CharacterInfo{},
	}
	return c.Get()
}

func FilterBest(name Name) {
	file, err := os.OpenFile("output.txt", os.O_CREATE|os.O_APPEND|os.O_SYNC, os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()
	log.SetOutput(file)

	var fg FiveGrid
	for fmax, smax := 3, 1; fmax < 33; smax++ {

		if len(name.cLast) > 1 {
			fg = MakeFiveGridFromStrokes(name.cLast[0].ScienceStrokes, name.cLast[1].ScienceStrokes, fmax, smax)
		} else {
			fg = MakeFiveGridFromStrokes(name.cLast[0].ScienceStrokes, 0, fmax, smax)
		}
		tt := NewThreeTalent(fg)

		if GetProperty().UseFivePhase() {

		}
		fp := model.NewFivePhase(string(tt.SkyTalent.ThreeTalentAttribute), string(tt.PersonTalent.ThreeTalentAttribute), string(tt.LandTalent.ThreeTalentAttribute))
		f := fp.GetFortune()
		//if f == "大吉" || f == "中吉" || f == "吉" {
		if f == "大吉" {
			if fg.ContainBest("吉", "半吉") {
				//if fg.ContainBest("吉") {
				tt.PrintThreeTalent()
				var sec []model.Character
				var trd []model.Character
				model.CharacterList("", fmax, &sec)
				model.CharacterList("", smax, &trd)
				if sec == nil || trd == nil {
					continue
				}
				for _, v := range sec {
					log.Printf("二:%+v\n", v.CharacterInfo)
				}
				for _, v := range trd {
					log.Printf("三:%+v\n", v.CharacterInfo)
				}
			}

		}

		if smax >= 33 {
			smax = 0
			fmax++
		}

	}
}
