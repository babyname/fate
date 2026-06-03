// Package analysis provides name analysis and report generation.
package analysis

import (
	"github.com/babyname/fate/v4/ent"
)

// shengMap defines the wuxing sheng (generating) cycle:
// each element generates the next one in the cycle.
var shengMap = map[string]string{
	"木": "火", "火": "土", "土": "金", "金": "水", "水": "木",
}

// keMap defines the wuxing ke (overcoming) cycle:
// each element overcomes another element in the cycle.
var keMap = map[string]string{
	"木": "土", "土": "水", "水": "火", "火": "金", "金": "木",
}

// YinYangWuXingAttr returns the yin-yang and wuxing attribute string
// derived from the given stroke number. The wuxing is determined by
// stroke % 10 and the yin-yang by stroke % 2.
func yinYangWuXingAttr(stroke int) string {
	sanCaiStr := "水木木火火土土金金水"
	yinYangStr := "阴阳"
	wx := string([]rune(sanCaiStr)[stroke%10])
	yy := string([]rune(yinYangStr)[stroke%2])
	return yy + wx
}

// GetTraditionalChar returns the traditional Chinese character for the
// given Character. If the character is already traditional it is returned
// directly; otherwise the simplified-to-traditional mapping is consulted.
func getTraditionalChar(c *ent.Character) string {
	if c.IsTraditional {
		return c.Char
	}
	if trad, ok := simplifiedToTraditional[c.Char]; ok {
		return trad
	}
	return c.Char
}

// TruncateStr truncates the string s to at most maxLen runes.
// If truncation occurs, an ellipsis "…" is appended.
func truncateStr(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "…"
}

// IsSheng reports whether wuxing element a sheng (generates) element b.
func isSheng(a, b string) bool {
	if v, ok := shengMap[a]; ok {
		return v == b
	}
	return false
}

// IsKe reports whether wuxing element a ke (overcomes) element b.
func isKe(a, b string) bool {
	if v, ok := keMap[a]; ok {
		return v == b
	}
	return false
}

// GetZodiacWuXing returns the wuxing element associated with the given
// Chinese zodiac sign. It returns an empty string if the zodiac is not
// recognized.
func getZodiacWuXing(zodiac string) string {
	m := map[string]string{
		"鼠": "水", "牛": "土", "虎": "木", "兔": "木",
		"龙": "土", "蛇": "火", "马": "火", "羊": "土",
		"猴": "金", "鸡": "金", "狗": "土", "猪": "水",
	}
	if v, ok := m[zodiac]; ok {
		return v
	}
	return ""
}
