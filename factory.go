package fate

import (
	"github.com/godcong/fate/model"
)

type factory struct {
	surname string
	sur     model.Character
	isFive  bool
	//fives        []model.Five
	//threeFives   []model.ThreeFive
	//curThreeFive model.ThreeFive
	second string
	sec    model.Character
	third  string
	trd    model.Character
}

//
//func NewFactory(sur string) *factory {
//	f := new(factory)
//	f.surname = sur
//	model.FindByNameChar(&f.sur, f.surname)
//	return f
//}
//
//func (f *factory) LoadThreeFive(isFive int) {
//	f.isFive = false
//	var tf []model.ThreeFive
//	if isFive == 1 {
//		f.isFive = true
//		model.FindFiveWithFirstByMass(&f.fives, model.MakeSanCai(f.sur.FixStrokes+1), []string{"优"})
//
//		for _, value := range f.fives {
//			tf = append(tf, model.FindStrokesWithFive(value)...)
//		}
//
//	} else {
//		tf = model.FindSecondThreeFive(f.sur.FixStrokes)
//	}
//	f.threeFives = tf
//
//}
//
//func (f *factory) SecondName(zodiac model.Zodiac) {
//	var cs []model.Character
//	var sec map[int]model.Character
//
//	best := strings.Split(zodiac.Best, ",")
//	for _, value := range f.threeFives {
//		var i int
//		if value.SecondStrokes == 0 {
//			continue
//		}
//		if zodiac.Best == "" {
//			model.FindCharactersByStrokeBest(&cs, value.SecondStrokes, nil)
//		} else {
//			model.FindCharactersByStrokeBest(&cs, value.SecondStrokes, best)
//		}
//
//		sec = make(map[int]model.Character)
//		for idx, c := range cs {
//			sec[idx+1] = c
//		}
//		NameOutput(sec)
//		fmt.Println("按回车或者0继续检索，按序列确认")
//		fmt.Scanln(&i)
//		if i != 0 {
//			f.second = sec[i].NameChar
//			f.sec = sec[i]
//			f.curThreeFive = value
//			break
//		}
//	}
//}
//func (f *factory) ThirdName(zodiac model.Zodiac) {
//	var i int
//	var cs []model.Character
//	var trd map[int]model.Character
//	best := strings.Split(zodiac.Best, ",")
//	//model.FindCharactersByStroke(&cs, f.curThreeFive.ThirdStrokes)
//	if f.curThreeFive.ThirdStrokes == 0 {
//		return
//	}
//
//	model.FindCharactersByStrokeBest(&cs, f.curThreeFive.ThirdStrokes, best)
//	trd = make(map[int]model.Character)
//	for idx, c := range cs {
//		trd[idx+1] = c
//	}
//	NameOutput(trd)
//	fmt.Println("按序列确认，按0跳过")
//	fmt.Scanln(&i)
//	if i != 0 {
//		f.third = trd[i].NameChar
//		f.trd = trd[i]
//	}
//}
//
//func NameOutput(cs map[int]model.Character) {
//	idx := len(cs)
//	fmt.Println("共计：", idx)
//	for i := 1; i <= idx; i++ {
//		fmt.Print(i, ":", cs[i].NameChar, "(", cs[i].NameType, ")", "  ")
//	}
//	fmt.Println()
//}
//
//func (f *factory) GetName() string {
//	log.Println(f.sur.FixStrokes, f.sec.FixStrokes, f.trd.FixStrokes)
//	tf := model.NewThreeFive(f.sur.FixStrokes, f.sec.FixStrokes, f.trd.FixStrokes)
//	tf.PrintString()
//	log.Println(f.sur.Pinyin, f.sec.Pinyin, f.trd.Pinyin)
//	return f.surname + f.second + f.third
//}
