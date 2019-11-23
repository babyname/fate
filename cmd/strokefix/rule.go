package main

import (
	"github.com/godcong/fate"
	"strings"
)

var numberCharList = `一二三四五六七八九十`

func NumberChar(ch *fate.Character) bool {
	if i := strings.Index(numberCharList, ch.Ch); i != -1 {
		if ch.Stroke != 0 {
			ch.Stroke += i
			return true
		}

	}
	return false
}

var radicalCharList = map[string]int{
	"扌": 1,
	"忄": 1,
	"氵": 1,
	"犭": 1,
	"礻": 1,
	"王": 1,
	"艹": 3,
	"衤": 1,
	"月": 3,
	"辶": 4,
}

func RadicalChar(ch *fate.Character) bool {
	for k, v := range radicalCharList {
		if strings.Compare(ch.Radical, k) == 0 {
			ch.Stroke += v
			return true
		}
		if strings.Compare(ch.SimpleRadical, k) == 0 {
			ch.SimpleTotalStroke += v
			return true
		}
	}
	return false
}

var fixTable = []func(character *fate.Character) bool{
	RadicalChar,
	NumberChar,
}

func fixChar(character *fate.Character) {
	for _, f := range fixTable {
		if f(character) {
			return
		}
	}
}
