package fate

import (
	"log"

	"github.com/godcong/fate/model"
)

func WheelStart(surname string) {
	slist := make(map[int]*int)
	c := model.Character{}
	model.FindByNameChar(&c, surname)
	s := model.FindSecondStrokesByCharacter(c)
	var cs []model.Character
	var nc []string
	model.FindCharactersByStrokes(&cs, s)
	for _, v := range cs {
		nc = append(nc, v.NameChar)

		if v.FixStrokes == 0 {
			slist[v.FixStrokes] = nil
			continue
		}
		slist[v.Strokes] = nil

	}
	stk := c.Strokes
	if c.FixStrokes != 0 {
		stk = c.FixStrokes
	}
	log.Println("second name:", nc)
	var thirdStrokes []int
	for key := range slist {
		thirdStrokes = append(thirdStrokes, model.FindThirdStrokes(stk, key)...)

	}

	model.FindCharactersByStrokes(&cs, thirdStrokes)
	nc = nil
	for _, v := range cs {
		nc = append(nc, v.NameChar)
	}
	log.Println("third name:", nc)
}

func WheelStartWithMass(sur, mass string) {
	c := model.Character{}
	model.FindByNameChar(&c, sur)
	//s := model.FindSecondStrokesByCharacter(c)
	var cs []model.Character
	var nc2 []string
	var nc3 []string
	var fives []model.Five
	model.FindFiveWithFirstByMass(&fives, "金", mass)
	var tf []model.ThreeFive
	log.Println(fives)
	for _, v := range fives {
		tf = model.FindStrokesWithFive(v)
		log.Println(tf)
		var secondStk []int
		var thirdStk []int
		for _, value := range tf {
			secondStk = append(secondStk, value.SecondStrokes)
			thirdStk = append(thirdStk, value.ThirdStrokes)
		}
		log.Println(secondStk)
		model.FindCharactersByStrokes(&cs, secondStk)
		nc2 = nil
		for _, v := range cs {
			nc2 = append(nc2, v.NameChar)
		}
		log.Println("second name", nc2)
		log.Println(thirdStk)
		model.FindCharactersByStrokes(&cs, thirdStk)
		nc3 = nil
		for _, v := range cs {
			nc3 = append(nc3, v.NameChar)
		}
		log.Println("third name", nc3)
	}

}
