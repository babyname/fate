package main

import (
	"github.com/godcong/fate"
	"strings"
)

var charCharList = []string{
	`阧障陇陾階阾降隄隓陞隂陈隯陖隗陠阱陲阼隩阬陳陌隀阠陿隉陉隥陧陁陱阶陂陽陼陪陗陊阪阞陑隭阬陆隡陫陦陂阳防阽陔隄防隕陨陇陑陚隬陉隣陕隝阺阫隢陹限附陰阮隚隥阪阭陏陳陜陸阦阷阞隫隞隘陒隴隍陶隀隔陀阤隤陬隨除陲邽降陮鄩阽邔阩阴鄰陋鄇队隳隊邟郝邲郃郯邝郜鄑郝隊郀隭陮隋邭郞鄜邙隝郪隴邷陘鄝郲阿邾郈阥郈郙鄌郸隆鄄鄼陡邰陏鄛陫邮郮邨郠陙陈邦鄁鄵隧郱阸鄈郓邜陯隑隍陡陛阸邸鄤邯隔邢郟鄯阻隑鄧邗隯邩郤陘郕邚郠陃陊郫除鄀陖隘郂隫鄝阨邩酈郚鄹鄚鄱邞阨隌邽陔部阧陵鄽附陟隞邲邠队鄳郕随隌郿邡郥陙邶邪鄡陃`,
	`邤酆那鄟郐邝邯酃郎鄐郸鄁邓酇郲鄽郩鄴郔郭郜鄻邦鄺郘鄋郦郡鄷鄇都鄲鄂鄥郣郐郻鄧鄕邒郭郳鄣郷鄆邞鄶邿鄮鄗鄏郗鄒邬郛郄郂鄸邨郴鄲邴邵邥部邴郧郞郟邘邟鄍郹郖郍鄿邺郰鄘鄶郙邭邛郏郡鄙鄺郑邸鄗鄉邠郀郬邼邖邼鄄郹鄤酁鄈郦邳郴郣鄐邱鄛邹郢邻鄷郆郅酂酄鄊郊邓郏郇郊邗郖鄭鄫鄖郉酆酄郵郁郋邧都邶邡`,
}

var numberCharList = `一二三四五六七八九十`

func NumberChar(ch *fate.Character) bool {
	if i := strings.Index(numberCharList, ch.Ch); i != -1 {
		if ch.Stroke != 0 {
			ch.Stroke += i
			return true
		}
		if ch.SimpleTotalStroke != 0 {
			ch.SimpleTotalStroke += i
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
