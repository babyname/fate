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
		inferred, source := e.inferWuXingForChar(old.Char, old.Radical)
		if inferred != "" {
			e.recordChange(old.Char, "wu_xing", "", inferred, source, "n_character")
			wx = inferred
			wxSource = "inferred:" + source
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
	} else if def, ok := e.definitions[old.Char]; ok && def != "" {
		e.recordChange(old.Char, "meaning", "", def[:min(100, len(def))], "unihan_definition", "Unihan")
		sc.Meaning = def
	}

	if sc.ScienceStroke == 0 || sc.KangxiStroke == 0 {
		if strokes, ok := e.totalStrokes[old.Char]; ok && strokes > 0 {
			if sc.ScienceStroke == 0 {
				e.recordChange(old.Char, "science_stroke", "0", fmt.Sprintf("%d", strokes), "unihan_stroke", "Unihan")
				sc.ScienceStroke = strokes
			}
			if sc.KangxiStroke == 0 {
				e.recordChange(old.Char, "kangxi_stroke", "0", fmt.Sprintf("%d", strokes), "unihan_stroke", "Unihan")
				sc.KangxiStroke = strokes
			}
		}
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
		inferred, source := e.inferWuXingForChar(old.Ch, old.Radical)
		if inferred != "" {
			e.recordChange(old.Ch, "wu_xing", "", inferred, source, "character")
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
	// 基本笔画与构件
	"一": "土", "丨": "木", "丶": "火", "丿": "木", "乙": "木",
	"亅": "金", "二": "火", "亠": "土",
	// 人部
	"人": "金", "亻": "金", "儿": "金", "入": "金", "八": "金",
	// 刀部
	"刀": "金", "刂": "金", "力": "金", "勹": "火",
	// 匕部
	"匕": "金", "匚": "金", "匸": "金", "十": "金",
	// 卜部
	"卜": "水", "卩": "金", "厂": "土", "厶": "金", "又": "土",
	// 口部
	"口": "木", "囗": "土",
	// 土部
	"土": "土", "士": "金", "夂": "火", "夊": "木", "夕": "金",
	"大": "火", "女": "水", "子": "水", "宀": "木", "寸": "金",
	"小": "金", "尢": "土", "尣": "土", "尸": "土", "屮": "木",
	"山": "土", "巛": "水", "工": "木", "己": "土", "巾": "木",
	"干": "木", "幺": "火", "广": "木", "廴": "火", "廾": "木",
	"弋": "金", "弓": "木", "彐": "火", "彡": "金", "彳": "金",
	// 心部
	"心": "火", "忄": "火", "戈": "金", "户": "水", "手": "火",
	"扌": "火", "支": "火", "攴": "火", "攵": "火", "文": "水",
	"斗": "火", "斤": "木", "方": "水", "无": "水",
	// 日部
	"日": "火", "曰": "火", "月": "木", "木": "木", "欠": "木",
	"止": "金", "歹": "火", "殳": "火", "毋": "土", "比": "水",
	"毛": "水", "氏": "火", "气": "木",
	// 水部
	"水": "水", "氵": "水", "氺": "水",
	// 火部
	"火": "火", "灬": "火",
	// 爪部
	"爪": "金", "父": "水", "爻": "火", "爿": "木", "片": "木",
	"牙": "木", "牛": "土", "牜": "土", "犬": "土", "犭": "土",
	"玄": "水", "玉": "土", "王": "土", "瓜": "木", "瓦": "土",
	"甘": "土", "生": "木", "用": "木", "田": "土", "疋": "金",
	"疒": "土", "癶": "火", "白": "金", "皮": "木", "皿": "水",
	"目": "火", "矛": "金", "矢": "金", "石": "土", "示": "金",
	"礻": "金", "禸": "木", "禾": "木", "穴": "土", "立": "金",
	"竹": "木", "米": "木", "糸": "火", "纟": "火", "缶": "土",
	"网": "土", "羊": "土", "羽": "水", "老": "火", "而": "火",
	"耒": "木", "耳": "火", "聿": "木", "肉": "木", "臣": "金",
	"自": "火", "至": "火", "臼": "金", "舌": "火", "舛": "金",
	"舟": "木", "艮": "土", "色": "金", "艹": "木", "虎": "木",
	"虫": "火", "血": "水", "行": "金", "衣": "木", "衤": "木",
	"西": "金",
	// 角部
	"見": "火", "角": "木", "言": "木", "讠": "木",
	"谷": "木", "豆": "木", "豕": "水", "豸": "火", "貝": "金",
	"赤": "火", "走": "火", "足": "火", "身": "金",
	"車": "金", "辛": "金",
	// 辰部
	"辰": "土", "辵": "火", "辶": "火", "邑": "土", "阝": "土",
	"酉": "水", "釆": "木", "里": "火",
	// 金部
	"金": "金", "钅": "金", "長": "火",
	"門": "金", "阜": "土", "隶": "水", "隹": "土",
	"雨": "水", "靑": "木", "非": "火",
	// 音部
	"面": "金", "革": "木", "韋": "土", "韭": "木",
	"音": "土", "頁": "金", "風": "水", "飛": "水",
	"食": "水", "首": "火", "香": "木",
	// 马部
	"馬": "火", "骨": "金", "高": "木", "髟": "木",
	"鬥": "木", "鬲": "火", "鬼": "水", "魚": "水",
	"鳥": "火", "鹵": "火", "鹿": "火", "麥": "木",
	"麻": "木", "黃": "土", "黍": "木", "黑": "水",
	"黹": "火", "黽": "水", "鼎": "火", "鼓": "木", "鼠": "水",
	"鼻": "金", "齊": "金", "龍": "土", "龜": "水", "龠": "木",
	// 简体异体
	"齐": "金", "龙": "土", "龟": "水",
	"马": "火", "鱼": "水", "鸟": "火", "麦": "木", "黄": "土",
	"见": "火", "页": "金", "风": "火", "飞": "水", "齿": "土",
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (e *Exporter) inferWuXingForChar(char, radical string) (string, string) {
	if wx, ok := e.wuxingMap[char]; ok && wx != "" {
		return wx, "yw11_wuxing"
	}
	wx := inferWuXing(radical)
	if wx != "" {
		return wx, "infer_from_radical:" + radical
	}
	return "", ""
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
