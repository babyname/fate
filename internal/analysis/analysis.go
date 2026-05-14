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
	Name     string `json:"name"`
	Stroke   int    `json:"stroke"`
	Lucky    string `json:"lucky"`
	DaYan    string `json:"da_yan"`
	SkyNine  string `json:"sky_nine"`
}

type WuGeResult struct {
	TianGe  GeItem `json:"tian_ge"`
	RenGe   GeItem `json:"ren_ge"`
	DiGe    GeItem `json:"di_ge"`
	WaiGe   GeItem `json:"wai_ge"`
	ZongGe  GeItem `json:"zong_ge"`
}

type CharInfo struct {
	Char           string   `json:"char"`
	Pinyin         string   `json:"pinyin"`
	WuXing         string   `json:"wu_xing"`
	ScienceStroke  int      `json:"science_stroke"`
	KangxiStroke   int      `json:"kangxi_stroke"`
	Meaning        string   `json:"meaning"`
	IsXiYong       bool     `json:"is_xi_yong"`
}

type NameResult struct {
	Rank       int         `json:"rank"`
	FullName   string      `json:"full_name"`
	Surname    string      `json:"surname"`
	FirstName  string      `json:"first_name"`
	Strokes    string      `json:"strokes"`
	Char1      CharInfo    `json:"char1"`
	Char2      CharInfo    `json:"char2"`
	WuGe       *WuGeResult `json:"wu_ge"`
	SanCai     string      `json:"san_cai"`
	SanCaiLuck string      `json:"san_cai_luck"`
	Score      float64     `json:"score"`
	Grade      string      `json:"grade"`
	WuXingScore  float64   `json:"wu_xing_score"`
	BiHuaScore   float64   `json:"bi_hua_score"`
	YinYunScore  float64   `json:"yin_yun_score"`
	Interpret  string      `json:"interpret"`
}

type BaziSection struct {
	Sizhu        [4]string `json:"sizhu"`
	Wuxing       [4]string `json:"wuxing"`
	Nayin        [4]string `json:"nayin"`
	Zodiac       string    `json:"zodiac"`
	Constellation string   `json:"constellation"`
}

type WuXingSection struct {
	DayGan       string  `json:"day_gan"`
	DayWuxing    string  `json:"day_wuxing"`
	QiangRuo     string  `json:"qiang_ruo"`
	XiWuxing     []string `json:"xi_wuxing"`
	YongWuxing   string  `json:"yong_wuxing"`
	JiWuxing     []string `json:"ji_wuxing"`
	Analysis     string  `json:"analysis"`
}

