package analysis

import (
	"fmt"
	"io"
	"strings"

	"github.com/godcong/chronos/v2"
	"github.com/babyname/fate/ent"
)

type NameResult struct {
	FullName   string
	Surname    string
	FirstName  string
	WuGe       *WuGeResult
	WuXing     *WuXingResult
	Score      float64
	Interpret  string
}

type WuGeResult struct {
	TianGe  GeItem
	RenGe   GeItem
	DiGe    GeItem
	WaiGe   GeItem
	ZongGe  GeItem
}

type GeItem struct {
	Name     string
	Stroke   int
	Lucky    bool
	ShiShen  string
	Analysis string
}

type WuXingResult struct {
	SanCai     string
	SanCaiLuck string
	DayMaster  string
	XiYong     []string
	JiXing     []string
}

type FateAnalysis struct {
	Bazi       *chronos.BaziInfo
	WuxingXiji *chronos.WuxingXijiInfo
	Names      []NameResult
}

type Formatter interface {
	Format(w io.Writer, analysis *FateAnalysis) error
}

type TextFormatter struct{}

func (f *TextFormatter) Format(w io.Writer, a *FateAnalysis) error {
	fmt.Fprintln(w, strings.Repeat("=", 60))
	fmt.Fprintln(w, "                    姓名分析报告")
	fmt.Fprintln(w, strings.Repeat("=", 60))

	if a.Bazi != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "【八字信息】")
		fmt.Fprintf(w, "  四柱: %s %s %s %s\n", a.Bazi.Sizhu[0], a.Bazi.Sizhu[1], a.Bazi.Sizhu[2], a.Bazi.Sizhu[3])
		fmt.Fprintf(w, "  五行: %s %s %s %s\n", a.Bazi.Wuxing[0], a.Bazi.Wuxing[1], a.Bazi.Wuxing[2], a.Bazi.Wuxing[3])
		fmt.Fprintf(w, "  纳音: %s %s %s %s\n", a.Bazi.Nayin[0], a.Bazi.Nayin[1], a.Bazi.Nayin[2], a.Bazi.Nayin[3])
		fmt.Fprintf(w, "  生肖: %s  星座: %s\n", a.Bazi.Zodiac, a.Bazi.Constellation)
	}

	if a.WuxingXiji != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "【五行喜忌分析】")
		fmt.Fprintf(w, "  日主: %s  五行: %s\n", a.WuxingXiji.DayGan, a.WuxingXiji.DayWuxing)
		fmt.Fprintf(w, "  喜用神: %s\n", strings.Join(a.WuxingXiji.XiWuxing, "、"))
		fmt.Fprintf(w, "  忌神: %s\n", strings.Join(a.WuxingXiji.JiWuxing, "、"))
		fmt.Fprintf(w, "  分析: %s\n", a.WuxingXiji.Analysis)
		if len(a.WuxingXiji.WuXingFen) > 0 {
			fmt.Fprintf(w, "  五行分布: ")
			parts := make([]string, 0, 5)
			for _, wx := range []string{"金", "木", "水", "火", "土"} {
				if v, ok := a.WuxingXiji.WuXingFen[wx]; ok {
					parts = append(parts, fmt.Sprintf("%s%d", wx, v/1000))
				} else {
					parts = append(parts, fmt.Sprintf("%s0", wx))
				}
			}
			fmt.Fprintf(w, "%s\n", strings.Join(parts, " "))
		}
	}

	if len(a.Names) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "【推荐名字】")
		fmt.Fprintln(w, strings.Repeat("-", 60))
		for i, name := range a.Names {
			fmt.Fprintf(w, "\n  %d. %s\n", i+1, name.FullName)
			fmt.Fprintf(w, "     评分: %.1f分\n", name.Score)
			if name.WuGe != nil {
				fmt.Fprintf(w,     "     天格: %d(%v)  人格: %d(%v)  地格: %d(%v)\n",
					name.WuGe.TianGe.Stroke, luckyStr(name.WuGe.TianGe.Lucky),
					name.WuGe.RenGe.Stroke, luckyStr(name.WuGe.RenGe.Lucky),
					name.WuGe.DiGe.Stroke, luckyStr(name.WuGe.DiGe.Lucky))
				fmt.Fprintf(w,     "     外格: %d(%v)  总格: %d(%v)\n",
					name.WuGe.WaiGe.Stroke, luckyStr(name.WuGe.WaiGe.Lucky),
					name.WuGe.ZongGe.Stroke, luckyStr(name.WuGe.ZongGe.Lucky))
			}
			if name.WuXing != nil {
				fmt.Fprintf(w, "     三才: %s(%s)\n", name.WuXing.SanCai, name.WuXing.SanCaiLuck)
				fmt.Fprintf(w, "     喜用: %s\n", strings.Join(name.WuXing.XiYong, "、"))
			}
			if name.Interpret != "" {
				fmt.Fprintf(w, "     解读: %s\n", name.Interpret)
			}
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("=", 60))
	return nil
}

func luckyStr(lucky bool) string {
	if lucky {
		return "吉"
	}
	return "凶"
}

type MarkdownFormatter struct{}

