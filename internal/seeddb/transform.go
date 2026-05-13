package seeddb

import (
	"fmt"

	"github.com/babyname/fate/internal/wuge"
)

func transformNCharacters(oldList []oldNCharacter) []SeedCharacter {
	result := make([]SeedCharacter, 0, len(oldList))
	for _, old := range oldList {
		result = append(result, transformNCharacter(old))
	}
	return result
}

func transformNCharacter(old oldNCharacter) SeedCharacter {
	wx := old.WuXing
	if wx == "岁" {
		wx = ""
	}
	if wx == "" {
		wx = inferWuXing(old.Radical)
	}

	sc := SeedCharacter{
		Char:          old.Char,
		IsSimplified:  old.IsSimplified,
		IsTraditional: old.IsTraditional,
		IsKangxi:      old.IsKangXi,
		IsVariant:     old.IsVariant,
		Regular:       old.IsRegular,
		Nameable:      determineNameable(old.IsScience, old.ScienceStroke, old.CharStroke, old.NeedFix),
		WuXing:        wx,
		Source:        "n_character",
	}

	sc.Pinyin = parseJSONString(old.PinYin)

	if old.Radical != "" {
		sc.Radical = old.Radical
		sc.RadicalStroke = old.RadicalStroke
	}

	sc.ScienceStroke = old.ScienceStroke
	sc.KangxiStroke = old.KangXiStroke

	if old.CharStroke > 0 && old.ScienceStroke == 0 {
		sc.ScienceStroke = old.CharStroke
	}

	simplifiedIDs := parseJSONInts(old.SimplifiedID)
	if len(simplifiedIDs) > 0 {
		sc.SimplifiedOfChar = fmt.Sprintf("id:%d", simplifiedIDs[0])
	}

	if old.Explanation != "" {
		sc.Meaning = old.Explanation
	}

	commentParts := parseJSONString(old.Comment)
	if len(commentParts) > 0 {
		if len(commentParts) == 1 {
			if len(commentParts[0]) > 200 {
				sc.Comment = commentParts[0][:200] + "..."
			} else {
				sc.Comment = commentParts[0]
			}
		} else {
			sc.Comment = fmt.Sprintf("%v", commentParts)
		}
	}

	if old.Lucky != "" {
		if sc.Comment != "" {
			sc.Comment = fmt.Sprintf("lucky=%s; %s", old.Lucky, sc.Comment)
		} else {
			sc.Comment = fmt.Sprintf("lucky=%s", old.Lucky)
		}
	}

	if old.NeedFix {
		if sc.Comment != "" {
			sc.Comment += "; [need_fix]"
		} else {
			sc.Comment = "[need_fix]"
		}
	}

	return sc
}

func transformCharacters(oldList []oldCharacter) []SeedCharacter {
	result := make([]SeedCharacter, 0, len(oldList))
	for _, old := range oldList {
		result = append(result, transformCharacter(old))
	}
	return result
}

func transformCharacter(old oldCharacter) SeedCharacter {
	wx := old.WuXing
	if wx == "岁" {
		wx = ""
	}
	if wx == "" {
		wx = inferWuXing(old.Radical)
	}

	sc := SeedCharacter{
		Char:      old.Ch,
		IsKangxi:  old.IsKangXi,
		Regular:   old.Regular,
		Nameable:  determineNameable(old.NameScience, old.ScienceStroke, old.Stroke, false),
		WuXing:    wx,
		Source:    "character",
	}

	sc.Pinyin = parseJSONString(old.PinYin)

	if old.Radical != "" {
		sc.Radical = old.Radical
		sc.RadicalStroke = old.RadicalStroke
	}

	sc.SimplifiedStroke = old.SimpleTotalStroke
	sc.TraditionalStroke = old.TraditionalTotalStroke
	sc.KangxiStroke = old.KangXiStroke
	sc.ScienceStroke = old.ScienceStroke

	tradChars := parseJSONString(old.TraditionalCharacter)
	if len(tradChars) > 0 {
		sc.IsSimplified = true
		sc.SimplifiedOfChar = tradChars[0]
	}

	if old.TraditionalTotalStroke > 0 && old.SimpleTotalStroke != old.TraditionalTotalStroke && len(tradChars) == 0 {
		sc.IsTraditional = true
	}

	if old.KangXi != "" && old.KangXi != old.Ch {
		sc.IsVariant = true
		sc.VariantOfChar = old.KangXi
	}

	variantChars := parseJSONString(old.VariantCharacter)
	if len(variantChars) > 0 {
		sc.IsVariant = true
		if sc.VariantOfChar == "" {
			sc.VariantOfChar = variantChars[0]
		}
	}

	commentParts := parseJSONString(old.Comment)
	if old.Lucky != "" {
		sc.Comment = fmt.Sprintf("lucky=%s", old.Lucky)
	}
	if len(commentParts) > 0 {
		if sc.Comment != "" {
			sc.Comment += "; "
		}
		if len(commentParts) == 1 {
			if len(commentParts[0]) > 200 {
				sc.Comment += commentParts[0][:200] + "..."
			} else {
				sc.Comment += commentParts[0]
			}
		}
	}

	return sc
}

