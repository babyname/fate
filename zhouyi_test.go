package fate

import (
	"testing"
)

func TestJiaoHu(t *testing.T) {
	for i, j := 0, 0; i <= 7; j++ {
		t.Log(fu[i], fu[j], fu[hu(i, j)])
		t.Log(fu[i], fu[j], fu[jiao(i, j)])
		t.Log(fu[i], fu[cuo(i)])
		if j == 7 {
			i++
			j = 0
		}
	}
}

func TestCuo(t *testing.T) {
	for i := 0; i <= 7; i++ {
		t.Log(fu[i], fu[cuo(i)])

	}
}

func TestZong(t *testing.T) {
	for i, j := 0, 0; i <= 7; j++ {
		shang,xia:=zong(i, j)
		t.Log(fu[i], fu[j], fu[shang],fu[xia])
		if j == 7 {
			i++
			j = 0
		}
	}

}