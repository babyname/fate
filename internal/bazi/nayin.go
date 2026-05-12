package bazi

import v2 "github.com/godcong/chronos/v2"

type NaYin struct {
	calendar v2.Calendar
}

func NewNaYin(calendar v2.Calendar) *NaYin {
	return &NaYin{
		calendar: calendar,
	}
}
