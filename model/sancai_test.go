package model_test

import (
	"testing"

	"log"

	"github.com/godcong/fate/model"
)

var name = [][]int{
	{17, 8, 7},
	{17, 8, 10},
	{17, 8, 16},
	{17, 12, 6},
	{17, 18, 6},
	{17, 18, 17},

	//{2, 9, 4},
	//2 3 20
	//2 1 10
	//2 9 6
	//2 4 2
	//2 1 20
	//2 9 7
	//2 4 12
	//2 3 10
	//3 2 6
	//3 2 16
	//3 3 5
	//3 3 10
	//3 18 14
	//3 8 24
	//3 3 12
	//3 18 17
	//3 10 22
	//4 2 5
	//4 2 16
	//4 3 22
	//4 9 2
	//4 19 22
	//4 19 6
	//4 11 6
	//4 9 4
	//4 20 5
	//5 18 14
	//5 10 14
	//5 2 4
	//5 11 7
	//5 2 14
	//5 12 14
	//5 8 5
	//5 18 6
	//5 8 24
	//6 19 16
	//6 10 15
	//6 9 6
	//6 10 23
	//6 9 14
	//6 19 4
	//6 9 16
	//6 19 6
	//6 10 7
	//7 22 10
	//7 9 15
	//7 8 10
	//7 9 16
	//7 8 16
	//7 18 6
	//7 8 17
	//7 18 7
	//7 9 7
	//8 13 2
	//8 9 16
	//8 3 2
	//8 13 12
	//8 10 5
	//8 3 12
	//8 13 16
	//8 10 6
	//8 9 6
	//8 10 15
	//8 9 7
	//9 9 7
	//9 2 4
	//9 12 4
	//9 2 14
	//9 12 20
	//9 8 7
	//9 20 12
	//9 9 6
	//9 22 10
	//10 11 14
	//10 3 22
	//10 3 2
	//10 11 20
	//10 11 4
	//10 3 10
	//10 13 2
	//10 11 10
	//10 3 12
	//10 13 10
	//10 11 12
	//10 3 20
	//10 19 12
	//10 13 12
	//10 14 7
	//10 14 17
	//10 19 2
}

func TestThreeFive_InitSave(t *testing.T) {
	for _, v := range name {
		five := model.ThreeFive{
			SurStrokes:    v[0],
			SecondStrokes: v[1],
			ThirdStrokes:  v[2],
		}
		five.InitSave()
	}

}
func TestFindSecondStrokes(t *testing.T) {
	s := model.FindSecondStrokes(17)
	log.Println(s)
}

func TestFindThirdStrokes(t *testing.T) {
	s := model.FindThirdStrokes(17, 8)
	log.Println(s)
}
