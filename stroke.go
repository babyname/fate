package fate

type Stroke struct {
	LastStroke  []int
	FirstStroke []int
	wuge        *WuGe
	sancai      *SanCai
}

//First 名字笔画
func First(s []*Stroke, i int) []int {
	var rlt []int
	for idx := range s {
		if len(s[idx].FirstStroke) > i {
			rlt = append(rlt, s[idx].FirstStroke[i])
		}

	}
	return rlt
}

//Last 姓氏笔画
func Last(s []*Stroke, i int) []int {
	var rlt []int
	for idx := range s {
		if len(s[idx].LastStroke) > i {
			rlt = append(rlt, s[idx].LastStroke[i])
		}
	}
	return rlt
}
