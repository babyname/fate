package bazi

import v2 "github.com/godcong/chronos/v2"

// NaYin 表示纳音五行，基于日历计算干支对应的纳音属性。
type NaYin struct {
	calendar v2.Calendar
}

// NewNaYin 根据日历创建纳音实例。
func NewNaYin(calendar v2.Calendar) *NaYin {
	return &NaYin{
		calendar: calendar,
	}
}
