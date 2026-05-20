package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	_ "github.com/sqlite3ent/sqlite3"
)

const dsn = "file:fate?cache=shared&_journal=WAL&_fk=1"

var validWuXing = map[string]bool{
	"金": true, "木": true, "水": true, "火": true, "土": true,
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

var pinyinPattern = regexp.MustCompile(`^[a-z][a-züāáǎàēéěèīíǐìōóǒòūúǔùǖǘǚǜ0-9]*$`)

var (
	okCount   int
	warnCount int
	errCount  int
)

func main() {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	var tableExists int
	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='character'").Scan(&tableExists)
	if tableExists == 0 {
		log.Fatal("character 表不存在于数据库中")
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("  Character 数据质量审计报告")
	fmt.Println(strings.Repeat("=", 60))

	auditOverallStats(db)
	auditWuXing(db)
	auditPinyin(db)
	auditStrokes(db)
	auditSimplifiedTraditional(db)
	auditSurnames(db)
	auditGB2312(db)

	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("  审计完成: %d OK, %d WARN, %d ERROR\n", okCount, warnCount, errCount)
	fmt.Println(strings.Repeat("=", 60))
}

func printResult(level, msg string) {
	fmt.Printf("  [%s] %s\n", level, msg)
	switch level {
	case "OK":
		okCount++
	case "WARN":
		warnCount++
	case "ERROR":
		errCount++
	}
}

func printSection(title string) {
	fmt.Println()
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("  %s\n", title)
	fmt.Println(strings.Repeat("-", 60))
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

func auditOverallStats(db *sql.DB) {
	printSection("1. 总体统计")

	var total int
	if err := db.QueryRow("SELECT COUNT(*) FROM character").Scan(&total); err != nil {
		printResult("ERROR", fmt.Sprintf("查询总记录数失败: %v", err))
		return
	}
	printResult("OK", fmt.Sprintf("总记录数: %d", total))

	type fieldCheck struct {
		name  string
		query string
	}

	checks := []fieldCheck{
		{"char", "SELECT COUNT(*) FROM character WHERE char IS NOT NULL AND char != ''"},
		{"unicode", "SELECT COUNT(*) FROM character WHERE unicode IS NOT NULL AND unicode != ''"},
		{"pinyin", "SELECT COUNT(*) FROM character WHERE pinyin IS NOT NULL AND pinyin != '' AND pinyin != '[]' AND pinyin != 'null'"},
		{"radical", "SELECT COUNT(*) FROM character WHERE radical IS NOT NULL AND radical != ''"},
		{"radical_stroke", "SELECT COUNT(*) FROM character WHERE radical_stroke IS NOT NULL AND radical_stroke > 0"},
		{"wu_xing", "SELECT COUNT(*) FROM character WHERE wu_xing IS NOT NULL AND wu_xing != ''"},
		{"simplified_stroke", "SELECT COUNT(*) FROM character WHERE simplified_stroke IS NOT NULL AND simplified_stroke > 0"},
		{"traditional_stroke", "SELECT COUNT(*) FROM character WHERE traditional_stroke IS NOT NULL AND traditional_stroke > 0"},
		{"kangxi_stroke", "SELECT COUNT(*) FROM character WHERE kangxi_stroke IS NOT NULL AND kangxi_stroke > 0"},
		{"science_stroke", "SELECT COUNT(*) FROM character WHERE science_stroke IS NOT NULL AND science_stroke > 0"},
		{"meaning", "SELECT COUNT(*) FROM character WHERE meaning IS NOT NULL AND meaning != ''"},
		{"regular", "SELECT COUNT(*) FROM character WHERE regular = 1"},
		{"nameable", "SELECT COUNT(*) FROM character WHERE nameable = 1"},
		{"is_simplified", "SELECT COUNT(*) FROM character WHERE is_simplified = 1"},
		{"is_traditional", "SELECT COUNT(*) FROM character WHERE is_traditional = 1"},
		{"is_kangxi", "SELECT COUNT(*) FROM character WHERE is_kangxi = 1"},
		{"source", "SELECT COUNT(*) FROM character WHERE source IS NOT NULL AND source != ''"},
	}

	fmt.Println()
	fmt.Printf("  %-20s %8s %8s %8s\n", "字段", "非空数", "覆盖率", "缺失数")
	fmt.Println("  " + strings.Repeat("-", 48))

	for _, fc := range checks {
		var count int
		if err := db.QueryRow(fc.query).Scan(&count); err != nil {
			printResult("ERROR", fmt.Sprintf("查询 %s 失败: %v", fc.name, err))
			continue
		}
		coverage := float64(count) / float64(total) * 100
		missing := total - count
		fmt.Printf("  %-20s %8d %7.1f%% %8d\n", fc.name, count, coverage, missing)
	}
}

func auditWuXing(db *sql.DB) {
	printSection("2. 五行(wu_xing)审计")

	var total, missing int
	db.QueryRow("SELECT COUNT(*) FROM character").Scan(&total)
	db.QueryRow("SELECT COUNT(*) FROM character WHERE wu_xing IS NULL OR wu_xing = ''").Scan(&missing)

	if missing == 0 {
		printResult("OK", fmt.Sprintf("所有字符均有五行属性 (共 %d 条)", total))
	} else {
		pct := float64(missing) / float64(total) * 100
		level := "WARN"
		if pct > 20 {
			level = "ERROR"
		}
		printResult(level, fmt.Sprintf("缺失五行的字符: %d (%.1f%%)", missing, pct))

		rows, err := db.Query(
			"SELECT char, COALESCE(radical,''), COALESCE(source,'') FROM character WHERE wu_xing IS NULL OR wu_xing = '' LIMIT 20",
		)
		if err == nil {
			fmt.Println("    缺失样本:")
			for rows.Next() {
				var ch, radical, source string
				rows.Scan(&ch, &radical, &source)
				fmt.Printf("      %s (部首=%s, 来源=%s)\n", ch, radical, source)
			}
			rows.Close()
		}
	}

	fmt.Println()
	fmt.Println("  五行值分布:")
	rows, err := db.Query(
		"SELECT wu_xing, COUNT(*) as cnt FROM character WHERE wu_xing IS NOT NULL AND wu_xing != '' GROUP BY wu_xing ORDER BY cnt DESC",
	)
	if err != nil {
		printResult("ERROR", fmt.Sprintf("查询五行分布失败: %v", err))
	} else {
		for rows.Next() {
			var wx string
			var cnt int
			rows.Scan(&wx, &cnt)
			level := "OK"
			if !validWuXing[wx] {
				level = "ERROR"
			}
			fmt.Printf("    [%s] %s: %d (%.1f%%)\n", level, wx, cnt, float64(cnt)/float64(total)*100)
		}
		rows.Close()
	}

	rows, err = db.Query(
		"SELECT DISTINCT wu_xing FROM character WHERE wu_xing IS NOT NULL AND wu_xing != '' AND wu_xing NOT IN ('金','木','水','火','土')",
	)
	if err == nil {
		var invalidValues []string
		for rows.Next() {
			var v string
			rows.Scan(&v)
			invalidValues = append(invalidValues, v)
		}
		rows.Close()
		if len(invalidValues) > 0 {
			printResult("ERROR", fmt.Sprintf("非法五行值: %s", strings.Join(invalidValues, ", ")))
			for _, v := range invalidValues {
				var cnt int
				db.QueryRow("SELECT COUNT(*) FROM character WHERE wu_xing = ?", v).Scan(&cnt)
				fmt.Printf("      '%s': %d 条记录\n", v, cnt)
				sampleRows, _ := db.Query("SELECT char FROM character WHERE wu_xing = ? LIMIT 5", v)
				var samples []string
				for sampleRows.Next() {
					var ch string
					sampleRows.Scan(&ch)
					samples = append(samples, ch)
				}
				sampleRows.Close()
				if len(samples) > 0 {
					fmt.Printf("        样本: %s\n", strings.Join(samples, " "))
				}
			}
		} else {
			printResult("OK", "未发现非法五行值")
		}
	}
}

func auditPinyin(db *sql.DB) {
	printSection("3. 拼音(pinyin)审计")

	var total, missing int
	db.QueryRow("SELECT COUNT(*) FROM character").Scan(&total)
	db.QueryRow(
		"SELECT COUNT(*) FROM character WHERE pinyin IS NULL OR pinyin = '' OR pinyin = '[]' OR pinyin = 'null'",
	).Scan(&missing)

	if missing == 0 {
		printResult("OK", fmt.Sprintf("所有字符均有拼音 (共 %d 条)", total))
	} else {
		pct := float64(missing) / float64(total) * 100
		level := "WARN"
		if pct > 20 {
			level = "ERROR"
		}
		printResult(level, fmt.Sprintf("缺失拼音的字符: %d (%.1f%%)", missing, pct))

		rows, err := db.Query(
			"SELECT char, COALESCE(wu_xing,''), COALESCE(source,'') FROM character WHERE pinyin IS NULL OR pinyin = '' OR pinyin = '[]' OR pinyin = 'null' LIMIT 20",
		)
		if err == nil {
			fmt.Println("    缺失样本:")
			for rows.Next() {
				var ch, wx, source string
				rows.Scan(&ch, &wx, &source)
				fmt.Printf("      %s (五行=%s, 来源=%s)\n", ch, wx, source)
			}
			rows.Close()
		}
	}

	fmt.Println()
	fmt.Println("  拼音格式检查:")
	rows, err := db.Query(
		"SELECT char, pinyin FROM character WHERE pinyin IS NOT NULL AND pinyin != '' AND pinyin != '[]' AND pinyin != 'null'",
	)
	if err != nil {
		printResult("ERROR", fmt.Sprintf("查询拼音失败: %v", err))
		return
	}
	defer rows.Close()

	var anomalyCount int
	var anomalies []string
	for rows.Next() {
		var ch, pinyinJSON string
		rows.Scan(&ch, &pinyinJSON)

		var pinyins []string
		if err := json.Unmarshal([]byte(pinyinJSON), &pinyins); err != nil {
			anomalyCount++
			if len(anomalies) < 10 {
				anomalies = append(anomalies, fmt.Sprintf("%s: JSON解析失败 (%s)", ch, truncate(pinyinJSON, 50)))
			}
			continue
		}

		for _, p := range pinyins {
			if p == "" {
				anomalyCount++
				if len(anomalies) < 10 {
					anomalies = append(anomalies, fmt.Sprintf("%s: 空拼音字符串", ch))
				}
				continue
			}
			if !pinyinPattern.MatchString(p) {
				anomalyCount++
				if len(anomalies) < 10 {
					anomalies = append(anomalies, fmt.Sprintf("%s: 格式异常 '%s'", ch, p))
				}
			}
		}
	}

	if anomalyCount == 0 {
		printResult("OK", "所有拼音格式正常")
	} else {
		printResult("WARN", fmt.Sprintf("拼音格式异常: %d 条", anomalyCount))
		for _, a := range anomalies {
			fmt.Printf("      %s\n", a)
		}
		if anomalyCount > len(anomalies) {
			fmt.Printf("      ... 还有 %d 条未显示\n", anomalyCount-len(anomalies))
		}
	}
}

func auditStrokes(db *sql.DB) {
	printSection("4. 笔画审计")

	type strokeField struct {
		name string
		col  string
	}

	fields := []strokeField{
		{"simplified_stroke", "simplified_stroke"},
		{"traditional_stroke", "traditional_stroke"},
		{"kangxi_stroke", "kangxi_stroke"},
		{"science_stroke", "science_stroke"},
	}

	fmt.Println("  各笔画字段缺失/异常统计:")
	fmt.Printf("    %-20s %8s %8s %8s\n", "字段", "缺失数", "值为0", "值为负")
	fmt.Println("    " + strings.Repeat("-", 48))

	for _, sf := range fields {
		var nullCount, zeroCount, negCount int
		db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM character WHERE %s IS NULL", sf.col)).Scan(&nullCount)
		db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM character WHERE %s = 0", sf.col)).Scan(&zeroCount)
		db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM character WHERE %s < 0", sf.col)).Scan(&negCount)

		level := "OK"
		if nullCount > 0 || negCount > 0 {
			level = "WARN"
		}
		fmt.Printf("    [%s] %-17s %8d %8d %8d\n", level, sf.name, nullCount, zeroCount, negCount)
	}

	fmt.Println()
	fmt.Println("  笔画值为0或负数的异常样本:")
	for _, sf := range fields {
		rows, err := db.Query(
			fmt.Sprintf("SELECT char, %s FROM character WHERE %s <= 0 AND %s IS NOT NULL LIMIT 10", sf.col, sf.col, sf.col),
		)
		if err == nil {
			var samples []string
			for rows.Next() {
				var ch string
				var val int
				rows.Scan(&ch, &val)
				samples = append(samples, fmt.Sprintf("%s(%d)", ch, val))
			}
			rows.Close()
			if len(samples) > 0 {
				printResult("WARN", fmt.Sprintf("%s: %s", sf.name, strings.Join(samples, " ")))
			} else {
				printResult("OK", fmt.Sprintf("%s: 无异常", sf.name))
			}
		}
	}

	fmt.Println()
	fmt.Println("  kangxi_stroke 与 science_stroke 差异过大 (差值>3):")
	rows, err := db.Query(`
		SELECT char, kangxi_stroke, science_stroke,
		       ABS(kangxi_stroke - science_stroke) as diff
		FROM character
		WHERE kangxi_stroke IS NOT NULL AND science_stroke IS NOT NULL
		    AND kangxi_stroke > 0 AND science_stroke > 0
		    AND ABS(kangxi_stroke - science_stroke) > 3
		ORDER BY diff DESC
		LIMIT 30
	`)
	if err != nil {
		printResult("ERROR", fmt.Sprintf("查询笔画差异失败: %v", err))
	} else {
		var count int
		var diffSamples []string
		for rows.Next() {
			var ch string
			var ks, ss, diff int
			rows.Scan(&ch, &ks, &ss, &diff)
			count++
			if count <= 20 {
				diffSamples = append(diffSamples, fmt.Sprintf("%s(康=%d,姓=%d,差=%d)", ch, ks, ss, diff))
			}
		}
		rows.Close()

		if count == 0 {
			printResult("OK", "未发现kangxi_stroke与science_stroke差异过大的字符")
		} else {
			printResult("WARN", fmt.Sprintf("差异过大的字符: %d 条", count))
			for i := 0; i < len(diffSamples); i++ {
				fmt.Printf("      %s", diffSamples[i])
				if (i+1)%3 == 0 || i == len(diffSamples)-1 {
					fmt.Println()
				} else {
					fmt.Print("  ")
				}
			}
			if count > 20 {
				fmt.Printf("      ... 还有 %d 条未显示\n", count-20)
			}
		}
	}
}

func auditSimplifiedTraditional(db *sql.DB) {
	printSection("5. 简繁体标记审计")

	var total int
	db.QueryRow("SELECT COUNT(*) FROM character").Scan(&total)

	var simpCount, tradCount, kangxiCount, neitherCount, bothCount int
	db.QueryRow("SELECT COUNT(*) FROM character WHERE is_simplified = 1").Scan(&simpCount)
	db.QueryRow("SELECT COUNT(*) FROM character WHERE is_traditional = 1").Scan(&tradCount)
	db.QueryRow("SELECT COUNT(*) FROM character WHERE is_kangxi = 1").Scan(&kangxiCount)
	db.QueryRow("SELECT COUNT(*) FROM character WHERE is_simplified = 0 AND is_traditional = 0").Scan(&neitherCount)
	db.QueryRow("SELECT COUNT(*) FROM character WHERE is_simplified = 1 AND is_traditional = 1").Scan(&bothCount)

	fmt.Printf("  is_simplified=true: %d (%.1f%%)\n", simpCount, float64(simpCount)/float64(total)*100)
	fmt.Printf("  is_traditional=true: %d (%.1f%%)\n", tradCount, float64(tradCount)/float64(total)*100)
	fmt.Printf("  is_kangxi=true: %d (%.1f%%)\n", kangxiCount, float64(kangxiCount)/float64(total)*100)

	if neitherCount == 0 {
		printResult("OK", "所有字符均有简体或繁体标记")
	} else {
		pct := float64(neitherCount) / float64(total) * 100
		level := "WARN"
		if pct > 50 {
			level = "ERROR"
		}
		printResult(level, fmt.Sprintf("既非简体又非繁体的字符: %d (%.1f%%)", neitherCount, pct))

		rows, err := db.Query(
			"SELECT char, is_kangxi, is_variant, is_ancient, COALESCE(source,'') FROM character WHERE is_simplified = 0 AND is_traditional = 0 LIMIT 20",
		)
		if err == nil {
			fmt.Println("    样本:")
			for rows.Next() {
				var ch, source string
				var isKangxi, isVariant, isAncient int
				rows.Scan(&ch, &isKangxi, &isVariant, &isAncient, &source)
				flags := []string{}
				if isKangxi == 1 {
					flags = append(flags, "康熙")
				}
				if isVariant == 1 {
					flags = append(flags, "异体")
				}
				if isAncient == 1 {
					flags = append(flags, "古字")
				}
				flagStr := strings.Join(flags, ",")
				if flagStr == "" {
					flagStr = "无标记"
				}
				fmt.Printf("      %s (标记=%s, 来源=%s)\n", ch, flagStr, source)
			}
			rows.Close()
		}
	}

	fmt.Println()
	if bothCount == 0 {
		printResult("OK", "无同时标记为简体和繁体的字符")
	} else {
		printResult("WARN", fmt.Sprintf("同时标记为简体和繁体的字符: %d", bothCount))
		rows, err := db.Query("SELECT char FROM character WHERE is_simplified = 1 AND is_traditional = 1 LIMIT 10")
		if err == nil {
			var samples []string
			for rows.Next() {
				var ch string
				rows.Scan(&ch)
				samples = append(samples, ch)
			}
			rows.Close()
			fmt.Printf("      样本: %s\n", strings.Join(samples, " "))
		}
	}
}

func auditSurnames(db *sql.DB) {
	printSection("6. 常用姓氏检查")

	fmt.Println("  百家姓前100个姓氏检查:")
	var missingSurnames []string
	for _, r := range baijiaxing100 {
		ch := string(r)
		var cnt int
		db.QueryRow("SELECT COUNT(*) FROM character WHERE char = ?", ch).Scan(&cnt)
		if cnt == 0 {
			missingSurnames = append(missingSurnames, ch)
		}
	}

	if len(missingSurnames) == 0 {
		printResult("OK", "百家姓前100个姓氏全部在数据库中")
	} else {
		printResult("ERROR", fmt.Sprintf("百家姓前100个姓氏中缺失 %d 个: %s",
			len(missingSurnames), strings.Join(missingSurnames, " ")))
	}

	fmt.Println()
	fmt.Println("  复姓用字检查:")
	var missingCompound []string
	for _, ch := range compoundSurnameChars {
		var cnt int
		db.QueryRow("SELECT COUNT(*) FROM character WHERE char = ?", ch).Scan(&cnt)
		if cnt == 0 {
			missingCompound = append(missingCompound, ch)
		}
	}

	if len(missingCompound) == 0 {
		printResult("OK", "所有复姓用字均在数据库中")
	} else {
		printResult("ERROR", fmt.Sprintf("复姓用字缺失 %d 个: %s",
			len(missingCompound), strings.Join(missingCompound, " ")))
	}

	fmt.Println()
	fmt.Println("  百家姓姓氏数据完整性:")
	var completeCount int
	for _, r := range baijiaxing100 {
		ch := string(r)
		var wx string
		var ks, ss, nameable int
		err := db.QueryRow(
			"SELECT COALESCE(wu_xing,''), COALESCE(kangxi_stroke,0), COALESCE(science_stroke,0), COALESCE(nameable,0) FROM character WHERE char = ?",
			ch,
		).Scan(&wx, &ks, &ss, &nameable)
		if err != nil {
			fmt.Printf("    [ERROR] %s: 未找到\n", ch)
			continue
		}
		issues := []string{}
		if wx == "" {
			issues = append(issues, "缺五行")
		}
		if ks == 0 {
			issues = append(issues, "缺康熙笔画")
		}
		if ss == 0 {
			issues = append(issues, "缺姓名学笔画")
		}
		if nameable == 0 {
			issues = append(issues, "不可用于起名")
		}
		if len(issues) > 0 {
			printResult("WARN", fmt.Sprintf("%s: 五行=%s 康熙=%d 姓名=%d 可起名=%v (%s)",
				ch, wx, ks, ss, nameable == 1, strings.Join(issues, ", ")))
		} else {
			completeCount++
		}
	}
	printResult("OK", fmt.Sprintf("百家姓中数据完整的姓氏: %d / %d", completeCount, len(baijiaxing100)))
}

func auditGB2312(db *sql.DB) {
	printSection("7. GB2312一级字库覆盖率检查")

	gb2312Chars := generateGB2312Level1()
	fmt.Printf("  GB2312一级字库总字数: %d\n", len(gb2312Chars))

	existingChars := make(map[string]bool)
	rows, err := db.Query("SELECT char FROM character")
	if err != nil {
		printResult("ERROR", fmt.Sprintf("查询字符失败: %v", err))
		return
	}
	for rows.Next() {
		var ch string
		rows.Scan(&ch)
		existingChars[ch] = true
	}
	rows.Close()

	var missingChars []string
	for _, ch := range gb2312Chars {
		if !existingChars[ch] {
			missingChars = append(missingChars, ch)
		}
	}

	covered := len(gb2312Chars) - len(missingChars)
	coverage := float64(covered) / float64(len(gb2312Chars)) * 100
	fmt.Printf("  数据库中已有: %d / %d\n", covered, len(gb2312Chars))
	fmt.Printf("  覆盖率: %.1f%%\n", coverage)

	if len(missingChars) == 0 {
		printResult("OK", "GB2312一级字库完全覆盖")
	} else {
		level := "WARN"
		if coverage < 80 {
			level = "ERROR"
		}
		printResult(level, fmt.Sprintf("缺失 %d 个常用字 (%.1f%%)", len(missingChars), 100-coverage))

		fmt.Println()
		fmt.Println("  缺失字列表 (每行20个):")
		for i := 0; i < len(missingChars); i++ {
			fmt.Printf("%s", missingChars[i])
			if (i+1)%20 == 0 || i == len(missingChars)-1 {
				fmt.Println()
			} else {
				fmt.Print(" ")
			}
		}
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
