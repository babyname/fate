package analysis

import (
	"bytes"
	"strings"
	"testing"
	"time"

	v2 "github.com/godcong/chronos/v2"
)

func TestTextFormatter(t *testing.T) {
	fateData, _ := v2.GetFateData(&v2.FateInput{
		BirthDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local),
		Gender:    1,
		Surname:   "张",
	})

	report := &FateReport{
		GeneratedAt: "2024-01-15 10:30:00",
		Surname:     "张",
		Sex:         "男",
		Born:        "2024年01月15日 10:30",
		TotalNames:  100,
	}

	if fateData != nil {
		if fateData.Bazi != nil {
			report.Bazi = &BaziSection{
				Sizhu:         fateData.Bazi.Sizhu,
				Wuxing:        fateData.Bazi.Wuxing,
				Nayin:         fateData.Bazi.Nayin,
				Zodiac:        fateData.Bazi.Zodiac,
				Constellation: fateData.Bazi.Constellation,
			}
		}
		if fateData.WuxingXiji != nil {
			report.WuXing = &WuXingSection{
				DayGan:     fateData.WuxingXiji.DayGan,
				DayWuxing:  fateData.WuxingXiji.DayWuxing,
				QiangRuo:   fateData.WuxingXiji.QiangRuo,
				XiWuxing:   fateData.WuxingXiji.XiWuxing,
				YongWuxing: fateData.WuxingXiji.YongWuxing,
				JiWuxing:   fateData.WuxingXiji.JiWuxing,
				Analysis:   fateData.WuxingXiji.Analysis,
			}
		}
	}

	report.TopNames = []NameResult{
		{
			Rank:      1,
			FullName:  "张伟安",
			Surname:   "张",
			FirstName: "伟安",
			Strokes:   "11,11,6",
			Char1: CharInfo{
				Char:              "伟",
				TraditionalChar:   "偉",
				Pinyin:            "wěi",
				WuXing:            "土",
				SimplifiedStroke:  6,
				TraditionalStroke: 11,
				ScienceStroke:     11,
				KangxiStroke:      11,
				Radical:           "亻",
				Meaning:           "伟大、壮美",
				IsXiYong:          true,
			},
			Char2: CharInfo{
				Char:              "安",
				TraditionalChar:   "安",
				Pinyin:            "ān",
				WuXing:            "土",
				SimplifiedStroke:  6,
				TraditionalStroke: 6,
				ScienceStroke:     6,
				KangxiStroke:      6,
				Radical:           "宀",
				Meaning:           "安定、平安",
				IsXiYong:          true,
			},
			WuGe: &WuGeResult{
				TianGe: GeItem{Name: "天格", Stroke: 12, Lucky: "凶", DaYan: "掘井无泉", SkyNine: "无理之数，发展薄弱，虽生不足，难酬志向。", YinYangWuXing: "阴木", Analysis: "无理之数，发展薄弱，虽生不足，难酬志向。"},
				RenGe:  GeItem{Name: "人格", Stroke: 22, Lucky: "凶", DaYan: "秋草逢霜", SkyNine: "秋草逢霜，困难疾弱，虽出豪杰，人生波折。", YinYangWuXing: "阴木", Analysis: "秋草逢霜，困难疾弱，虽出豪杰，人生波折。"},
				DiGe:   GeItem{Name: "地格", Stroke: 17, Lucky: "半吉", DaYan: "刚强", SkyNine: "权威刚强，突破万难，如能容忍，必获成功。", YinYangWuXing: "阳金", Analysis: "权威刚强，突破万难，如能容忍，必获成功。"},
				WaiGe:  GeItem{Name: "外格", Stroke: 7, Lucky: "吉", DaYan: "七政之数", SkyNine: "七政之数，精悍严谨，天赋之力，吉星照耀。", YinYangWuXing: "阳火", Analysis: "七政之数，精悍严谨，天赋之力，吉星照耀。"},
				ZongGe: GeItem{Name: "总格", Stroke: 28, Lucky: "凶", DaYan: "阔水浮萍", SkyNine: "遭难之数，豪杰气概，四海漂泊，终世浮躁。", YinYangWuXing: "阴金", Analysis: "遭难之数，豪杰气概，四海漂泊，终世浮躁。"},
			},
			SanCai:       "木木金",
			SanCaiLuck:   "凶多吉少",
			SanCaiDetail: "命运被压抑，导致不良的结果，易生脑部疾病和神经衰弱等，甚至有失财丧命之虑。",
			JiChuYun:     "基础坚实，身适安泰，但天格为3或4时，则内部易产生分离倾向。",
			ChengGongYun: "受上级的引进，得成功顺利发展，基础强固，身心平安，能得长寿，属幸福理想的配置。",
			RenJiGuanXi:  "性格温和，富于理性，能得朋友助力，社交圆满。",
			Score:        85.5,
			Grade:        "上吉",
			ScoreDetail: ScoreDetail{
				WenHuaYinXiang: 80,
				WuXingBaZi:     80,
				ShengXiao:      70,
				WuGeShuLi:      96,
			},
			Interpret: "此名三才配置为木木金，人格为凶数，总格凶，适合八字喜土之人。",
		},
	}

	var buf bytes.Buffer
	f := &TextFormatter{}
	err := f.Format(&buf, report)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "姓名分析报告") {
		t.Error("Output should contain report title")
	}
	if !strings.Contains(output, "张伟安") {
		t.Error("Output should contain the name")
	}
	if !strings.Contains(output, "字基本信息") {
		t.Error("Output should contain char basic info section")
	}
	if !strings.Contains(output, "五格图") {
		t.Error("Output should contain wuge section")
	}
	if !strings.Contains(output, "阴木") {
		t.Error("Output should contain yin-yang wuxing")
	}
	if !strings.Contains(output, "三才解析") {
		t.Error("Output should contain sancai detail")
	}
	if !strings.Contains(output, "基础运") {
		t.Error("Output should contain jichuyun")
	}
	if !strings.Contains(output, "成功运") {
		t.Error("Output should contain chenggongyun")
	}
	if !strings.Contains(output, "人际关系") {
		t.Error("Output should contain renjiguanxi")
	}
	if !strings.Contains(output, "文化") {
		t.Error("Output should contain wenhua score")
	}
	if !strings.Contains(output, "五行八字") {
		t.Error("Output should contain wuxing bazi score")
	}
	if !strings.Contains(output, "偏旁") {
		t.Error("Output should contain radical")
	}
}

