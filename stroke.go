package fate

import (
	"errors"
	"github.com/godcong/fate/mongo"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
)

type Stroke struct {
	LastStroke  []int
	FirstStroke []int
}

//CalculatorBestStroke 计算最佳笔画数
func calculatorBestStroke(character []*mongo.Character) ([]*Stroke, error) {
	if character == nil || len(character) > 2 || len(character) == 0 {
		return nil, errors.New("CalculatorBestStroke: check character error")
	}

	l1, _ := strconv.Atoi(character[0].TotalStrokes)

	l2 := 0
	if len(character) == 2 {
		l2, _ = strconv.Atoi(character[1].TotalStrokes)
	}
	var wuXing mongo.WuXing
	f1, f2 := 1, 1
	for ; f1 < 32; f2++ {
		if f2 > 32 {
			f2 = 1
			f1++
		}

		wuGe := mongo.MakeWuGe(l1, l2, f1, f2)
		_ = checkWuGe(wuGe.TianGe)
		rg := checkWuGe(wuGe.RenGe)
		dg := checkWuGe(wuGe.DiGe)
		wg := checkWuGe(wuGe.WaiGe)
		zg := checkWuGe(wuGe.ZongGe)
		if !( rg && dg && wg && zg) {
			//log.Print(l1, l2, f1, f2, ":")
			//log.Println(tg, rg, dg, wg, zg)
			continue
		}

		sanCai := mongo.MakeSanCai(wuGe)
		mongo.C("wuxing").Find(bson.M{
			"wu_xing": []string{sanCai.TianCai, sanCai.RenCai, sanCai.DiCai},
		}).One(&wuXing)
		if wuXing.Fortune == "大吉" || wuXing.Fortune == "中吉" || wuXing.Fortune == "吉" {
			log.Println(l1, f1, f2, wuXing)
		}
	}

	return nil, nil
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
