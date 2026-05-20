package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/babyname/fate/internal/seeddb"
	"golang.org/x/text/encoding/simplifiedchinese"
	_ "github.com/sqlite3ent/sqlite3"
)

var radicalWuXing = map[string]string{
	"一": "土", "丨": "木", "丶": "火", "丿": "木", "乙": "木",
	"亅": "金", "二": "火", "亠": "土",
	"人": "金", "亻": "金", "儿": "金", "入": "金", "八": "金",
	"刀": "金", "刂": "金", "力": "金", "勹": "火",
	"匕": "金", "匚": "金", "匸": "金", "十": "金",
	"卜": "水", "卩": "金", "厂": "土", "厶": "金", "又": "土",
	"口": "木", "囗": "土",
	"土": "土", "士": "金", "夂": "火", "夊": "木", "夕": "金",
	"大": "火", "女": "水", "子": "水", "宀": "木", "寸": "金",
	"小": "金", "尢": "土", "尣": "土", "尸": "土", "屮": "木",
	"山": "土", "巛": "水", "工": "木", "己": "土", "巾": "木",
	"干": "木", "幺": "火", "广": "木", "廴": "火", "廾": "木",
	"弋": "金", "弓": "木", "彐": "火", "彡": "金", "彳": "金",
	"心": "火", "忄": "火", "戈": "金", "户": "水", "手": "火",
	"扌": "火", "支": "火", "攴": "火", "攵": "火", "文": "水",
	"斗": "火", "斤": "木", "方": "水", "无": "水",
	"日": "火", "曰": "火", "月": "木", "木": "木", "欠": "木",
	"止": "金", "歹": "火", "殳": "火", "毋": "土", "比": "水",
	"毛": "水", "氏": "火", "气": "木",
	"水": "水", "氵": "水", "氺": "水",
	"火": "火", "灬": "火",
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
	"見": "火", "角": "木", "言": "木", "讠": "木",
	"谷": "木", "豆": "木", "豕": "水", "豸": "火", "貝": "金",
	"赤": "火", "走": "火", "足": "火", "身": "金",
	"車": "金", "辛": "金",
	"辰": "土", "辵": "火", "辶": "火", "邑": "土", "阝": "土",
	"酉": "水", "釆": "木", "里": "火",
	"金": "金", "钅": "金", "長": "火",
	"門": "金", "阜": "土", "隶": "水", "隹": "土",
	"雨": "水", "靑": "木", "非": "火",
	"面": "金", "革": "木", "韋": "土", "韭": "木",
	"音": "土", "頁": "金", "風": "水", "飛": "水",
	"食": "水", "首": "火", "香": "木",
	"馬": "火", "骨": "金", "高": "木", "髟": "木",
	"鬥": "木", "鬲": "火", "鬼": "水", "魚": "水",
	"鳥": "火", "鹵": "火", "鹿": "火", "麥": "木",
	"麻": "木", "黃": "土", "黍": "木", "黑": "水",
	"黹": "火", "黽": "水", "鼎": "火", "鼓": "木", "鼠": "水",
	"鼻": "金", "齊": "金", "龍": "土", "龜": "水", "龠": "木",
	"齐": "金", "龙": "土", "龟": "水",
	"马": "火", "鱼": "水", "鸟": "火", "麦": "木", "黄": "土",
	"见": "火", "页": "金", "风": "火", "飞": "水", "齿": "土",
	"釒": "金", "糹": "火", "飠": "水", "齒": "土",
	"冫": "水", "虍": "火", "罒": "水", "几": "木",
	"冖": "木", "冂": "木", "镸": "火", "覀": "土",
	"艸": "木", "旡": "木", "彑": "木", "耂": "火",
	"凵": "木", "巜": "水", "丬": "木", "爫": "金",
	"龸": "金", "舁": "金",
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

var simplifiedRadicals = map[string]bool{
	"钅": true, "讠": true, "饣": true, "门": true, "韦": true,
	"马": true, "鸟": true, "鱼": true, "龙": true, "齿": true,
	"车": true, "贝": true, "见": true, "页": true, "风": true,
	"飞": true, "纟": true, "专": true, "业": true, "东": true,
	"长": true,
}

var validWuXing = map[string]bool{
	"金": true, "木": true, "水": true, "火": true, "土": true,
}

var rareRadicals = map[string]bool{
	"鬼": true, "鼠": true, "鬥": true, "鹵": true, "黹": true,
	"黽": true, "龠": true, "鼎": true, "鬲": true, "殳": true,
	"歹": true, "疒": true,
}

var baijiaxing100 = []rune(
	"赵钱孙李周吴郑王冯陈褚卫蒋沈韩杨" +
		"朱秦尤许何吕施张孔曹严华金魏陶姜" +
		"戚谢邹喻柏水窦章云苏潘葛奚范彭郎" +
		"鲁韦昌马苗凤花方俞任袁柳酆鲍史唐" +
		"费廉岑薛雷贺倪汤滕殷罗毕郝邬安常" +
		"乐于时傅皮卞齐康伍余元卜顾孟平黄" +
		"和穆萧尹",
)

var compoundSurnameChars = []string{
	"西", "门", "东", "方", "南", "宫", "北", "堂",
	"上", "官", "司", "马", "诸", "葛", "欧", "阳",
}

type FixStats struct {
	Total                int
	WuXingPreserved      int
	WuXingInferred       int
	WuXingInvalidCleared int
	WuXingStillMissing   int
	SimpFixed            int
	TradFixed            int
	KangxiFixed          int
	NameableToTrue       int
	NameableToFalse      int
	RegularToTrue        int
	PinyinFixed          int
}

func parseJSONString(s string) []string {
	if s == "" || s == "null" {
		return nil
	}
	var result []string
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		return []string{s}
	}
	return result
}

