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

	if sc.Nameable && len(sc.Pinyin) == 0 {
		e.recordChange(old.Char, "nameable", "true", "false", "no_pinyin_not_nameable", "n_character")
		sc.Nameable = false
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
		Char:     old.Ch,
		IsKangxi: old.IsKangXi,
		Regular:  old.Regular,
		Nameable: newNameable,
		WuXing:   wx,
		Source:   "character",
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

	if sc.Nameable && len(sc.Pinyin) == 0 {
		e.recordChange(old.Ch, "nameable", "true", "false", "no_pinyin_not_nameable", "character")
		sc.Nameable = false
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
		result = append(result, SeedWuXing(old))
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
	// 繁体偏旁（与简体对应）
	"釒": "金", "糹": "火", "飠": "水", "齒": "土",
	// 缺失的康熙部首
	"冫": "水", "虍": "火", "罒": "水", "几": "木",
	"冖": "木", "冂": "木", "镸": "火", "覀": "土",
	"艸": "木", "旡": "木", "彑": "木", "耂": "火",
	"凵": "木", "巜": "水", "丬": "木", "爫": "金",
	"龸": "金", "舁": "金",
	// 补充常见缺漏
	"丷": "火", "襾": "土", "韮": "木",
}

var kangxiRadicalNames = map[int]string{
	1: "一", 2: "丨", 3: "丶", 4: "丿", 5: "乙", 6: "亅",
	7: "二", 8: "亠", 9: "人", 10: "儿", 11: "入", 12: "八",
	13: "冂", 14: "冖", 15: "冫", 16: "几", 17: "凵", 18: "刀",
	19: "力", 20: "勹", 21: "匕", 22: "匚", 23: "匸", 24: "十",
	25: "卜", 26: "卩", 27: "厂", 28: "厶", 29: "又", 30: "口",
	31: "囗", 32: "土", 33: "士", 34: "夂", 35: "夊", 36: "夕",
	37: "大", 38: "女", 39: "子", 40: "宀", 41: "寸", 42: "小",
	43: "尢", 44: "尸", 45: "屮", 46: "山", 47: "巛", 48: "工",
	49: "己", 50: "巾", 51: "干", 52: "幺", 53: "广", 54: "廴",
	55: "廾", 56: "弋", 57: "弓", 58: "彐", 59: "彡", 60: "彳",
	61: "心", 62: "戈", 63: "户", 64: "手", 65: "支", 66: "攴",
	67: "文", 68: "斗", 69: "斤", 70: "方", 71: "无", 72: "日",
	73: "曰", 74: "月", 75: "木", 76: "欠", 77: "止", 78: "歹",
	79: "殳", 80: "毋", 81: "比", 82: "毛", 83: "氏", 84: "气",
	85: "水", 86: "火", 87: "爪", 88: "父", 89: "爻", 90: "爿",
	91: "片", 92: "牙", 93: "牛", 94: "犬", 95: "玄", 96: "玉",
	97: "瓜", 98: "瓦", 99: "甘", 100: "生", 101: "用", 102: "田",
	103: "疋", 104: "疒", 105: "癶", 106: "白", 107: "皮", 108: "皿",
	109: "目", 110: "矛", 111: "矢", 112: "石", 113: "示", 114: "禸",
	115: "禾", 116: "穴", 117: "立", 118: "竹", 119: "米", 120: "糸",
	121: "缶", 122: "网", 123: "羊", 124: "羽", 125: "老", 126: "而",
	127: "耒", 128: "耳", 129: "聿", 130: "肉", 131: "臣", 132: "自",
	133: "至", 134: "臼", 135: "舌", 136: "舛", 137: "舟", 138: "艮",
	139: "色", 140: "艸", 141: "虍", 142: "虫", 143: "血", 144: "行",
	145: "衣", 146: "西", 147: "見", 148: "角", 149: "言", 150: "谷",
	151: "豆", 152: "豕", 153: "豸", 154: "貝", 155: "赤", 156: "走",
	157: "足", 158: "身", 159: "車", 160: "辛", 161: "辰", 162: "辵",
	163: "邑", 164: "酉", 165: "釆", 166: "里", 167: "金", 168: "長",
	169: "門", 170: "阜", 171: "隶", 172: "隹", 173: "雨", 174: "靑",
	175: "非", 176: "面", 177: "革", 178: "韋", 179: "韭", 180: "音",
	181: "頁", 182: "風", 183: "飛", 184: "食", 185: "首", 186: "香",
	187: "馬", 188: "骨", 189: "高", 190: "髟", 191: "鬥", 192: "鬲",
	193: "鬼", 194: "魚", 195: "鳥", 196: "鹵", 197: "鹿", 198: "麥",
	199: "麻", 200: "黃", 201: "黍", 202: "黑", 203: "黹", 204: "黽",
	205: "鼎", 206: "鼓", 207: "鼠", 208: "鼻", 209: "齊", 210: "齒",
	211: "龍", 212: "龜", 213: "龠",
}

func (e *Exporter) inferWuXingForChar(char, radical string) (string, string) {
	if wx, ok := e.wuxingMap[char]; ok && wx != "" {
		return wx, "yw11_wuxing"
	}
	wx := inferWuXing(radical)
	if wx != "" {
		return wx, "infer_from_radical:" + radical
	}
	if radical == "" {
		if rsNum, ok := e.rsUnicode[char]; ok {
			if rsName, ok2 := kangxiRadicalNames[rsNum]; ok2 {
				wx := inferWuXing(rsName)
				if wx != "" {
					return wx, "infer_from_rs:" + rsName
				}
			}
		}
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