func transformWuGeLucky(oldList []oldWuGeLucky) []SeedWuGeLucky {
	result := make([]SeedWuGeLucky, 0, len(oldList))
	for _, old := range oldList {
		result = append(result, transformWuGeLuckyOne(old))
	}
	return result
}

func transformWuGeLuckyOne(old oldWuGeLucky) SeedWuGeLucky {
	ge := wuge.CalcWuGe(old.LastStroke1, old.LastStroke2, old.FirstStroke1, old.FirstStroke2)

	tianDaYan := wuge.Find(ge.TianGe())
	renDaYan := wuge.Find(ge.RenGe())
	diDaYan := wuge.Find(ge.DiGe())
	waiDaYan := wuge.Find(ge.WaiGe())
	zongDaYan := wuge.Find(ge.ZongGe())

	return SeedWuGeLucky{
		LastStroke1:  old.LastStroke1,
		LastStroke2:  old.LastStroke2,
		FirstStroke1: old.FirstStroke1,
		FirstStroke2: old.FirstStroke2,
		TianGe:       ge.TianGe(),
		TianDaYan:    formatDaYan(tianDaYan),
		RenGe:        ge.RenGe(),
		RenDaYan:     formatDaYan(renDaYan),
		DiGe:         ge.DiGe(),
		DiDaYan:      formatDaYan(diDaYan),
		WaiGe:        ge.WaiGe(),
		WaiDaYan:     formatDaYan(waiDaYan),
		ZongGe:       ge.ZongGe(),
		ZongDaYan:    formatDaYan(zongDaYan),
		ZongLucky:    old.ZongLucky,
		ZongSex:      bool(zongDaYan.Sex),
		ZongMax:      zongDaYan.Max,
	}
}

func formatDaYan(dy wuge.DaYan) string {
	return fmt.Sprintf("%d|%s|%s", dy.Number, dy.Lucky, dy.SkyNine)
}

func transformWuXing(oldList []oldWuXing) []SeedWuXing {
	result := make([]SeedWuXing, 0, len(oldList))
	for _, old := range oldList {
		result = append(result, SeedWuXing{
			ID:      old.ID,
			First:   old.First,
			Second:  old.Second,
			Third:   old.Third,
			Fortune: old.Fortune,
		})
	}
	return result
}

func determineNameable(isScience bool, scienceStroke int, fallbackStroke int, needFix bool) bool {
	if isScience {
		return true
	}
	stroke := scienceStroke
	if stroke <= 0 {
		stroke = fallbackStroke
	}
	if stroke <= 0 {
		return false
	}
	if needFix {
		return false
	}
	return stroke >= 1 && stroke <= 30
}

var radicalWuXing = map[string]string{
	"木": "木", "禾": "木", "竹": "木", "艹": "木", "草": "木", "耒": "木",
	"火": "火", "灬": "火", "日": "火", "光": "火", "炎": "火",
	"土": "土", "山": "土", "石": "土", "田": "土", "邑": "土", "阝": "土",
	"金": "金", "钅": "金", "刀": "金", "刂": "金", "力": "金", "戈": "金",
	"水": "水", "氵": "水", "雨": "水", "冫": "水", "鱼": "水", "酉": "水",
}

func inferWuXing(radical string) string {
	if radical == "" {
		return ""
	}
	if wx, ok := radicalWuXing[radical]; ok {
		return wx
	}
	return ""
}
