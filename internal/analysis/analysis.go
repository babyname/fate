package analysis

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/rating"
	"github.com/babyname/fate/internal/wuge"
	"github.com/babyname/fate/internal/wuxing"
	v2 "github.com/godcong/chronos/v2"
)

type GeItem struct {
	Name         string `json:"name"`
	Stroke       int    `json:"stroke"`
	Lucky        string `json:"lucky"`
	DaYan        string `json:"da_yan"`
	SkyNine      string `json:"sky_nine"`
	YinYangWuXing string `json:"yin_yang_wu_xing"`
	Analysis     string `json:"analysis"`
}

type WuGeResult struct {
	TianGe GeItem `json:"tian_ge"`
	RenGe  GeItem `json:"ren_ge"`
	DiGe   GeItem `json:"di_ge"`
	WaiGe  GeItem `json:"wai_ge"`
	ZongGe GeItem `json:"zong_ge"`
}

type CharInfo struct {
	Char              string `json:"char"`
	TraditionalChar   string `json:"traditional_char"`
	Pinyin            string `json:"pinyin"`
	WuXing            string `json:"wu_xing"`
	SimplifiedStroke  int    `json:"simplified_stroke"`
	TraditionalStroke int    `json:"traditional_stroke"`
	ScienceStroke     int    `json:"science_stroke"`
	KangxiStroke      int    `json:"kangxi_stroke"`
	Radical           string `json:"radical"`
	Meaning           string `json:"meaning"`
	IsXiYong          bool   `json:"is_xi_yong"`
}

type ScoreDetail struct {
	WenHuaYinXiang float64 `json:"wen_hua_yin_xiang"`
	WuXingBaZi     float64 `json:"wu_xing_ba_zi"`
	ShengXiao      float64 `json:"sheng_xiao"`
	WuGeShuLi      float64 `json:"wu_ge_shu_li"`
}

type NameResult struct {
	Rank         int          `json:"rank"`
	FullName     string       `json:"full_name"`
	Surname      string       `json:"surname"`
	FirstName    string       `json:"first_name"`
	Strokes      string       `json:"strokes"`
	Char1        CharInfo     `json:"char1"`
	Char2        CharInfo     `json:"char2"`
	WuGe         *WuGeResult  `json:"wu_ge"`
	SanCai       string       `json:"san_cai"`
	SanCaiLuck   string       `json:"san_cai_luck"`
	SanCaiDetail string       `json:"san_cai_detail"`
	JiChuYun     string       `json:"ji_chu_yun"`
	ChengGongYun string       `json:"cheng_gong_yun"`
	RenJiGuanXi  string       `json:"ren_ji_guan_xi"`
	ZhouYi       *ZhouYiResult `json:"zhou_yi"`
	Score        float64      `json:"score"`
	Grade        string       `json:"grade"`
	ScoreDetail  ScoreDetail  `json:"score_detail"`
	Interpret    string       `json:"interpret"`
}

type BaziSection struct {
	Sizhu         [4]string `json:"sizhu"`
	Wuxing        [4]string `json:"wuxing"`
	Nayin         [4]string `json:"nayin"`
	Zodiac        string    `json:"zodiac"`
	Constellation string    `json:"constellation"`
}

type WuXingSection struct {
	DayGan      string   `json:"day_gan"`
	DayWuxing   string   `json:"day_wuxing"`
	QiangRuo    string   `json:"qiang_ruo"`
	XiWuxing    []string `json:"xi_wuxing"`
	YongWuxing  string   `json:"yong_wuxing"`
	JiWuxing    []string `json:"ji_wuxing"`
	ChouWuxing  []string `json:"chou_wuxing"`
	XianWuxing  []string `json:"xian_wuxing"`
	Method      string   `json:"method"`
	MethodName  string   `json:"method_name"`
	GeJuName    string   `json:"geju_name"`
	Analysis    string   `json:"analysis"`
}

type FateReport struct {
	GeneratedAt string         `json:"generated_at"`
	Surname     string         `json:"surname"`
	Born        string         `json:"born"`
	Sex         string         `json:"sex"`
	Bazi        *BaziSection   `json:"bazi"`
	WuXing      *WuXingSection `json:"wu_xing"`
	TotalNames  int            `json:"total_names"`
	TopNames    []NameResult   `json:"top_names"`
}

type Formatter interface {
	Format(w io.Writer, report *FateReport) error
	Extension() string
}