func inferWuXingFromRadical(radical string) string {
	if radical == "" {
		return ""
	}
	if wx, ok := radicalWuXing[radical]; ok {
		return wx
	}
	return ""
}

func appendComment(sc *seeddb.SeedCharacter, s string) {
	if sc.Comment != "" {
		sc.Comment += "; " + s
	} else {
		sc.Comment = s
	}
}

func generateGB2312Level1() []string {
	var chars []string
	decoder := simplifiedchinese.GBK.NewDecoder()
	for hi := 0xB0; hi <= 0xD7; hi++ {
		loEnd := 0xFE
		if hi == 0xD7 {
			loEnd = 0xF9
		}
		for lo := 0xA1; lo <= loEnd; lo++ {
			gbBytes := []byte{byte(hi), byte(lo)}
			utf8Bytes, err := decoder.Bytes(gbBytes)
			if err == nil && len(utf8Bytes) > 0 {
				chars = append(chars, string(utf8Bytes))
			}
		}
	}
	return chars
}

func fixWuXing(sc *seeddb.SeedCharacter, stats *FixStats) {
	if validWuXing[sc.WuXing] {
		stats.WuXingPreserved++
		return
	}
	if sc.WuXing != "" {
		stats.WuXingInvalidCleared++
		appendComment(sc, fmt.Sprintf("wx_fix:clear_invalid(%s)", sc.WuXing))
	}
	sc.WuXing = ""

	if sc.Radical != "" {
		inferred := inferWuXingFromRadical(sc.Radical)
		if inferred != "" {
			sc.WuXing = inferred
			stats.WuXingInferred++
			appendComment(sc, fmt.Sprintf("wx_fix:radical(%s)→%s", sc.Radical, inferred))
			return
		}
	}

	if sc.Radical == "" && sc.RadicalStroke > 0 {
		for num, name := range kangxiRadicalNames {
			if num == sc.RadicalStroke {
				inferred := inferWuXingFromRadical(name)
				if inferred != "" {
					sc.WuXing = inferred
					stats.WuXingInferred++
					appendComment(sc, fmt.Sprintf("wx_fix:kangxi_rs(%d→%s)→%s", sc.RadicalStroke, name, inferred))
					return
				}
			}
		}
	}

	stats.WuXingStillMissing++
}

