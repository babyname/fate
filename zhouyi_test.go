package fate

import (
	"testing"
)

func TestJiaoHu(t *testing.T) {
	for i, j := 0, 0; i <= 7; j++ {
		t.Log(fu[i], fu[j], fu[hu(i, j)])
		if j == 7 {
			i++
			j = 0
		}
	}
}