func NewReport(surname, born, sex string, fateData *v2.FateData, totalNames int) *FateReport {
	report := &FateReport{
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
		Surname:     surname,
		Born:        born,
		Sex:         sex,
		TotalNames:  totalNames,
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
				ChouWuxing: fateData.WuxingXiji.ChouWuxing,
				XianWuxing: fateData.WuxingXiji.XianWuxing,
				Method:     fateData.WuxingXiji.MethodName,
				MethodName: fateData.WuxingXiji.MethodName,
				Analysis:   fateData.WuxingXiji.Analysis,
			}
			if fateData.WuxingXiji.GeJu != nil {
				report.WuXing.GeJuName = fateData.WuxingXiji.GeJu.Name
			}
		}
	}

	return report
}

func yinYangWuXingAttr(stroke int) string {
	sanCaiStr := "水木木火火土土金金水"
	yinYangStr := "阴阳"
	wx := string([]rune(sanCaiStr)[stroke%10])
	yy := string([]rune(yinYangStr)[stroke%2])
	return yy + wx
}

func BuildNameResult(rank int, surname string, c1, c2 *ent.Character, l1, l2 int, fateData *v2.FateData) NameResult {
	f1 := c1.ScienceStroke
	f2 := c2.ScienceStroke

	ge := wuge.CalcWuGe(l1, l2, f1, f2)
	tianGe := ge.TianGe()
	renGe := ge.RenGe()
	diGe := ge.DiGe()
	waiGe := ge.WaiGe()
	zongGe := ge.ZongGe()

	tianDaYan := wuge.Find(tianGe)
	renDaYan := wuge.Find(renGe)
	diDaYan := wuge.Find(diGe)
	waiDaYan := wuge.Find(waiGe)
	zongDaYan := wuge.Find(zongGe)

	sanCaiObj := wuxing.NewSanCai(tianGe, renGe, diGe)
	sanCaiStr := sanCaiObj.String()
	sanCaiLuck, _ := wuxing.GetWuXing(sanCaiStr)
	sanCaiDetail := getSanCaiDetail(sanCaiStr)

	tianWuXing := yinYangWuXingAttr(tianGe)
	renWuXing := yinYangWuXingAttr(renGe)
	diWuXing := yinYangWuXingAttr(diGe)
	waiWuXing := yinYangWuXingAttr(waiGe)
	zongWuXing := yinYangWuXingAttr(zongGe)

	renWX := string([]rune(renWuXing)[1:])
	diWX := string([]rune(diWuXing)[1:])
	tianWX := string([]rune(tianWuXing)[1:])
	waiWX := string([]rune(waiWuXing)[1:])

	jiChuYun := getJiChuYun(renWX, diWX)
	chengGongYun := getChengGongYun(renWX, tianWX)
	renJiGuanXi := getRenJiGuanXi(renWX, waiWX)

	p1 := ""
	if len(c1.Pinyin) > 0 {
		p1 = c1.Pinyin[0]
	}
	p2 := ""
	if len(c2.Pinyin) > 0 {
		p2 = c2.Pinyin[0]
	}

	isXiYong1 := false
	isXiYong2 := false
	if fateData != nil && fateData.WuxingXiji != nil {
		for _, wx := range fateData.WuxingXiji.XiWuxing {
			if c1.WuXing == wx {
				isXiYong1 = true
			}
			if c2.WuXing == wx {
				isXiYong2 = true
			}
		}
	}

	var nr *rating.NameRating
	rater := rating.NewRater(fateData)
	nr = rater.RateName(surname, c1, c2)

	strokeStr := fmt.Sprintf("%d,%d,%d", l1, f1, f2)
	if l2 > 0 {
		strokeStr = fmt.Sprintf("%d,%d,%d,%d", l1, l2, f1, f2)
	}

	traditionalChar1 := getTraditionalChar(c1)
	traditionalChar2 := getTraditionalChar(c2)

	zhouyiResult := CalcZhouYi(l1, l2, f1, f2)

	scoreDetail := calcScoreDetail(nr, zhouyiResult, c1, c2, fateData)

	return NameResult{
		Rank:      rank,
		FullName:  surname + c1.Char + c2.Char,
		Surname:   surname,
		FirstName: c1.Char + c2.Char,
		Strokes:   strokeStr,
		Char1: CharInfo{
			Char:              c1.Char,
			TraditionalChar:   traditionalChar1,
			Pinyin:            p1,
			WuXing:            c1.WuXing,
			SimplifiedStroke:  c1.SimplifiedStroke,
			TraditionalStroke: c1.TraditionalStroke,
			ScienceStroke:     c1.ScienceStroke,
			KangxiStroke:      c1.KangxiStroke,
			Radical:           c1.Radical,
			Meaning:           c1.Meaning,
			IsXiYong:          isXiYong1,
		},
		Char2: CharInfo{
			Char:              c2.Char,
			TraditionalChar:   traditionalChar2,
			Pinyin:            p2,
			WuXing:            c2.WuXing,
			SimplifiedStroke:  c2.SimplifiedStroke,
			TraditionalStroke: c2.TraditionalStroke,
			ScienceStroke:     c2.ScienceStroke,
			KangxiStroke:      c2.KangxiStroke,
			Radical:           c2.Radical,
			Meaning:           c2.Meaning,
			IsXiYong:          isXiYong2,
		},
		WuGe: &WuGeResult{
			TianGe: GeItem{
				Name:          "天格",
				Stroke:        tianGe,
				Lucky:         tianDaYan.Lucky,
				DaYan:         tianDaYan.SkyNine,
				SkyNine:       tianDaYan.Comment,
				YinYangWuXing: tianWuXing,
				Analysis:      tianDaYan.Comment,
			},
			RenGe: GeItem{
				Name:          "人格",
				Stroke:        renGe,
				Lucky:         renDaYan.Lucky,
				DaYan:         renDaYan.SkyNine,
				SkyNine:       renDaYan.Comment,
				YinYangWuXing: renWuXing,
				Analysis:      renDaYan.Comment,
			},
			DiGe: GeItem{
				Name:          "地格",
				Stroke:        diGe,
				Lucky:         diDaYan.Lucky,
				DaYan:         diDaYan.SkyNine,
				SkyNine:       diDaYan.Comment,
				YinYangWuXing: diWuXing,
				Analysis:      diDaYan.Comment,
			},
			WaiGe: GeItem{
				Name:          "外格",
				Stroke:        waiGe,
				Lucky:         waiDaYan.Lucky,
				DaYan:         waiDaYan.SkyNine,
				SkyNine:       waiDaYan.Comment,
				YinYangWuXing: waiWuXing,
				Analysis:      waiDaYan.Comment,
			},
			ZongGe: GeItem{
				Name:          "总格",
				Stroke:        zongGe,
				Lucky:         zongDaYan.Lucky,
				DaYan:         zongDaYan.SkyNine,
				SkyNine:       zongDaYan.Comment,
				YinYangWuXing: zongWuXing,
				Analysis:      zongDaYan.Comment,
			},
		},
		SanCai:       sanCaiStr,
		SanCaiLuck:   sanCaiLuck,
		SanCaiDetail: sanCaiDetail,
		JiChuYun:     jiChuYun,
		ChengGongYun: chengGongYun,
		RenJiGuanXi:  renJiGuanXi,
		ZhouYi:       zhouyiResult,
		Score:        nr.TotalScore,
		Grade:        nr.Grade,
		ScoreDetail:  scoreDetail,
		Interpret:    nr.Interpret,
	}
}