type FateReport struct {
	GeneratedAt string        `json:"generated_at"`
	Surname     string        `json:"surname"`
	Born        string        `json:"born"`
	Sex         string        `json:"sex"`
	Bazi        *BaziSection  `json:"bazi"`
	WuXing      *WuXingSection `json:"wu_xing"`
	TotalNames  int           `json:"total_names"`
	TopNames    []NameResult  `json:"top_names"`
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
				Sizhu:        fateData.Bazi.Sizhu,
				Wuxing:       fateData.Bazi.Wuxing,
				Nayin:        fateData.Bazi.Nayin,
				Zodiac:       fateData.Bazi.Zodiac,
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

	return report
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

	sanCai := wuxing.NewSanCai(tianGe, renGe, diGe)
	sanCaiStr := sanCai.String()
	sanCaiLuck, _ := wuxing.GetWuXing(sanCaiStr)

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

	return NameResult{
		Rank:      rank,
		FullName:  surname + c1.Char + c2.Char,
		Surname:   surname,
		FirstName: c1.Char + c2.Char,
		Strokes:   strokeStr,
		Char1: CharInfo{
			Char:          c1.Char,
			Pinyin:        p1,
			WuXing:        c1.WuXing,
			ScienceStroke: c1.ScienceStroke,
			KangxiStroke:  c1.KangxiStroke,
			Meaning:       c1.Meaning,
			IsXiYong:      isXiYong1,
		},
		Char2: CharInfo{
			Char:          c2.Char,
			Pinyin:        p2,
			WuXing:        c2.WuXing,
			ScienceStroke: c2.ScienceStroke,
			KangxiStroke:  c2.KangxiStroke,
			Meaning:       c2.Meaning,
			IsXiYong:      isXiYong2,
		},
		WuGe: &WuGeResult{
			TianGe: GeItem{Name: "天格", Stroke: tianGe, Lucky: tianDaYan.Lucky, DaYan: tianDaYan.SkyNine, SkyNine: tianDaYan.Comment},
			RenGe:  GeItem{Name: "人格", Stroke: renGe, Lucky: renDaYan.Lucky, DaYan: renDaYan.SkyNine, SkyNine: renDaYan.Comment},
			DiGe:   GeItem{Name: "地格", Stroke: diGe, Lucky: diDaYan.Lucky, DaYan: diDaYan.SkyNine, SkyNine: diDaYan.Comment},
			WaiGe:  GeItem{Name: "外格", Stroke: waiGe, Lucky: waiDaYan.Lucky, DaYan: waiDaYan.SkyNine, SkyNine: waiDaYan.Comment},
			ZongGe: GeItem{Name: "总格", Stroke: zongGe, Lucky: zongDaYan.Lucky, DaYan: zongDaYan.SkyNine, SkyNine: zongDaYan.Comment},
		},
		SanCai:     sanCaiStr,
		SanCaiLuck: sanCaiLuck,
		Score:      nr.TotalScore,
		Grade:      nr.Grade,
		WuXingScore:  nr.WuXingScore,
		BiHuaScore:   nr.BiHuaScore,
		YinYunScore:  nr.YinYunScore,
		Interpret:  nr.Interpret,
	}
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
	fmt.Fprintln(w, strings.Repeat("═", 64))
	fmt.Fprintln(w, "                    姓名分析报告")
	fmt.Fprintln(w, strings.Repeat("═", 64))
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
		fmt.Fprintf(w, "  日主: %s（%s）    强弱: %s\n", r.WuXing.DayGan, r.WuXing.DayWuxing, r.WuXing.QiangRuo)
		fmt.Fprintf(w, "  喜用神: %s\n", strings.Join(r.WuXing.XiWuxing, "、"))
		fmt.Fprintf(w, "  用  神: %s\n", r.WuXing.YongWuxing)
		fmt.Fprintf(w, "  忌  神: %s\n", strings.Join(r.WuXing.JiWuxing, "、"))
		fmt.Fprintf(w, "  分  析: %s\n", r.WuXing.Analysis)
	}

	if len(r.TopNames) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "【推荐名字 TOP %d（共 %d 个吉名可选）】\n", len(r.TopNames), r.TotalNames)
		fmt.Fprintln(w, strings.Repeat("─", 64))

		for _, nm := range r.TopNames {
			fmt.Fprintf(w, "\n  %2d. %s  （%s）\n", nm.Rank, nm.FullName, nm.Grade)
			fmt.Fprintf(w, "      笔画: %s\n", nm.Strokes)
			fmt.Fprintf(w, "      ┌───────────────────────────────────────────────┐\n")
			if nm.WuGe != nil {
				fmt.Fprintf(w, "      │ %s %2d  %s  %s\n", nm.WuGe.TianGe.Name, nm.WuGe.TianGe.Stroke, nm.WuGe.TianGe.Lucky, nm.WuGe.TianGe.SkyNine)
				fmt.Fprintf(w, "      │ %s %2d  %s  %s\n", nm.WuGe.RenGe.Name, nm.WuGe.RenGe.Stroke, nm.WuGe.RenGe.Lucky, nm.WuGe.RenGe.SkyNine)
				fmt.Fprintf(w, "      │ %s %2d  %s  %s\n", nm.WuGe.DiGe.Name, nm.WuGe.DiGe.Stroke, nm.WuGe.DiGe.Lucky, nm.WuGe.DiGe.SkyNine)
				fmt.Fprintf(w, "      │ %s %2d  %s  %s\n", nm.WuGe.WaiGe.Name, nm.WuGe.WaiGe.Stroke, nm.WuGe.WaiGe.Lucky, nm.WuGe.WaiGe.SkyNine)
				fmt.Fprintf(w, "      │ %s %2d  %s  %s\n", nm.WuGe.ZongGe.Name, nm.WuGe.ZongGe.Stroke, nm.WuGe.ZongGe.Lucky, nm.WuGe.ZongGe.SkyNine)
			}
			fmt.Fprintf(w, "      └───────────────────────────────────────────────┘\n")
			fmt.Fprintf(w, "      三才: %s（%s）\n", nm.SanCai, nm.SanCaiLuck)
			fmt.Fprintf(w, "      五行: %s%s", nm.Char1.WuXing, nm.Char2.WuXing)
			if nm.Char1.IsXiYong {
				fmt.Fprintf(w, "  「%s」为喜用", nm.Char1.Char)
			}
			if nm.Char2.IsXiYong {
				fmt.Fprintf(w, "  「%s」为喜用", nm.Char2.Char)
			}
			fmt.Fprintln(w)
			fmt.Fprintf(w, "      读音: %s %s\n", nm.Char1.Pinyin, nm.Char2.Pinyin)
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
			fmt.Fprintf(w, "      评分: %.1f（五行%.0f 笔画%.0f 音韵%.0f）\n",
				nm.Score, nm.WuXingScore, nm.BiHuaScore, nm.YinYunScore)
			fmt.Fprintf(w, "      解读: %s\n", nm.Interpret)
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("═", 64))
	return nil
}

type MarkdownFormatter struct{}

func (f *MarkdownFormatter) Extension() string { return ".md" }