func TestMarkdownFormatter(t *testing.T) {
	report := &FateReport{
		GeneratedAt: "2024-01-15 10:30:00",
		Surname:     "张",
		Sex:         "男",
		Born:        "2024年01月15日 10:30",
		TotalNames:  100,
		TopNames: []NameResult{
			{
				Rank:      1,
				FullName:  "张伟安",
				Surname:   "张",
				FirstName: "伟安",
				Strokes:   "11,11,6",
				Char1: CharInfo{
					Char: "伟", TraditionalChar: "偉", Pinyin: "wěi", WuXing: "土",
					SimplifiedStroke: 6, TraditionalStroke: 11, ScienceStroke: 11,
					KangxiStroke: 11, Radical: "亻", Meaning: "伟大", IsXiYong: true,
				},
				Char2: CharInfo{
					Char: "安", TraditionalChar: "安", Pinyin: "ān", WuXing: "土",
					SimplifiedStroke: 6, TraditionalStroke: 6, ScienceStroke: 6,
					KangxiStroke: 6, Radical: "宀", Meaning: "平安", IsXiYong: false,
				},
				WuGe: &WuGeResult{
					TianGe: GeItem{Name: "天格", Stroke: 12, Lucky: "凶", YinYangWuXing: "阴木", DaYan: "掘井无泉", SkyNine: "无理之数"},
					RenGe:  GeItem{Name: "人格", Stroke: 22, Lucky: "凶", YinYangWuXing: "阴木", DaYan: "秋草逢霜", SkyNine: "秋草逢霜"},
					DiGe:   GeItem{Name: "地格", Stroke: 17, Lucky: "半吉", YinYangWuXing: "阳金", DaYan: "刚强", SkyNine: "权威刚强"},
					WaiGe:  GeItem{Name: "外格", Stroke: 7, Lucky: "吉", YinYangWuXing: "阳火", DaYan: "七政之数", SkyNine: "七政之数"},
					ZongGe: GeItem{Name: "总格", Stroke: 28, Lucky: "凶", YinYangWuXing: "阴金", DaYan: "阔水浮萍", SkyNine: "遭难之数"},
				},
				SanCai:       "木木金",
				SanCaiLuck:   "凶多吉少",
				SanCaiDetail: "命运被压抑，导致不良的结果。",
				JiChuYun:     "基础坚实，身适安泰。",
				ChengGongYun: "受上级的引进，得成功顺利发展。",
				RenJiGuanXi:  "性格温和，富于理性。",
				Score:        85.5,
				Grade:        "上吉",
				ScoreDetail:  ScoreDetail{WenHuaYinXiang: 80, WuXingBaZi: 80, ShengXiao: 70, WuGeShuLi: 96},
				Interpret:    "此名三才配置为木木金。",
			},
		},
	}

	var buf bytes.Buffer
	f := &MarkdownFormatter{}
	err := f.Format(&buf, report)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "# 姓名分析报告") {
		t.Error("Output should contain markdown title")
	}
	if !strings.Contains(output, "字基本信息") {
		t.Error("Output should contain char info section")
	}
	if !strings.Contains(output, "五格图") {
		t.Error("Output should contain wuge section")
	}
	if !strings.Contains(output, "运势解析") {
		t.Error("Output should contain yunshi section")
	}
	if !strings.Contains(output, "阴阳五行") {
		t.Error("Output should contain yinyang wuxing column")
	}
	if !strings.Contains(output, "繁体笔画") {
		t.Error("Output should contain traditional stroke column")
	}
	if !strings.Contains(output, "文化印象") {
		t.Error("Output should contain wenhua score")
	}
}