func getTraditionalChar(c *ent.Character) string {
	if c.IsTraditional {
		return c.Char
	}
	if trad, ok := simplifiedToTraditional[c.Char]; ok {
		return trad
	}
	return c.Char
}

func calcScoreDetail(nr *rating.NameRating, zyResult *ZhouYiResult, c1, c2 *ent.Character, fateData *v2.FateData) ScoreDetail {
	sd := ScoreDetail{
		WenHuaYinXiang: calcWenHuaScore(c1, c2),
		WuXingBaZi:     nr.WuXingScore,
		ShengXiao:      calcShengXiaoScore(c1, c2, fateData),
		WuGeShuLi:      nr.BiHuaScore,
	}
	return sd
}

func calcWenHuaScore(c1, c2 *ent.Character) float64 {
	score := 70.0
	if c1.Regular {
		score += 5
	}
	if c2.Regular {
		score += 5
	}
	if c1.CommonLevel > 0 && c1.CommonLevel <= 3 {
		score += 5
	}
	if c2.CommonLevel > 0 && c2.CommonLevel <= 3 {
		score += 5
	}
	if c1.Meaning != "" {
		score += 5
	}
	if c2.Meaning != "" {
		score += 5
	}
	if score > 100 {
		score = 100
	}
	return score
}