func fixSimplifiedTraditional(sc *seeddb.SeedCharacter, stats *FixStats) {
	if sc.SimplifiedStroke > 0 && !sc.IsSimplified {
		sc.IsSimplified = true
		stats.SimpFixed++
	}
	if sc.TraditionalStroke > 0 && !sc.IsTraditional {
		sc.IsTraditional = true
		stats.TradFixed++
	}
	if !sc.IsSimplified && !sc.IsTraditional && sc.ScienceStroke > 0 {
		hasSimpRadical := false
		for rad := range simplifiedRadicals {
			if strings.Contains(sc.Char, rad) {
				hasSimpRadical = true
				break
			}
		}
		sc.IsSimplified = true
		stats.SimpFixed++
		if hasSimpRadical {
			appendComment(sc, "st_fix:simp_radical_detected")
		} else {
			appendComment(sc, "st_fix:default_simplified")
		}
	}
}

func fixKangxiStroke(sc *seeddb.SeedCharacter, stats *FixStats) {
	if sc.KangxiStroke > 0 {
		return
	}
	if sc.IsSimplified && sc.TraditionalStroke > 0 {
		sc.KangxiStroke = sc.TraditionalStroke
		stats.KangxiFixed++
		appendComment(sc, fmt.Sprintf("kx_fix:simp→trad(%d)", sc.TraditionalStroke))
	} else if sc.ScienceStroke > 0 {
		sc.KangxiStroke = sc.ScienceStroke
		stats.KangxiFixed++
		appendComment(sc, fmt.Sprintf("kx_fix:fallback_science(%d)", sc.ScienceStroke))
	}
}

func fixNameable(sc *seeddb.SeedCharacter, stats *FixStats) {
	hasPinyin := len(sc.Pinyin) > 0
	strokeOK := sc.ScienceStroke >= 1 && sc.ScienceStroke <= 30
	if !strokeOK && sc.KangxiStroke >= 1 && sc.KangxiStroke <= 30 {
		strokeOK = true
	}
	isRareRadical := rareRadicals[sc.Radical]

	shouldNameable := hasPinyin && strokeOK && !isRareRadical

	if shouldNameable != sc.Nameable {
		if shouldNameable {
			stats.NameableToTrue++
		} else {
			stats.NameableToFalse++
		}
		sc.Nameable = shouldNameable
	}
}

func fixRegular(sc *seeddb.SeedCharacter, gb2312Set, baijiaxingSet map[string]bool, stats *FixStats) {
	if gb2312Set[sc.Char] || baijiaxingSet[sc.Char] {
		if !sc.Regular {
			sc.Regular = true
			stats.RegularToTrue++
		}
	}
}

func fixPinyin(sc *seeddb.SeedCharacter, stats *FixStats) {
	fixed := false
	for i, p := range sc.Pinyin {
		orig := p
		p = strings.ReplaceAll(p, "\u3000", " ")
		p = strings.TrimSpace(p)
		for strings.Contains(p, "  ") {
			p = strings.ReplaceAll(p, "  ", " ")
		}
		parts := strings.Split(p, ",")
		if len(parts) > 1 {
			for j, part := range parts {
				parts[j] = strings.TrimSpace(part)
			}
			p = strings.Join(parts, ",")
		}
		if p != orig {
			sc.Pinyin[i] = p
			fixed = true
		}
	}
	if fixed {
		stats.PinyinFixed++
	}
}

