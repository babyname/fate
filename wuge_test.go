package fate_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/godcong/fate"
)

// TestWuGe_WaiGe ...
func TestWuGe_WaiGe(t *testing.T) {
	l1, l2, f1, f2 := 1, 1, 1, 1
	for i := 0; i < 80000; i++ {
		if f2 >= fate.BiHuaMax {
			f1++
			f2 = 1
		}
		if f1 >= fate.BiHuaMax {
			l2++
			f1 = 1
		}
		if l2 >= fate.BiHuaMax {
			l1++
			l2 = 1
		}
		nk := fate.GetNameStroke(l1, l2, f1, f2)
		wg := nk.GetWuGes(true)[0]
		sum := l1 + l2 + f1 + f2
		if wg.ZongGe != sum {
			log.Println(wg.ZongGe == sum, l1, l2, f1, f2, wg.ZongGe)
		}
		fmt.Println("result:", fate.Check(wg))
		f2++
	}
}