func calcShengXiaoScore(c1, c2 *ent.Character, fateData *v2.FateData) float64 {
	score := 60.0
	if fateData != nil && fateData.Bazi != nil {
		zodiac := fateData.Bazi.Zodiac
		wx1 := c1.WuXing
		wx2 := c2.WuXing
		zodiacWuXing := getZodiacWuXing(zodiac)
		if zodiacWuXing != "" {
			if isSheng(zodiacWuXing, wx1) || isSheng(wx1, zodiacWuXing) {
				score += 10
			}
			if isSheng(zodiacWuXing, wx2) || isSheng(wx2, zodiacWuXing) {
				score += 10
			}
			if isKe(zodiacWuXing, wx1) {
				score -= 5
			}
			if isKe(zodiacWuXing, wx2) {
				score -= 5
			}
		}
	}
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return score
}

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

var shengMap = map[string]string{
	"木": "火", "火": "土", "土": "金", "金": "水", "水": "木",
}

var keMap = map[string]string{
	"木": "土", "土": "水", "水": "火", "火": "金", "金": "木",
}

func isSheng(a, b string) bool {
	if v, ok := shengMap[a]; ok {
		return v == b
	}
	return false
}

func isKe(a, b string) bool {
	if v, ok := keMap[a]; ok {
		return v == b
	}
	return false
}

func CollectTopNames(names []NameSource, surname string, l1, l2 int, fateData *v2.FateData, topN int, filterFunc func(c1, c2 *ent.Character) bool) []NameResult {
	var all []NameResult
	for _, nm := range names {
		c1 := nm.C1
		c2 := nm.C2

		if filterFunc != nil && !filterFunc(c1, c2) {
			continue
		}

		nr := BuildNameResult(0, surname, c1, c2, l1, l2, fateData)
		all = append(all, nr)
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Score > all[j].Score
	})

	if topN > len(all) {
		topN = len(all)
	}

	result := all[:topN]
	for i := range result {
		result[i].Rank = i + 1
	}
	return result
}

type NameSource struct {
	C1 *ent.Character
	C2 *ent.Character
}

type TextFormatter struct{}

func (f *TextFormatter) Extension() string { return ".txt" }

