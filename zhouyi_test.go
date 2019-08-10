package fate

import "testing"

func TestQiGua(t *testing.T) {
	yi := QiGua(7, 7)
	get := yi.Get(0)
	t.Log(get)
}
