package fate

import (
	"log"
	"testing"

	"github.com/godcong/fate/model"
)

func TestGenerateThreeTalent(t *testing.T) {
	for i := 1; i < 30; i++ {
		attr := model.GenerateAttribute(i)
		yy := model.GenerateYinYang(i)
		log.Println(i, attr, yy)
	}

}

//
//func TestThreeFive_InitSave(t *testing.T) {
//	for _, v := range name {
//		five := model.ThreeFive{
//			SurStrokes:    v[0],
//			SecondStrokes: v[1],
//			ThirdStrokes:  v[2],
//		}
//		five.InitSave()
//	}
//
//}
//func TestFindSecondStrokes(t *testing.T) {
//	s := model.FindSecondStrokes(17)
//	log.Println(s)
//}
//
//func TestFindThirdStrokes(t *testing.T) {
//	s := model.FindThirdStrokes(17, 8)
//	log.Println(s)
//}
