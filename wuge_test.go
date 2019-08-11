package fate_test

import (
	"github.com/godcong/fate"
	"log"
	"testing"
)

func TestWuGe_WaiGe(t *testing.T) {
	l1, l2, f1, f2 := 1, 1, 1, 1
	for i := 0; i < 80000; i++ {
		if f2 >= fate.WuGeMax {
			f1++
			f2 = 1
		}
		if f1 >= fate.WuGeMax {
			l2++
			f1 = 1
		}
		if l2 >= fate.WuGeMax {
			l1++
			l2 = 1
		}
		wg := fate.NewWuGe(l1, l2, f1, f2)
		sum := l1 + l2 + f1 + f2
		if wg.ZongGe() != sum {
			log.Println(wg.ZongGe() == sum, l1, l2, f1, f2, wg.ZongGe())
		}
		f2++
	}
}