func (f *MarkdownFormatter) Format(w io.Writer, r *FateReport) error {
	fmt.Fprintf(w, "# 姓名分析报告 — %s\n\n", r.Surname)
	fmt.Fprintf(w, "> 生成时间: %s | 性别: %s | 出生: %s\n\n", r.GeneratedAt, r.Sex, r.Born)

	if r.Bazi != nil {
		fmt.Fprintln(w, "## 八字信息\n")
		fmt.Fprintln(w, "| 项目 | 年柱 | 月柱 | 日柱 | 时柱 |")
		fmt.Fprintln(w, "|------|------|------|------|------|")
		fmt.Fprintf(w, "| 四柱 | %s | %s | %s | %s |\n", r.Bazi.Sizhu[0], r.Bazi.Sizhu[1], r.Bazi.Sizhu[2], r.Bazi.Sizhu[3])
		fmt.Fprintf(w, "| 五行 | %s | %s | %s | %s |\n", r.Bazi.Wuxing[0], r.Bazi.Wuxing[1], r.Bazi.Wuxing[2], r.Bazi.Wuxing[3])
		fmt.Fprintf(w, "| 纳音 | %s | %s | %s | %s |\n\n", r.Bazi.Nayin[0], r.Bazi.Nayin[1], r.Bazi.Nayin[2], r.Bazi.Nayin[3])
		fmt.Fprintf(w, "**生肖**: %s | **星座**: %s\n\n", r.Bazi.Zodiac, r.Bazi.Constellation)
	}

	if r.WuXing != nil {
		fmt.Fprintln(w, "## 五行喜忌分析\n")
		fmt.Fprintf(w, "| 项目 | 内容 |\n|------|------|\n")
		fmt.Fprintf(w, "| 日主 | %s（%s）|\n", r.WuXing.DayGan, r.WuXing.DayWuxing)
		fmt.Fprintf(w, "| 强弱 | %s |\n", r.WuXing.QiangRuo)
		fmt.Fprintf(w, "| 喜用神 | %s |\n", strings.Join(r.WuXing.XiWuxing, "、"))
		fmt.Fprintf(w, "| 用神 | %s |\n", r.WuXing.YongWuxing)
		fmt.Fprintf(w, "| 忌神 | %s |\n\n", strings.Join(r.WuXing.JiWuxing, "、"))
		fmt.Fprintf(w, "> %s\n\n", r.WuXing.Analysis)
	}

	if len(r.TopNames) > 0 {
		fmt.Fprintf(w, "## 推荐名字 TOP %d（共 %d 个吉名可选）\n\n", len(r.TopNames), r.TotalNames)

		for _, nm := range r.TopNames {
			fmt.Fprintf(w, "### %d. %s <sup>%s %.1f分</sup>\n\n", nm.Rank, nm.FullName, nm.Grade, nm.Score)
			fmt.Fprintf(w, "**笔画**: %s | **三才**: %s（%s）| **读音**: %s %s\n\n", nm.Strokes, nm.SanCai, nm.SanCaiLuck, nm.Char1.Pinyin, nm.Char2.Pinyin)

			if nm.WuGe != nil {
				fmt.Fprintln(w, "| 格 | 笔画 | 吉凶 | 大衍 | 九星 |")
				fmt.Fprintln(w, "|-----|------|------|------|------|")
				items := []GeItem{nm.WuGe.TianGe, nm.WuGe.RenGe, nm.WuGe.DiGe, nm.WuGe.WaiGe, nm.WuGe.ZongGe}
				for _, g := range items {
					fmt.Fprintf(w, "| %s | %d | %s | %s | %s |\n", g.Name, g.Stroke, g.Lucky, g.DaYan, g.SkyNine)
				}
				fmt.Fprintln(w)
			}

			fmt.Fprintf(w, "| 字 | 五行 | 喜用 | 康熙笔画 | 字义 |\n")
			fmt.Fprintf(w, "|----|------|------|----------|------|\n")
			xiYong1 := ""
			if nm.Char1.IsXiYong {
				xiYong1 = "✓"
			}
			xiYong2 := ""
			if nm.Char2.IsXiYong {
				xiYong2 = "✓"
			}
			m1 := nm.Char1.Meaning
			m2 := nm.Char2.Meaning
			if len(m1) > 20 {
				m1 = m1[:20] + "…"
			}
			if len(m2) > 20 {
				m2 = m2[:20] + "…"
			}
			fmt.Fprintf(w, "| %s | %s | %s | %d | %s |\n", nm.Char1.Char, nm.Char1.WuXing, xiYong1, nm.Char1.KangxiStroke, m1)
			fmt.Fprintf(w, "| %s | %s | %s | %d | %s |\n\n", nm.Char2.Char, nm.Char2.WuXing, xiYong2, nm.Char2.KangxiStroke, m2)

			fmt.Fprintf(w, "评分: 五行%.0f + 笔画%.0f + 音韵%.0f = **%.1f**\n\n", nm.WuXingScore, nm.BiHuaScore, nm.YinYunScore, nm.Score)
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
