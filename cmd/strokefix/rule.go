package main

import (
	"github.com/godcong/fate"
	"strings"
)

var charCharList = []string{
	`阧障陇陾階阾降隄隓鄟郐邝陞那邤隂酆陈酇鄴隯郲鄽郩陖隗陠阱郔陲阼都郜鄻邦郁郋邧隩阬陳郘鄺陌隀鄋酄郵阠郦郣郐郻陿郡鄷隉陉鄲隥鄇陧都陁鄂陱鄥阶陂陽陼鄧陪酆鄕陗邒郖鄭鄫鄖郉陊郭郊邗阪郇郳郊邓郏阞鄣陑郷隭鄊鄆酄阬邞邿郅酂陆鄶隡陫鄮鄛邹郢邻鄷郆陦鄗陂鄏阳郗防邱鄒阽陔邬郛郣鄐隄郄郂防隕郴鄸鄈郦邳陨邨陇鄲郴陑陚酁郡邴隬邵陉邥部隣陕邴隝郧阺郞阫郏郟郰郙邛隢邘邺陹鄿邟鄍郹限附陰郖郍阮隚隥阪阭鄘陏鄶陳陜邭陸阦阷阞邸郑鄙隫鄺隞隘陒鄗隴鄉隍郀邠陶隀郬邼隔陀邖阤鄄鄤邼隤陬郹隨除陲邽降陮鄩阽邔阩阴鄰陋鄇队隳隊邟郝邲郃郯邝郜鄑郝隊郀隭陮隋邭郞鄜邙隝郪隴邷陘鄝郲阿邾郈阥郈郙鄌郸隆鄄鄼陡邰陏鄛陫邮郮邨郠陙陈邦鄁鄵隧郱阸鄈郓邜陯隑隍陡陛阸邸鄤邯隔邢郟鄯阻隑鄧邗隯邩郤陘郕邚郠陃陊郫除鄀陖隘郂隫鄝阨邩酈郚鄹鄚鄱邞阨隌邽陔部阧陵鄽附陟隞邲邠队鄳郕随隌郿邡郥陙邶邪鄡陃`,
	`邯酃郎鄐郸鄁邓郭邶邡`,
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