func (f *TextFormatter) Format(w io.Writer, r *FateReport) error {
	fmt.Fprintln(w, strings.Repeat("═", 72))
	fmt.Fprintln(w, "                      姓名分析报告")
	fmt.Fprintln(w, strings.Repeat("═", 72))
	fmt.Fprintf(w, "  姓氏: %s    性别: %s    出生: %s\n", r.Surname, r.Sex, r.Born)

	if r.Bazi != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "【八字信息】")
		fmt.Fprintf(w, "  四柱: %s  %s  %s  %s\n", r.Bazi.Sizhu[0], r.Bazi.Sizhu[1], r.Bazi.Sizhu[2], r.Bazi.Sizhu[3])
		fmt.Fprintf(w, "  五行: %s  %s  %s  %s\n", r.Bazi.Wuxing[0], r.Bazi.Wuxing[1], r.Bazi.Wuxing[2], r.Bazi.Wuxing[3])
		fmt.Fprintf(w, "  纳音: %s  %s  %s  %s\n", r.Bazi.Nayin[0], r.Bazi.Nayin[1], r.Bazi.Nayin[2], r.Bazi.Nayin[3])
		fmt.Fprintf(w, "  生肖: %s    星座: %s\n", r.Bazi.Zodiac, r.Bazi.Constellation)
	}

	if r.WuXing != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "【五行喜忌分析】")
		fmt.Fprintf(w, "  算法: %s", r.WuXing.MethodName)
		if r.WuXing.GeJuName != "" {
			fmt.Fprintf(w, "  格局: %s", r.WuXing.GeJuName)
		}
		fmt.Fprintln(w)
		fmt.Fprintf(w, "  日主: %s（%s）    强弱: %s\n", r.WuXing.DayGan, r.WuXing.DayWuxing, r.WuXing.QiangRuo)
		fmt.Fprintf(w, "  用  神: %s\n", r.WuXing.YongWuxing)
		fmt.Fprintf(w, "  喜  神: %s\n", strings.Join(r.WuXing.XiWuxing, "、"))
		fmt.Fprintf(w, "  忌  神: %s\n", strings.Join(r.WuXing.JiWuxing, "、"))
		fmt.Fprintf(w, "  仇  神: %s\n", strings.Join(r.WuXing.ChouWuxing, "、"))
		fmt.Fprintf(w, "  闲  神: %s\n", strings.Join(r.WuXing.XianWuxing, "、"))
		fmt.Fprintf(w, "  分  析: %s\n", r.WuXing.Analysis)
	}

	if len(r.TopNames) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "【推荐名字 TOP %d（共 %d 个吉名可选）】\n", len(r.TopNames), r.TotalNames)
		fmt.Fprintln(w, strings.Repeat("─", 72))

		for _, nm := range r.TopNames {
			fmt.Fprintf(w, "\n  %2d. %s  （%s %.1f分）\n", nm.Rank, nm.FullName, nm.Grade, nm.Score)
			fmt.Fprintf(w, "      笔画: %s\n", nm.Strokes)

			fmt.Fprintln(w, "      ┌─ 字基本信息 ─────────────────────────────────┐")
			fmt.Fprintf(w, "      │ %s  繁:%s  简:%d画  繁:%d画  姓名学:%d画  五行:%s  偏旁:%s  拼音:%s\n",
				nm.Char1.Char, nm.Char1.TraditionalChar, nm.Char1.SimplifiedStroke, nm.Char1.TraditionalStroke, nm.Char1.ScienceStroke, nm.Char1.WuXing, nm.Char1.Radical, nm.Char1.Pinyin)
			fmt.Fprintf(w, "      │ %s  繁:%s  简:%d画  繁:%d画  姓名学:%d画  五行:%s  偏旁:%s  拼音:%s\n",
				nm.Char2.Char, nm.Char2.TraditionalChar, nm.Char2.SimplifiedStroke, nm.Char2.TraditionalStroke, nm.Char2.ScienceStroke, nm.Char2.WuXing, nm.Char2.Radical, nm.Char2.Pinyin)
			fmt.Fprintln(w, "      └──────────────────────────────────────────────┘")

			if nm.WuGe != nil {
				fmt.Fprintln(w, "      ┌─ 五格图 ─────────────────────────────────────┐")
				fmt.Fprintf(w, "      │ %s %2d画 %s %s  %s\n", nm.WuGe.TianGe.Name, nm.WuGe.TianGe.Stroke, nm.WuGe.TianGe.YinYangWuXing, nm.WuGe.TianGe.Lucky, nm.WuGe.TianGe.SkyNine)
				fmt.Fprintf(w, "      │ %s %2d画 %s %s  %s\n", nm.WuGe.RenGe.Name, nm.WuGe.RenGe.Stroke, nm.WuGe.RenGe.YinYangWuXing, nm.WuGe.RenGe.Lucky, nm.WuGe.RenGe.SkyNine)
				fmt.Fprintf(w, "      │ %s %2d画 %s %s  %s\n", nm.WuGe.DiGe.Name, nm.WuGe.DiGe.Stroke, nm.WuGe.DiGe.YinYangWuXing, nm.WuGe.DiGe.Lucky, nm.WuGe.DiGe.SkyNine)
				fmt.Fprintf(w, "      │ %s %2d画 %s %s  %s\n", nm.WuGe.WaiGe.Name, nm.WuGe.WaiGe.Stroke, nm.WuGe.WaiGe.YinYangWuXing, nm.WuGe.WaiGe.Lucky, nm.WuGe.WaiGe.SkyNine)
				fmt.Fprintf(w, "      │ %s %2d画 %s %s  %s\n", nm.WuGe.ZongGe.Name, nm.WuGe.ZongGe.Stroke, nm.WuGe.ZongGe.YinYangWuXing, nm.WuGe.ZongGe.Lucky, nm.WuGe.ZongGe.SkyNine)
				fmt.Fprintln(w, "      └──────────────────────────────────────────────┘")
			}

			fmt.Fprintf(w, "      三才: %s（%s）\n", nm.SanCai, nm.SanCaiLuck)
			fmt.Fprintf(w, "      三才解析: %s\n", nm.SanCaiDetail)
			fmt.Fprintf(w, "      基础运(人地): %s\n", nm.JiChuYun)
			fmt.Fprintf(w, "      成功运(人天): %s\n", nm.ChengGongYun)
			fmt.Fprintf(w, "      人际关系(人外): %s\n", nm.RenJiGuanXi)

			if nm.ZhouYi != nil {
				fmt.Fprintln(w, "      ┌─ 周易卦象 ───────────────────────────────────┐")
				fmt.Fprintf(w, "      │ 本卦: %s（%s）\n", nm.ZhouYi.BenGuaName, nm.ZhouYi.BenGuaJiXiong)
				fmt.Fprintf(w, "      │ %s\n", nm.ZhouYi.DongYaoDesc)
				fmt.Fprintf(w, "      │ 变卦: %s\n", nm.ZhouYi.BianGuaName)
				fmt.Fprintf(w, "      │ 大象: %s\n", nm.ZhouYi.DaXiang)
				if nm.ZhouYi.ShiYe != "" {
					fmt.Fprintf(w, "      │ 事业: %s\n", truncateStr(nm.ZhouYi.ShiYe, 50))
				}
				if nm.ZhouYi.JingShang != "" {
					fmt.Fprintf(w, "      │ 经商: %s\n", truncateStr(nm.ZhouYi.JingShang, 50))
				}
				if nm.ZhouYi.QiuMing != "" {
					fmt.Fprintf(w, "      │ 求名: %s\n", truncateStr(nm.ZhouYi.QiuMing, 50))
				}
				if nm.ZhouYi.HunLian != "" {
					fmt.Fprintf(w, "      │ 婚恋: %s\n", truncateStr(nm.ZhouYi.HunLian, 50))
				}
				if nm.ZhouYi.JueCe != "" {
					fmt.Fprintf(w, "      │ 决策: %s\n", truncateStr(nm.ZhouYi.JueCe, 50))
				}
				fmt.Fprintf(w, "      │ 卦象评分: %d分\n", nm.ZhouYi.Score)
				fmt.Fprintln(w, "      └──────────────────────────────────────────────┘")
			}

			fmt.Fprintf(w, "      五行: %s%s", nm.Char1.WuXing, nm.Char2.WuXing)
			if nm.Char1.IsXiYong {
				fmt.Fprintf(w, "  「%s」为喜用", nm.Char1.Char)
			}
			if nm.Char2.IsXiYong {
				fmt.Fprintf(w, "  「%s」为喜用", nm.Char2.Char)
			}
			fmt.Fprintln(w)
			if nm.Char1.Meaning != "" || nm.Char2.Meaning != "" {
				m1 := nm.Char1.Meaning
				m2 := nm.Char2.Meaning
				if len(m1) > 25 {
					m1 = m1[:25] + "…"
				}
				if len(m2) > 25 {
					m2 = m2[:25] + "…"
				}
				fmt.Fprintf(w, "      字义: %s—%s；%s—%s\n", nm.Char1.Char, m1, nm.Char2.Char, m2)
			}
			fmt.Fprintf(w, "      评分: %.1f（文化%.0f 五行八字%.0f 生肖%.0f 五格数理%.0f）\n",
				nm.Score, nm.ScoreDetail.WenHuaYinXiang, nm.ScoreDetail.WuXingBaZi, nm.ScoreDetail.ShengXiao, nm.ScoreDetail.WuGeShuLi)
			fmt.Fprintf(w, "      解读: %s\n", nm.Interpret)
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("═", 72))
	return nil
}

