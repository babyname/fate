package fate

import (
	"errors"
	"github.com/godcong/fate/mongo"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"log"
)

type Stroke struct {
	LastStroke  []int
	FirstStroke []int
}

//CalculatorBestStroke 计算最佳笔画数
func calculatorBestStroke(f *fate, character []*mongo.Character) ([]*Stroke, error) {
	var strokes []*Stroke
	martial := f.GetMartial()
	if character == nil || len(character) > 2 || len(character) == 0 {
		return nil, errors.New("CalculatorBestStroke: check character error")
	}
	log.Println(character)
	l1, _ := strconv.Atoi(character[0].TotalStrokes)

	l2 := 0
	if len(character) == 2 {
		l2, _ = strconv.Atoi(character[1].TotalStrokes)
	}
	var wuXing mongo.WuXing
	f1, f2 := 1, 1
	for ; f1 < 30; f2++ {
		if f2 > 30 {
			f2 = 1
			f1++
		}

		wuGe := MakeWuGe(l1, l2, f1, f2)
		_ = checkWuGe(wuGe.TianGe)
		rg := checkWuGe(wuGe.RenGe)
		dg := checkWuGe(wuGe.DiGe)
		wg := checkWuGe(wuGe.WaiGe)
		zg := checkWuGe(wuGe.ZongGe)

		if martial.BiHua {
			if !(rg && dg && wg && zg) {
				continue
			}
		}

		sanCai := MakeSanCai(wuGe)
		mongo.C("wuxing").Find(bson.M{
			"wu_xing": []string{sanCai.TianCai, sanCai.RenCai, sanCai.DiCai},
		}).One(&wuXing)

		if martial.SanCai {
			if wuXing.Fortune == "大吉" || wuXing.Fortune == "中吉" || wuXing.Fortune == "吉" {
				strokes = append(strokes, &Stroke{
					LastStroke:  []int{l1, l2},
					FirstStroke: []int{f1, f2},
				})
				log.Println(l1, l2, f1, f2, wuXing)
			}
		}
	}

	return strokes, nil
}

func checkWuGe(i int) bool {
	var dy mongo.DaYan
	mongo.C("dayan").Find(bson.M{
		"index": i,
	}).One(&dy)
	if !(dy.Fortune == "吉" || dy.Fortune == "半吉") {
		return false
	}

	return true

}
