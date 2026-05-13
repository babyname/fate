package seeddb

import (
	"fmt"

	"github.com/babyname/fate/internal/wuge"
)

func (e *Exporter) transformNCharacters(oldList []oldNCharacter, idLookup map[int]string) []SeedCharacter {
	result := make([]SeedCharacter, 0, len(oldList))
	for _, old := range oldList {
		result = append(result, e.transformNCharacter(old, idLookup))
	}
	return result
}

func (e *Exporter) transformNCharacter(old oldNCharacter, idLookup map[int]string) SeedCharacter {
	wx := old.WuXing
	wxSource := "original"
	if wx == "岁" {
		e.recordChange(old.Char, "wu_xing", "岁", "", "dirty_data_cleared", "n_character")
		wx = ""
	}
	if wx == "" {
		inferred := inferWuXing(old.Radical)
		if inferred != "" {
			e.recordChange(old.Char, "wu_xing", "", inferred, "infer_from_radical:"+old.Radical, "n_character")
			wx = inferred
			wxSource = "inferred:radical"
		}
	}

	origNameable := fmt.Sprintf("%v", old.IsScience)
	newNameable := determineNameable(old.IsScience, old.ScienceStroke, old.CharStroke, old.NeedFix)
	if fmt.Sprintf("%v", newNameable) != origNameable {
		e.recordChange(old.Char, "nameable", origNameable, fmt.Sprintf("%v", newNameable), "determineNameable", "n_character")
	}

	sc := SeedCharacter{
		Char:          old.Char,
		IsSimplified:  old.IsSimplified,
		IsTraditional: old.IsTraditional,
		IsKangxi:      old.IsKangXi,
		IsVariant:     old.IsVariant,
		Regular:       old.IsRegular,
		Nameable:      newNameable,
		WuXing:        wx,
		Source:        "n_character",
	}

	sc.Pinyin = parseJSONString(old.PinYin)
	if len(sc.Pinyin) == 0 {
		if pinyins, ok := e.pinyinMap[old.Char]; ok && len(pinyins) > 0 {
			e.recordChange(old.Char, "pinyin", "", fmt.Sprintf("%v", pinyins), "unihan_pinyin", "Unihan")
			sc.Pinyin = pinyins
		}
	}

	if old.Radical != "" {
		sc.Radical = old.Radical
		sc.RadicalStroke = old.RadicalStroke
	}

	sc.ScienceStroke = old.ScienceStroke
	sc.KangxiStroke = old.KangXiStroke

	if old.CharStroke > 0 && old.ScienceStroke == 0 {
		e.recordChange(old.Char, "science_stroke", "0", fmt.Sprintf("%d", old.CharStroke), "fallback_from_char_stroke", "n_character")
		sc.ScienceStroke = old.CharStroke
	}

	simplifiedIDs := parseJSONInts(old.SimplifiedID)
	if len(simplifiedIDs) > 0 {
		resolved := resolveIDs(simplifiedIDs, idLookup)
		if len(resolved) > 0 && resolved[0] != fmt.Sprintf("id:%d", simplifiedIDs[0]) {
			e.recordChange(old.Char, "simplified_of_char", fmt.Sprintf("id:%d", simplifiedIDs[0]), resolved[0], "resolve_id_to_char", "n_character")
			sc.SimplifiedOfChar = resolved[0]
		} else {
			sc.SimplifiedOfChar = resolved[0]
		}
	}

	_ = wxSource

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

func (e *Exporter) transformCharacters(oldList []oldCharacter) []SeedCharacter {
	result := make([]SeedCharacter, 0, len(oldList))
	for _, old := range oldList {
		result = append(result, e.transformCharacter(old))
	}
	return result
}

func (e *Exporter) transformCharacter(old oldCharacter) SeedCharacter {
	wx := old.WuXing
	if wx == "岁" {
		e.recordChange(old.Ch, "wu_xing", "岁", "", "dirty_data_cleared", "character")
		wx = ""
	}
	if wx == "" {
		inferred := inferWuXing(old.Radical)
		if inferred != "" {
			e.recordChange(old.Ch, "wu_xing", "", inferred, "infer_from_radical:"+old.Radical, "character")
			wx = inferred
		}
	}

	origNameable := fmt.Sprintf("%v", old.NameScience)
	newNameable := determineNameable(old.NameScience, old.ScienceStroke, old.Stroke, false)
	if fmt.Sprintf("%v", newNameable) != origNameable {
		e.recordChange(old.Ch, "nameable", origNameable, fmt.Sprintf("%v", newNameable), "determineNameable", "character")
	}

	sc := SeedCharacter{
		Char:      old.Ch,
		IsKangxi:  old.IsKangXi,
		Regular:   old.Regular,
		Nameable:  newNameable,
		WuXing:    wx,
		Source:    "character",
	}

	sc.Pinyin = parseJSONString(old.PinYin)
	if len(sc.Pinyin) == 0 {
		if pinyins, ok := e.pinyinMap[old.Ch]; ok && len(pinyins) > 0 {
			e.recordChange(old.Ch, "pinyin", "", fmt.Sprintf("%v", pinyins), "unihan_pinyin", "Unihan")
			sc.Pinyin = pinyins
		}
	}

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
	"口": "木", "言": "木", "讠": "木", "耳": "火", "心": "火", "忄": "火",
	"手": "火", "扌": "火", "足": "火", "走": "火", "辶": "火", "廴": "火",
	"目": "火", "见": "火", "贝": "金", "页": "金", "骨": "金", "韦": "金",
	"风": "水", "飞": "水", "马": "火", "牛": "土", "羊": "土", "犬": "土",
	"犭": "土", "豕": "水", "虫": "火", "衣": "木", "衤": "木", "巾": "木",
	"糸": "火", "纟": "火", "丝": "火", "革": "木", "角": "木", "皮": "木",
	"毛": "水", "老": "火", "长": "火", "高": "木", "门": "金", "户": "水",
	"广": "木", "宀": "木", "穴": "土", "瓦": "土", "皿": "水",
	"食": "水", "饣": "水", "鬼": "水", "示": "金", "礻": "金", "夕": "金",
	"大": "火", "女": "水", "子": "水", "又": "土", "寸": "金", "工": "木",
	"士": "金", "八": "金", "人": "金", "亻": "金", "入": "金", "匕": "金",
	"匚": "金", "弓": "木", "彳": "金", "攵": "火", "斗": "火", "方": "水",
	"父": "水", "谷": "木", "虍": "火", "臼": "金", "里": "火",
	"鹿": "火", "麦": "木", "麻": "木", "黍": "木", "鸟": "火", "齐": "金",
	"肉": "木", "色": "金", "舌": "火", "身": "金", "鼠": "水", "音": "土",
	"隹": "土", "爪": "金", "舟": "木",
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