func printStats(stats *FixStats, seeds []seeddb.SeedCharacter) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("  数据修复统计报告")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("\n  总字符数: %d\n", stats.Total)

	fmt.Println("\n  --- 1. 五行修复 ---")
	fmt.Printf("  保留原有合法五行: %d\n", stats.WuXingPreserved)
	fmt.Printf("  从部首推断五行:   %d\n", stats.WuXingInferred)
	fmt.Printf("  清除非法五行值:   %d\n", stats.WuXingInvalidCleared)
	fmt.Printf("  仍缺五行:         %d\n", stats.WuXingStillMissing)

	fmt.Println("\n  --- 2. 简繁体标记修复 ---")
	fmt.Printf("  修复为简体: %d\n", stats.SimpFixed)
	fmt.Printf("  修复为繁体: %d\n", stats.TradFixed)

	fmt.Println("\n  --- 3. 康熙笔画修复 ---")
	fmt.Printf("  修复康熙笔画: %d\n", stats.KangxiFixed)

	fmt.Println("\n  --- 4. 可起名修复 ---")
	fmt.Printf("  设为可起名:   %d\n", stats.NameableToTrue)
	fmt.Printf("  设为不可起名: %d\n", stats.NameableToFalse)

	fmt.Println("\n  --- 5. 常用字修复 ---")
	fmt.Printf("  设为常用字: %d\n", stats.RegularToTrue)

	fmt.Println("\n  --- 6. 拼音格式修复 ---")
	fmt.Printf("  修复拼音格式: %d\n", stats.PinyinFixed)

	var withWuXing, withPinyin, nameable, regular, simplified, traditional int
	wuXingDist := make(map[string]int)
	for _, sc := range seeds {
		if sc.WuXing != "" {
			withWuXing++
			wuXingDist[sc.WuXing]++
		}
		if len(sc.Pinyin) > 0 {
			withPinyin++
		}
		if sc.Nameable {
			nameable++
		}
		if sc.Regular {
			regular++
		}
		if sc.IsSimplified {
			simplified++
		}
		if sc.IsTraditional {
			traditional++
		}
	}

	fmt.Println("\n  --- 修复后数据概览 ---")
	total := stats.Total
	fmt.Printf("  有五行: %d (%.1f%%)\n", withWuXing, float64(withWuXing)/float64(total)*100)
	fmt.Printf("  有拼音: %d (%.1f%%)\n", withPinyin, float64(withPinyin)/float64(total)*100)
	fmt.Printf("  可起名: %d (%.1f%%)\n", nameable, float64(nameable)/float64(total)*100)
	fmt.Printf("  常用字: %d (%.1f%%)\n", regular, float64(regular)/float64(total)*100)
	fmt.Printf("  简体字: %d (%.1f%%)\n", simplified, float64(simplified)/float64(total)*100)
	fmt.Printf("  繁体字: %d (%.1f%%)\n", traditional, float64(traditional)/float64(total)*100)

	fmt.Println("\n  五行分布:")
	for _, wx := range []string{"金", "木", "水", "火", "土"} {
		fmt.Printf("    %s: %d\n", wx, wuXingDist[wx])
	}
	var otherWx int
	for wx, cnt := range wuXingDist {
		if !validWuXing[wx] {
			otherWx += cnt
		}
	}
	if otherWx > 0 {
		fmt.Printf("    其他(非法): %d\n", otherWx)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
}

