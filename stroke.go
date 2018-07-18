package fate

import (
	"github.com/godcong/fate/mongo"
)

type Stroke struct {
	LastStroke  []int
	FirstStroke []int
}

//CalculatorBestStroke 计算最佳笔画数
func CalculatorBestStroke(character []*mongo.Character) (error, []*Stroke) {
	if len(character) > 2 {
		//return fmt.Errorf("CalculatorBestStroke: len overflow", )
	}
	return nil, nil
}
