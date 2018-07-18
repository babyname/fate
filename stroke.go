package fate

import (
	"fmt"

	"github.com/godcong/fate/mongo"
)

//CalculatorBestStroke 计算最佳笔画数
func CalculatorBestStroke(character []*mongo.Character) error {
	if len(character) > 2 {
		return fmt.Errorf("CalculatorBestStroke: len overflow", )
	}
	return nil
}


