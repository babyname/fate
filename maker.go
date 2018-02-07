package fate

//
//func WheelStart(surname string) {
//	slist := make(map[int]*int)
//	c := model.Character{}
//	model.FindByNameChar(&c, surname)
//	s := model.FindSecondStrokesByCharacter(c)
//	var cs []model.Character
//	var nc []string
//	model.FindCharactersByStrokes(&cs, s)
//	for _, v := range cs {
//		nc = append(nc, v.NameChar)
//
//		if v.FixStrokes == 0 {
//			slist[v.FixStrokes] = nil
//			continue
//		}
//		slist[v.Strokes] = nil
//
//	}
//	stk := c.Strokes
//	if c.FixStrokes != 0 {
//		stk = c.FixStrokes
//	}
//	log.Println("second name:", nc)
//	var thirdStrokes []int
//	for key := range slist {
//		thirdStrokes = append(thirdStrokes, model.FindThirdStrokes(stk, key)...)
//
//	}
//
//	model.FindCharactersByStrokes(&cs, thirdStrokes)
//	nc = nil
//	for _, v := range cs {
//		nc = append(nc, v.NameChar)
//	}
//	log.Println("third name:", nc)
//}
//
//func WheelStartWithMass(sur, mass string) {
//	log.SetFlags(log.Llongfile)
//	c := model.Character{}
//	model.FindByNameChar(&c, sur)
//	//s := model.FindSecondStrokesByCharacter(c)
//	var cs []model.Character
//	var nc2 []string
//	var nc3 []string
//	var fives []model.Five
//	model.FindFiveWithFirstByMass(&fives, "金", []string{mass})
//	var tf []model.ThreeFive
//
//	for _, v := range fives {
//		log.Println(v.StringFive())
//		tf = model.FindStrokesWithFive(v)
//
//		var secondStk []int
//		var thirdStk []int
//		for _, value := range tf {
//			secondStk = append(secondStk, value.SecondStrokes)
//			thirdStk = append(thirdStk, value.ThirdStrokes)
//		}
//		model.FindCharactersWithFiveByStrokes(&cs, "土", secondStk)
//		nc2 = nil
//		for _, v := range cs {
//			nc2 = append(nc2, v.NameChar)
//		}
//		if nc2 == nil {
//			continue
//		}
//		log.Println(secondStk)
//		log.Println("second name", nc2)
//		model.FindCharactersByStrokes(&cs, thirdStk)
//		nc3 = nil
//		for _, v := range cs {
//			nc3 = append(nc3, v.NameChar)
//		}
//		log.Println(thirdStk)
//		log.Println("third name", nc3)
//	}
//
//}

//func SecondNameList(sur string) []string {
//	c := model.Character{}
//	model.FindByNameChar(&c, sur)
//	if c.FixStrokes == 0 {
//		var cs []model.Character
//		model.FindCharactersWithFiveByStrokes(&cs, "", []int{c.Strokes})
//	}
//}


//
//func TurnIntoFiveGrid(s int) string {
//
//}