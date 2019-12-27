package main

import (
	"github.com/godcong/fate"
	"strings"
)

var charCharList = []string{
	`阧障陇陾階阾降隄隓陞隂陈隯陖隗陠阱陲阼隩阬陳陌隀阠陿隉陉隥陧陁陱阶陂陽陼陪陗陊阪阞陑隭阬陆隡陫陦陂阳防阽陔隄防隕陨陇陑陚隬陉隣陕隝` +
		`阺阫隢陹限附陰阮隚隥阪阭陏陳陜陸阦阷阞隫隞隘陒隴隍陶隀隔陀阤隤陬隨除陲降陮阽阩阴陋队隳隊隊隭陮隋隝隴陘阥隆陡陏陫陙陈隧阸陯隑隍` +
		`陡陛阸隔阻隑隯陘陃陊除陖隘隫阨阨隌陔阧陵附陟隞队随隌陙陃`,
	`邤酆那鄟郐邝邯酃邸邢郟鄧郤郕邚郫鄀郠郂鄝邩酈邽部郚鄽鄹鄚邲鄳郕郥邶邪鄡邡郿邠鄱邞邩邗鄯鄤邯郎鄐郸鄁邓酇郲鄽郩鄴郔郭郜鄻邦鄺郘鄋郦郡` +
		`鄷鄇都鄲鄂鄥郣郐郻鄧鄕邒郭郳鄣郷鄆邞鄶邿鄮鄗鄏郗鄒邬郪邷郲郈鄼鄄郙邰鄛邦郱郓邜鄈鄁鄵邮郮邨郠鄌郸阿邾郈鄝郛郄郂鄸邨郴鄲邴邵邥部` +
		`邴郧郞郟邘邟鄍郹郖郍鄿邺郰鄘鄶郙邭邛郏郡鄙鄺郑邸鄗鄉邠郀郬邼邖邼鄄郹鄤酁邽鄈鄩邔郦鄰鄇邳郴郝郀邲郃郞鄜邙邭郯邝郜鄑郝邟郣鄐邱鄛` +
		`邹郢邻鄷郆郅酂酄鄊郊邓郏郇郊邗郖鄭鄫鄖郉酆酄郵郁郋邧都邶邡`,
}

func CharChar(ch *fate.Character) bool {
	if i := strings.Index(charCharList[0], ch.Ch); i != -1 {
		if ch.Stroke != 0 {
			ch.ScienceStroke = ch.Stroke + 8
			return true
		}
		if ch.SimpleTotalStroke != 0 {
			ch.ScienceStroke = ch.SimpleTotalStroke + 8
			return true
		}
	}

	return false
}

var numberCharList = `一二三四五六七八九十`

func NumberChar(ch *fate.Character) bool {
	if i := strings.Index(numberCharList, ch.Ch); i != -1 {
		if ch.Stroke != 0 {
			ch.ScienceStroke = ch.Stroke + i
			return true
		}
		if ch.SimpleTotalStroke != 0 {
			ch.ScienceStroke = ch.SimpleTotalStroke + i
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
			ch.ScienceStroke = ch.Stroke + v
			return true
		}
		if strings.Compare(ch.SimpleRadical, k) == 0 {
			ch.ScienceStroke = ch.SimpleTotalStroke + v
			return true
		}
	}
	return false
}

var fixTable = []func(character *fate.Character) bool{
	RadicalChar,
	NumberChar,
	CharChar,
}

func fixChar(character *fate.Character) bool {
	for _, f := range fixTable {
		if f(character) {
			return true
		}
	}
	return false
}