func TestYinYangWuXingAttr(t *testing.T) {
	tests := []struct {
		stroke   int
		expected string
	}{
		{1, "阳木"},
		{2, "阴木"},
		{3, "阳火"},
		{4, "阴火"},
		{5, "阳土"},
		{6, "阴土"},
		{7, "阳金"},
		{8, "阴金"},
		{9, "阳水"},
		{10, "阴水"},
		{11, "阳木"},
		{12, "阴木"},
		{13, "阳火"},
	}
	for _, tt := range tests {
		result := yinYangWuXingAttr(tt.stroke)
		if result != tt.expected {
			t.Errorf("yinYangWuXingAttr(%d) = %s, want %s", tt.stroke, result, tt.expected)
		}
	}
}

func TestSanCaiDetail(t *testing.T) {
	detail := getSanCaiDetail("木木木")
	if detail == "" || detail == "暂无解析" {
		t.Error("Should have detail for 木木木")
	}
	detail = getSanCaiDetail("火火火")
	if detail == "" || detail == "暂无解析" {
		t.Error("Should have detail for 火火火")
	}
}

func TestJiChuYun(t *testing.T) {
	result := getJiChuYun("木", "木")
	if result == "" || result == "暂无解析" {
		t.Error("Should have jichuyun for 木木")
	}
}

func TestChengGongYun(t *testing.T) {
	result := getChengGongYun("木", "火")
	if result == "" || result == "暂无解析" {
		t.Error("Should have chenggongyun for 木火")
	}
}

func TestRenJiGuanXi(t *testing.T) {
	result := getRenJiGuanXi("金", "水")
	if result == "" || result == "暂无解析" {
		t.Error("Should have renjiguanxi for 金水")
	}
}

func TestZodiacWuXing(t *testing.T) {
	tests := []struct {
		zodiac   string
		expected string
	}{
		{"鼠", "水"},
		{"虎", "木"},
		{"龙", "土"},
		{"蛇", "火"},
		{"猴", "金"},
		{"", ""},
	}
	for _, tt := range tests {
		result := getZodiacWuXing(tt.zodiac)
		if result != tt.expected {
			t.Errorf("getZodiacWuXing(%s) = %s, want %s", tt.zodiac, result, tt.expected)
		}
	}
}

func TestCalcZhouYi(t *testing.T) {
	result := CalcZhouYi(11, 0, 13, 12)
	if result == nil {
		t.Fatal("CalcZhouYi should not return nil")
	}
	if result.BenGuaName == "" {
		t.Error("BenGuaName should not be empty")
	}
	if result.BianGuaName == "" {
		t.Error("BianGuaName should not be empty")
	}
	t.Logf("本卦: %s（%s）, 变卦: %s, 动爻: %s, 大象: %s",
		result.BenGuaName, result.BenGuaJiXiong, result.BianGuaName, result.DongYaoDesc, result.DaXiang)
}
