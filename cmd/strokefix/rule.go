package main

import (
	"github.com/godcong/fate"
	"strings"
)

var charCharList = []string{
	`阧
	障
	邓
	陇
	陾
	階
	鄁
	阾
	降
	郸
	酃
	郎
	鄐
	隄
	隓
	邯
	鄟
	郐
	邝
	陞
	那
	邤
	隂
	酆
	陈
	酇
	鄴
	隯
	郲
	鄽
	郩
	陖
	隗
	陠
	阱
	郔
	陲
	阼
	都
	郜
	鄻
	邦
	郁
	郋
	邧
	隩
	阬
	陳
	郘
	鄺
	陌
	隀
	鄋
	酄
	郵
	阠
	郦
	郣
	郐
	郻
	陿
	郡
	鄷
	隉
	陉
	鄲
	隥
	鄇
	陧
	都
	陁
	鄂
	陱
	鄥
	阶
	陂
	陽
	陼
	鄧
	陪
	酆
	鄕
	陗
	邒
	郖
	鄭
	鄫
	鄖
	郉
	陊
	郭
	郊
	邗
	阪
	郇
	郳
	郊
	邓
	郏
	阞
	鄣
	陑
	郷
	隭
	鄊
	鄆
	酄
	阬
	邞
	邿
	郅
	酂
	陆
	鄶
	隡
	陫
	鄮
	鄛
	邹
	郢
	邻
	鄷
	郆
	陦
	鄗
	陂
	鄏
	阳
	郗
	防
	邱
	鄒
	阽
	陔
	邬
	郛
	郣
	鄐
	隄
	郄
	郂
	防
	隕
	郴
	鄸
	鄈
	郦
	邳
	陨
	邨
	陇
	鄲
	郴
	陑
	陚
	酁
	郡
	邴
	隬
	邵
	陉
	邥
	部
	隣
	陕
	邴
	隝
	郧
	阺
	郞
	阫
	郏
	郟
	郰
	郙
	邛
	隢
	邘
	邺
	陹
	鄿
	邟
	鄍
	郹
	限
	附
	陰
	郖
	郍
	阮
	隚
	隥
	阪
	阭
	鄘
	陏
	鄶
	陳
	陜
	邭
	陸
	阦
	阷
	阞
	邸
	郑
	鄙
	隫
	鄺
	隞
	隘
	陒
	鄗
	隴
	鄉
	隍
	郀
	邠
	陶
	隀
	郬
	邼
	隔
	陀
	邖
	阤
	鄄
	鄤
	邼
	隤
	陬
	郹
	隨
	除
	陲
	邽
	降
	陮
	鄩
	阽
	邔
	阩
	阴
	鄰
	陋
	鄇
	队
	隳
	隊
	邟
	郝
	邲
	郃
	郯
	邝
	郜
	鄑
	郝
	隊
	郀
	隭
	陮
	隋
	邭
	郞
	鄜
	邙
	隝
	郪
	隴
	邷
	陘
	鄝
	郲
	阿
	邾
	郈
	阥
	郈
	郙
	鄌
	郸
	隆
	鄄
	鄼
	陡
	邰
	陏
	鄛
	陫
	邮
	郮
	邨
	郠
	陙
	陈
	邦
	鄁
	鄵
	隧
	郱
	阸
	鄈
	郓
	邜
	陯
	隑
	隍
	陡
	陛
	阸
	邸
	鄤
	邯
	隔
	邢
	郟
	鄯
	阻
	隑
	鄧
	邗
	隯
	邩
	郤
	陘
	郕
	邚
	郠
	陃
	陊
	郫
	除
	鄀
	陖
	隘
	郂
	隫
	鄝
	阨
	邩
	酈
	郚
	鄹
	鄚
	鄱
	邞
	阨
	隌
	邽
	陔
	部
	阧
	陵
	鄽
	附
	陟
	隞
	邲
	邠
	队
	鄳
	郕
	随
	隌
	郿
	邡
	郥
	陙
	邶
	邪
	鄡
	陃
	郭
	邶
	邡`,
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