type MarkdownFormatter struct{}

func (f *MarkdownFormatter) Extension() string { return ".md" }

func (f *MarkdownFormatter) Format(w io.Writer, r *FateReport) error {
	fmt.Fprintf(w, "# 姓名分析报告 — %s\n\n", r.Surname)
	fmt.Fprintf(w, "> 生成时间: %s | 性别: %s | 出生: %s\n\n", r.GeneratedAt, r.Sex, r.Born)

	if r.Bazi != nil {
		fmt.Fprintf(w, "## 八字信息\n\n")
		fmt.Fprintln(w, "| 项目 | 年柱 | 月柱 | 日柱 | 时柱 |")
		fmt.Fprintln(w, "|------|------|------|------|------|")
		fmt.Fprintf(w, "| 四柱 | %s | %s | %s | %s |\n", r.Bazi.Sizhu[0], r.Bazi.Sizhu[1], r.Bazi.Sizhu[2], r.Bazi.Sizhu[3])
		fmt.Fprintf(w, "| 五行 | %s | %s | %s | %s |\n", r.Bazi.Wuxing[0], r.Bazi.Wuxing[1], r.Bazi.Wuxing[2], r.Bazi.Wuxing[3])
		fmt.Fprintf(w, "| 纳音 | %s | %s | %s | %s |\n\n", r.Bazi.Nayin[0], r.Bazi.Nayin[1], r.Bazi.Nayin[2], r.Bazi.Nayin[3])
		fmt.Fprintf(w, "**生肖**: %s | **星座**: %s\n\n", r.Bazi.Zodiac, r.Bazi.Constellation)
	}

	if r.WuXing != nil {
		fmt.Fprintf(w, "## 五行喜忌分析\n\n")
		fmt.Fprintf(w, "| 项目 | 内容 |\n|------|------|\n")
		fmt.Fprintf(w, "| 算法 | %s |\n", r.WuXing.MethodName)
		if r.WuXing.GeJuName != "" {
			fmt.Fprintf(w, "| 格局 | %s |\n", r.WuXing.GeJuName)
		}
		fmt.Fprintf(w, "| 日主 | %s（%s）|\n", r.WuXing.DayGan, r.WuXing.DayWuxing)
		fmt.Fprintf(w, "| 强弱 | %s |\n", r.WuXing.QiangRuo)
		fmt.Fprintf(w, "| 用神 | %s |\n", r.WuXing.YongWuxing)
		fmt.Fprintf(w, "| 喜神 | %s |\n", strings.Join(r.WuXing.XiWuxing, "、"))
		fmt.Fprintf(w, "| 忌神 | %s |\n", strings.Join(r.WuXing.JiWuxing, "、"))
		fmt.Fprintf(w, "| 仇神 | %s |\n", strings.Join(r.WuXing.ChouWuxing, "、"))
		fmt.Fprintf(w, "| 闲神 | %s |\n\n", strings.Join(r.WuXing.XianWuxing, "、"))
		fmt.Fprintf(w, "> %s\n\n", r.WuXing.Analysis)
	}

	if len(r.TopNames) > 0 {
		fmt.Fprintf(w, "## 推荐名字 TOP %d（共 %d 个吉名可选）\n\n", len(r.TopNames), r.TotalNames)

		for _, nm := range r.TopNames {
			fmt.Fprintf(w, "### %d. %s <sup>%s %.1f分</sup>\n\n", nm.Rank, nm.FullName, nm.Grade, nm.Score)
			fmt.Fprintf(w, "**笔画**: %s | **三才**: %s（%s）| **读音**: %s %s\n\n", nm.Strokes, nm.SanCai, nm.SanCaiLuck, nm.Char1.Pinyin, nm.Char2.Pinyin)

			fmt.Fprintf(w, "#### 字基本信息\n\n")
			fmt.Fprintln(w, "| 字 | 繁体 | 简体笔画 | 繁体笔画 | 姓名学笔画 | 五行 | 偏旁 | 拼音 | 喜用 |")
			fmt.Fprintln(w, "|----|------|----------|----------|------------|------|------|------|------|")
			xiYong1 := ""
			if nm.Char1.IsXiYong {
				xiYong1 = "✓"
			}
			xiYong2 := ""
			if nm.Char2.IsXiYong {
				xiYong2 = "✓"
			}
			fmt.Fprintf(w, "| %s | %s | %d | %d | %d | %s | %s | %s | %s |\n", nm.Char1.Char, nm.Char1.TraditionalChar, nm.Char1.SimplifiedStroke, nm.Char1.TraditionalStroke, nm.Char1.ScienceStroke, nm.Char1.WuXing, nm.Char1.Radical, nm.Char1.Pinyin, xiYong1)
			fmt.Fprintf(w, "| %s | %s | %d | %d | %d | %s | %s | %s | %s |\n\n", nm.Char2.Char, nm.Char2.TraditionalChar, nm.Char2.SimplifiedStroke, nm.Char2.TraditionalStroke, nm.Char2.ScienceStroke, nm.Char2.WuXing, nm.Char2.Radical, nm.Char2.Pinyin, xiYong2)

			if nm.WuGe != nil {
				fmt.Fprintf(w, "#### 五格图\n\n")
				fmt.Fprintln(w, "| 格 | 笔画 | 阴阳五行 | 吉凶 | 大衍 | 九星解析 |")
				fmt.Fprintln(w, "|-----|------|----------|------|------|----------|")
				items := []GeItem{nm.WuGe.TianGe, nm.WuGe.RenGe, nm.WuGe.DiGe, nm.WuGe.WaiGe, nm.WuGe.ZongGe}
				for _, g := range items {
					fmt.Fprintf(w, "| %s | %d | %s | %s | %s | %s |\n", g.Name, g.Stroke, g.YinYangWuXing, g.Lucky, g.DaYan, g.SkyNine)
				}
				fmt.Fprintln(w)
			}

			fmt.Fprintf(w, "#### 运势解析\n\n")
			fmt.Fprintf(w, "| 项目 | 解析 |\n|------|------|\n")
			fmt.Fprintf(w, "| 三才解析 | %s（%s）|\n", nm.SanCai, nm.SanCaiLuck)
			fmt.Fprintf(w, "| 三才详解 | %s |\n", nm.SanCaiDetail)
			fmt.Fprintf(w, "| 基础运(人地) | %s |\n", nm.JiChuYun)
			fmt.Fprintf(w, "| 成功运(人天) | %s |\n", nm.ChengGongYun)
			fmt.Fprintf(w, "| 人际关系(人外) | %s |\n\n", nm.RenJiGuanXi)

			if nm.ZhouYi != nil {
				fmt.Fprintf(w, "#### 周易卦象\n\n")
				fmt.Fprintf(w, "| 项目 | 内容 |\n|------|------|\n")
				fmt.Fprintf(w, "| 本卦 | %s（%s）|\n", nm.ZhouYi.BenGuaName, nm.ZhouYi.BenGuaJiXiong)
				fmt.Fprintf(w, "| 动爻 | %s |\n", nm.ZhouYi.DongYaoDesc)
				fmt.Fprintf(w, "| 变卦 | %s |\n", nm.ZhouYi.BianGuaName)
				fmt.Fprintf(w, "| 大象 | %s |\n", nm.ZhouYi.DaXiang)
				fmt.Fprintf(w, "| 卦象评分 | %d分 |\n\n", nm.ZhouYi.Score)

				if nm.ZhouYi.ShiYe != "" || nm.ZhouYi.HunLian != "" {
					fmt.Fprintln(w, "| 方面 | 解读 |")
					fmt.Fprintln(w, "|------|------|")
					if nm.ZhouYi.ShiYe != "" {
						fmt.Fprintf(w, "| 事业 | %s |\n", nm.ZhouYi.ShiYe)
					}
					if nm.ZhouYi.JingShang != "" {
						fmt.Fprintf(w, "| 经商 | %s |\n", nm.ZhouYi.JingShang)
					}
					if nm.ZhouYi.QiuMing != "" {
						fmt.Fprintf(w, "| 求名 | %s |\n", nm.ZhouYi.QiuMing)
					}
					if nm.ZhouYi.HunLian != "" {
						fmt.Fprintf(w, "| 婚恋 | %s |\n", nm.ZhouYi.HunLian)
					}
					if nm.ZhouYi.JueCe != "" {
						fmt.Fprintf(w, "| 决策 | %s |\n", nm.ZhouYi.JueCe)
					}
					fmt.Fprintln(w)
				}
			}

			fmt.Fprintf(w, "#### 评分\n")
			fmt.Fprintf(w, "| 项目 | 分数 |\n|------|------|\n")
			fmt.Fprintf(w, "| 文化印象 | %.0f |\n", nm.ScoreDetail.WenHuaYinXiang)
			fmt.Fprintf(w, "| 五行八字 | %.0f |\n", nm.ScoreDetail.WuXingBaZi)
			fmt.Fprintf(w, "| 生肖 | %.0f |\n", nm.ScoreDetail.ShengXiao)
			fmt.Fprintf(w, "| 五格数理 | %.0f |\n", nm.ScoreDetail.WuGeShuLi)
			fmt.Fprintf(w, "| **综合** | **%.1f** |\n\n", nm.Score)

			if nm.Char1.Meaning != "" || nm.Char2.Meaning != "" {
				m1 := nm.Char1.Meaning
				m2 := nm.Char2.Meaning
				if len(m1) > 40 {
					m1 = m1[:40] + "…"
				}
				if len(m2) > 40 {
					m2 = m2[:40] + "…"
				}
				fmt.Fprintf(w, "> **字义**: %s—%s；%s—%s\n\n", nm.Char1.Char, m1, nm.Char2.Char, m2)
			}
			fmt.Fprintf(w, "> %s\n\n", nm.Interpret)
		}
	}

	return nil
}

type JSONFormatter struct {
	Indent bool
}

func (f *JSONFormatter) Extension() string { return ".json" }

func (f *JSONFormatter) Format(w io.Writer, r *FateReport) error {
	enc := json.NewEncoder(w)
	if f.Indent {
		enc.SetIndent("", "  ")
	}
	enc.SetEscapeHTML(false)
	return enc.Encode(r)
}

func truncateStr(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "…"
}