func main() {
	dbPath := "model/db/data.db"
	db, err := sql.Open("sqlite3", dbPath+"?mode=ro")
	if err != nil {
		log.Fatalf("无法打开数据库 %s: %v", dbPath, err)
	}
	defer db.Close()

	log.Printf("正在从 v4 数据库读取数据: %s", dbPath)

	gb2312Chars := generateGB2312Level1()
	gb2312Set := make(map[string]bool, len(gb2312Chars))
	for _, ch := range gb2312Chars {
		gb2312Set[ch] = true
	}
	log.Printf("GB2312一级字库: %d 字", len(gb2312Chars))

	baijiaxingSet := make(map[string]bool)
	for _, r := range baijiaxing100 {
		baijiaxingSet[string(r)] = true
	}
	for _, ch := range compoundSurnameChars {
		baijiaxingSet[ch] = true
	}

	rows, err := db.Query(`SELECT id, pin_yin, ch, science_stroke, radical, radical_stroke, stroke, is_kangxi, kangxi, kangxi_stroke, simple_radical, simple_radical_stroke, simple_total_stroke, traditional_radical, traditional_radical_stroke, traditional_total_stroke, is_name_science, wu_xing, lucky, is_regular, traditional_character, variant_character, comment FROM characters`)
	if err != nil {
		log.Fatalf("查询 characters 表失败: %v", err)
	}
	defer rows.Close()

	var seeds []seeddb.SeedCharacter
	var stats FixStats

	for rows.Next() {
		var c struct {
			ID                       string
			PinYin                   string
			Ch                       string
			ScienceStroke            int
			Radical                  string
			RadicalStroke            int
			Stroke                   int
			IsKangxi                 bool
			Kangxi                   string
			KangxiStroke             int
			SimpleRadical            string
			SimpleRadicalStroke      int
			SimpleTotalStroke        int
			TraditionalRadical       string
			TraditionalRadicalStroke int
			TraditionalTotalStroke   int
			IsNameScience            bool
			WuXing                   string
			Lucky                    string
			IsRegular                bool
			TraditionalCharacter     string
			VariantCharacter         string
			Comment                  string
		}
		err := rows.Scan(
			&c.ID, &c.PinYin, &c.Ch, &c.ScienceStroke, &c.Radical, &c.RadicalStroke,
			&c.Stroke, &c.IsKangxi, &c.Kangxi, &c.KangxiStroke, &c.SimpleRadical,
			&c.SimpleRadicalStroke, &c.SimpleTotalStroke, &c.TraditionalRadical,
			&c.TraditionalRadicalStroke, &c.TraditionalTotalStroke, &c.IsNameScience,
			&c.WuXing, &c.Lucky, &c.IsRegular, &c.TraditionalCharacter,
			&c.VariantCharacter, &c.Comment,
		)
		if err != nil {
			log.Printf("scan error: %v", err)
			continue
		}

		stats.Total++

		sc := seeddb.SeedCharacter{
			Char:              c.Ch,
			IsKangxi:          c.IsKangxi,
			Regular:           c.IsRegular,
			Nameable:          c.IsNameScience,
			WuXing:            c.WuXing,
			Radical:           c.Radical,
			RadicalStroke:     c.RadicalStroke,
			SimplifiedStroke:  c.SimpleTotalStroke,
			TraditionalStroke: c.TraditionalTotalStroke,
			KangxiStroke:      c.KangxiStroke,
			ScienceStroke:     c.ScienceStroke,
			Source:            "v4_dbfix",
		}

		sc.Pinyin = parseJSONString(c.PinYin)

		tradChars := parseJSONString(c.TraditionalCharacter)
		if len(tradChars) > 0 {
			sc.IsSimplified = true
			sc.SimplifiedOfChar = tradChars[0]
		}

		if c.Kangxi != "" && c.Kangxi != c.Ch {
			sc.IsVariant = true
			sc.VariantOfChar = c.Kangxi
		}
		variantChars := parseJSONString(c.VariantCharacter)
		if len(variantChars) > 0 {
			sc.IsVariant = true
			if sc.VariantOfChar == "" {
				sc.VariantOfChar = variantChars[0]
			}
		}

		if c.Lucky != "" {
			sc.Comment = fmt.Sprintf("lucky=%s", c.Lucky)
		}
		commentParts := parseJSONString(c.Comment)
		if len(commentParts) > 0 {
			if sc.Comment != "" {
				sc.Comment += "; "
			}
			if len(commentParts[0]) > 200 {
				sc.Comment += commentParts[0][:200] + "..."
			} else {
				sc.Comment += commentParts[0]
			}
		}

		fixWuXing(&sc, &stats)
		fixSimplifiedTraditional(&sc, &stats)
		fixKangxiStroke(&sc, &stats)
		fixNameable(&sc, &stats)
		fixRegular(&sc, gb2312Set, baijiaxingSet, &stats)
		fixPinyin(&sc, &stats)

		seeds = append(seeds, sc)
	}

	log.Printf("从 v4 数据库读取 %d 字符", len(seeds))

	outputPaths := []string{
		"data/seed/character.json",
		"internal/seeddb/data/character.json",
	}

	for _, path := range outputPaths {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("创建目录失败 %s: %v", dir, err)
		}
		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("创建文件失败 %s: %v", path, err)
		}
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(seeds); err != nil {
			f.Close()
			log.Fatalf("写入JSON失败 %s: %v", path, err)
		}
		f.Close()
		log.Printf("已写入 %s (%d 字符)", path, len(seeds))
	}

	printStats(&stats, seeds)
}
