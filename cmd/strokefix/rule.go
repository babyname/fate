package main

import (
	"github.com/godcong/fate"
	"strings"
)

var numberCharList = `一二三四五六七八九十`

func NumberChar(ch *fate.Character) bool {
	if i := strings.Index(numberCharList, ch.Ch); i != -1 {
		return true
	}
	return false
}
