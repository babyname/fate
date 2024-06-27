package fate

import "github.com/godcong/chronos"

// NaYin 纳音
type NaYin struct {
	calendar *chronos.Calendar
}

// NewNaYin 创建纳音
func NewNaYin(calendar *chronos.Calendar) *NaYin {
	return &NaYin{
		calendar: calendar,
	}
}
