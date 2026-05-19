package zhouyi

import "testing"

func TestQiGua(t *testing.T) {
	yi := QiGua(7, 7)
	get := yi.GetGua(0)
	t.Log(get)
}
