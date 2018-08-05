package fate

import (
	"github.com/godcong/fate/mongo"
)

type Stroke struct {
	LastStroke  []int
	FirstStroke []int
	wuge        *WuGe
	sancai      *SanCai
}

func checkWuGe(dy []*mongo.DaYan, idx int) bool {
	if dy == nil || len(dy) <= idx {
		return false
	}
	switch dy[idx-1].Fortune {
	case "吉", "半吉":
		return true
	}
	return false
}

func First(s []*Stroke, i int) []int {
	var rlt []int
	for idx := range s {
		if len(s[idx].FirstStroke) > i {
			rlt = append(rlt, s[idx].FirstStroke[i])
		}

	}
	return rlt
}

func Last(s []*Stroke, i int) []int {
	var rlt []int
	for idx := range s {
		if len(s[idx].LastStroke) > i {
			rlt = append(rlt, s[idx].LastStroke[i])
		}
	}
	return rlt
}