func (f *MarkdownFormatter) Format(w io.Writer, a *FateAnalysis) error {
	fmt.Fprintln(w, "# 姓名分析报告")
	fmt.Fprintln(w)

	if a.Bazi != nil {
		fmt.Fprintln(w, "## 八字信息")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "| 项目 | 年柱 | 月柱 | 日柱 | 时柱 |")
		fmt.Fprintln(w, "|------|------|------|------|------|")
		fmt.Fprintf(w, "| 四柱 | %s | %s | %s | %s |\n", a.Bazi.Sizhu[0], a.Bazi.Sizhu[1], a.Bazi.Sizhu[2], a.Bazi.Sizhu[3])
		fmt.Fprintf(w, "| 五行 | %s | %s | %s | %s |\n", a.Bazi.Wuxing[0], a.Bazi.Wuxing[1], a.Bazi.Wuxing[2], a.Bazi.Wuxing[3])
		fmt.Fprintf(w, "| 纳音 | %s | %s | %s | %s |\n", a.Bazi.Nayin[0], a.Bazi.Nayin[1], a.Bazi.Nayin[2], a.Bazi.Nayin[3])
		fmt.Fprintf(w, "\n**生肖**: %s  **星座**: %s\n\n", a.Bazi.Zodiac, a.Bazi.Constellation)
	}

	if a.WuxingXiji != nil {
		fmt.Fprintln(w, "## 五行喜忌分析")
		fmt.Fprintln(w)
		fmt.Fprintf(w, "- **日主**: %s（五行属%s）\n", a.WuxingXiji.DayGan, a.WuxingXiji.DayWuxing)
		fmt.Fprintf(w, "- **喜用神**: %s\n", strings.Join(a.WuxingXiji.XiWuxing, "、"))
		fmt.Fprintf(w, "- **忌神**: %s\n", strings.Join(a.WuxingXiji.JiWuxing, "、"))
		fmt.Fprintf(w, "- **分析**: %s\n\n", a.WuxingXiji.Analysis)

		if len(a.WuxingXiji.WuXingFen) > 0 {
			fmt.Fprintln(w, "### 五行分布")
			fmt.Fprintln(w)
			fmt.Fprintln(w, "| 五行 | 金 | 木 | 水 | 火 | 土 |")
			fmt.Fprintln(w, "|------|----|----|----|----|----|")
			parts := make([]string, 0, 5)
			for _, wx := range []string{"金", "木", "水", "火", "土"} {
				if v, ok := a.WuxingXiji.WuXingFen[wx]; ok {
					parts = append(parts, fmt.Sprintf("%d", v/1000))
				} else {
					parts = append(parts, "0")
				}
			}
			fmt.Fprintf(w, "| 数量 | %s |\n\n", strings.Join(parts, " | "))
		}
	}

	if len(a.Names) > 0 {
		fmt.Fprintln(w, "## 推荐名字")
		fmt.Fprintln(w)
		for i, name := range a.Names {
			fmt.Fprintf(w, "### %d. %s（%.1f分）\n\n", i+1, name.FullName, name.Score)
			if name.WuGe != nil {
				fmt.Fprintln(w, "| 格 | 笔画 | 吉凶 |")
				fmt.Fprintln(w, "|-----|------|------|")
				fmt.Fprintf(w, "| 天格 | %d | %s |\n", name.WuGe.TianGe.Stroke, luckyStr(name.WuGe.TianGe.Lucky))
				fmt.Fprintf(w, "| 人格 | %d | %s |\n", name.WuGe.RenGe.Stroke, luckyStr(name.WuGe.RenGe.Lucky))
				fmt.Fprintf(w, "| 地格 | %d | %s |\n", name.WuGe.DiGe.Stroke, luckyStr(name.WuGe.DiGe.Lucky))
				fmt.Fprintf(w, "| 外格 | %d | %s |\n", name.WuGe.WaiGe.Stroke, luckyStr(name.WuGe.WaiGe.Lucky))
				fmt.Fprintf(w, "| 总格 | %d | %s |\n\n", name.WuGe.ZongGe.Stroke, luckyStr(name.WuGe.ZongGe.Lucky))
			}
			if name.WuXing != nil {
				fmt.Fprintf(w, "- **三才**: %s（%s）\n", name.WuXing.SanCai, name.WuXing.SanCaiLuck)
				fmt.Fprintf(w, "- **喜用**: %s\n\n", strings.Join(name.WuXing.XiYong, "、"))
			}
			if name.Interpret != "" {
				fmt.Fprintf(w, "> %s\n\n", name.Interpret)
			}
		}
	}

	return nil
}

func NewNameResult(surname string, c1, c2 *ent.Character, fateData *chronos.FateData) NameResult {
	r := NameResult{
		Surname:   surname,
		FirstName: c1.Char + c2.Char,
		FullName:  surname + c1.Char + c2.Char,
	}

	if fateData != nil && fateData.WuxingXiji != nil {
		r.WuXing = &WuXingResult{
			DayMaster: fateData.WuxingXiji.DayGan,
			XiYong:    fateData.WuxingXiji.XiWuxing,
			JiXing:    fateData.WuxingXiji.JiWuxing,
		}
	}

	return r
}
