package fate

type Stroke struct {
	LastStroke  []int
	FirstStroke []int
	wuge        *WuGe
	sancai      *SanCai
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
