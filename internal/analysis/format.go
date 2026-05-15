package analysis

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) write(format string, args ...any) {
	if ew.err != nil {
		return
	}
	_, ew.err = fmt.Fprintf(ew.w, format, args...)
}

// Formatter 定义报告格式化输出接口
type Formatter interface {
	Format(w io.Writer, report *FateReport) error
	Extension() string
}

// TextFormatter 以纯文本格式输出报告
type TextFormatter struct{}

// Extension 返回文本格式的文件扩展名
func (f *TextFormatter) Extension() string { return ".txt" }

// Format 将报告格式化为纯文本并写入指定输出流
func (f *TextFormatter) Format(w io.Writer, r *FateReport) error {
	ew := &errWriter{w: w}
	ew.write("%s\n", strings.Repeat("═", 72))
	ew.write("                      姓名分析报告\n")
	ew.write("%s\n", strings.Repeat("═", 72))
	ew.write("  姓氏: %s    性别: %s    出生: %s\n", r.Surname, r.Sex, r.Born)

	if r.Bazi != nil {
		ew.write("\n")
		ew.write("【八字信息】\n")
		ew.write("  四柱: %s  %s  %s  %s\n", r.Bazi.Sizhu[0], r.Bazi.Sizhu[1], r.Bazi.Sizhu[2], r.Bazi.Sizhu[3])
		ew.write("  五行: %s  %s  %s  %s\n", r.Bazi.Wuxing[0], r.Bazi.Wuxing[1], r.Bazi.Wuxing[2], r.Bazi.Wuxing[3])
		ew.write("  纳音: %s  %s  %s  %s\n", r.Bazi.Nayin[0], r.Bazi.Nayin[1], r.Bazi.Nayin[2], r.Bazi.Nayin[3])
		ew.write("  生肖: %s    星座: %s\n", r.Bazi.Zodiac, r.Bazi.Constellation)
	}

	if r.WuXing != nil {
		ew.write("\n")
		ew.write("【五行喜忌分析】\n")
		ew.write("  算法: %s", r.WuXing.MethodName)
		if r.WuXing.GeJuName != "" {
			ew.write("  格局: %s", r.WuXing.GeJuName)
		}
		ew.write("\n")
		ew.write("  日主: %s（%s）    强弱: %s\n", r.WuXing.DayGan, r.WuXing.DayWuxing, r.WuXing.QiangRuo)
		ew.write("  用  神: %s\n", r.WuXing.YongWuxing)
		ew.write("  喜  神: %s\n", strings.Join(r.WuXing.XiWuxing, "、"))
		ew.write("  忌  神: %s\n", strings.Join(r.WuXing.JiWuxing, "、"))
		ew.write("  仇  神: %s\n", strings.Join(r.WuXing.ChouWuxing, "、"))
		ew.write("  闲  神: %s\n", strings.Join(r.WuXing.XianWuxing, "、"))
		ew.write("  分  析: %s\n", r.WuXing.Analysis)
	}

	if len(r.TopNames) > 0 {
		ew.write("\n")
		ew.write("【推荐名字 TOP %d（共 %d 个吉名可选）】\n", len(r.TopNames), r.TotalNames)
		ew.write("%s\n", strings.Repeat("─", 72))

		for _, nm := range r.TopNames {
			ew.write("\n  %2d. %s  （%s %.1f分）\n", nm.Rank, nm.FullName, nm.Grade, nm.Score)
			ew.write("      笔画: %s\n", nm.Strokes)

			ew.write("      ┌─ 字基本信息 ─────────────────────────────────┐\n")
			ew.write("      │ %s  繁:%s  简:%d画  繁:%d画  姓名学:%d画  五行:%s  偏旁:%s  拼音:%s\n",
				nm.Char1.Char, nm.Char1.TraditionalChar, nm.Char1.SimplifiedStroke, nm.Char1.TraditionalStroke, nm.Char1.ScienceStroke, nm.Char1.WuXing, nm.Char1.Radical, nm.Char1.Pinyin)
			ew.write("      │ %s  繁:%s  简:%d画  繁:%d画  姓名学:%d画  五行:%s  偏旁:%s  拼音:%s\n",
				nm.Char2.Char, nm.Char2.TraditionalChar, nm.Char2.SimplifiedStroke, nm.Char2.TraditionalStroke, nm.Char2.ScienceStroke, nm.Char2.WuXing, nm.Char2.Radical, nm.Char2.Pinyin)
			ew.write("      └──────────────────────────────────────────────┘\n")

			if nm.WuGe != nil {
				ew.write("      ┌─ 五格图 ─────────────────────────────────────┐\n")
				ew.write("      │ %s %2d画 %s %s  %s\n", nm.WuGe.TianGe.Name, nm.WuGe.TianGe.Stroke, nm.WuGe.TianGe.YinYangWuXing, nm.WuGe.TianGe.Lucky, nm.WuGe.TianGe.SkyNine)
				ew.write("      │ %s %2d画 %s %s  %s\n", nm.WuGe.RenGe.Name, nm.WuGe.RenGe.Stroke, nm.WuGe.RenGe.YinYangWuXing, nm.WuGe.RenGe.Lucky, nm.WuGe.RenGe.SkyNine)
				ew.write("      │ %s %2d画 %s %s  %s\n", nm.WuGe.DiGe.Name, nm.WuGe.DiGe.Stroke, nm.WuGe.DiGe.YinYangWuXing, nm.WuGe.DiGe.Lucky, nm.WuGe.DiGe.SkyNine)
				ew.write("      │ %s %2d画 %s %s  %s\n", nm.WuGe.WaiGe.Name, nm.WuGe.WaiGe.Stroke, nm.WuGe.WaiGe.YinYangWuXing, nm.WuGe.WaiGe.Lucky, nm.WuGe.WaiGe.SkyNine)
				ew.write("      │ %s %2d画 %s %s  %s\n", nm.WuGe.ZongGe.Name, nm.WuGe.ZongGe.Stroke, nm.WuGe.ZongGe.YinYangWuXing, nm.WuGe.ZongGe.Lucky, nm.WuGe.ZongGe.SkyNine)
				ew.write("      └──────────────────────────────────────────────┘\n")
			}

			ew.write("      三才: %s（%s）\n", nm.SanCai, nm.SanCaiLuck)
			ew.write("      三才解析: %s\n", nm.SanCaiDetail)
			ew.write("      基础运(人地): %s\n", nm.JiChuYun)
			ew.write("      成功运(人天): %s\n", nm.ChengGongYun)
			ew.write("      人际关系(人外): %s\n", nm.RenJiGuanXi)

			if nm.ZhouYi != nil {
				ew.write("      ┌─ 周易卦象 ───────────────────────────────────┐\n")
				ew.write("      │ 本卦: %s（%s）\n", nm.ZhouYi.BenGuaName, nm.ZhouYi.BenGuaJiXiong)
				ew.write("      │ %s\n", nm.ZhouYi.DongYaoDesc)
				ew.write("      │ 变卦: %s\n", nm.ZhouYi.BianGuaName)
				ew.write("      │ 大象: %s\n", nm.ZhouYi.DaXiang)
				if nm.ZhouYi.ShiYe != "" {
					ew.write("      │ 事业: %s\n", truncateStr(nm.ZhouYi.ShiYe, 50))
				}
				if nm.ZhouYi.JingShang != "" {
					ew.write("      │ 经商: %s\n", truncateStr(nm.ZhouYi.JingShang, 50))
				}
				if nm.ZhouYi.QiuMing != "" {
					ew.write("      │ 求名: %s\n", truncateStr(nm.ZhouYi.QiuMing, 50))
				}
				if nm.ZhouYi.HunLian != "" {
					ew.write("      │ 婚恋: %s\n", truncateStr(nm.ZhouYi.HunLian, 50))
				}
				if nm.ZhouYi.JueCe != "" {
					ew.write("      │ 决策: %s\n", truncateStr(nm.ZhouYi.JueCe, 50))
				}
				ew.write("      │ 卦象评分: %d分\n", nm.ZhouYi.Score)
				ew.write("      └──────────────────────────────────────────────┘\n")
			}

			ew.write("      五行: %s%s", nm.Char1.WuXing, nm.Char2.WuXing)
			if nm.Char1.IsXiYong {
				ew.write("  「%s」为喜用", nm.Char1.Char)
			}
			if nm.Char2.IsXiYong {
				ew.write("  「%s」为喜用", nm.Char2.Char)
			}
			ew.write("\n")
			if nm.Char1.Meaning != "" || nm.Char2.Meaning != "" {
				m1 := nm.Char1.Meaning
				m2 := nm.Char2.Meaning
				if len(m1) > 25 {
					m1 = m1[:25] + "…"
				}
				if len(m2) > 25 {
					m2 = m2[:25] + "…"
				}
				ew.write("      字义: %s—%s；%s—%s\n", nm.Char1.Char, m1, nm.Char2.Char, m2)
			}
			ew.write("      评分: %.1f（文化%.0f 五行八字%.0f 生肖%.0f 五格数理%.0f）\n",
				nm.Score, nm.ScoreDetail.WenHuaYinXiang, nm.ScoreDetail.WuXingBaZi, nm.ScoreDetail.ShengXiao, nm.ScoreDetail.WuGeShuLi)
			ew.write("      解读: %s\n", nm.Interpret)
		}
	}

	ew.write("\n")
	ew.write("%s\n", strings.Repeat("═", 72))
	return ew.err
}

// MarkdownFormatter 以 Markdown 格式输出报告
type MarkdownFormatter struct{}

// Extension 返回 Markdown 格式的文件扩展名
func (f *MarkdownFormatter) Extension() string { return ".md" }

// Format 将报告格式化为 Markdown 并写入指定输出流
func (f *MarkdownFormatter) Format(w io.Writer, r *FateReport) error {
	ew := &errWriter{w: w}
	ew.write("# 姓名分析报告 — %s\n\n", r.Surname)
	ew.write("> 生成时间: %s | 性别: %s | 出生: %s\n\n", r.GeneratedAt, r.Sex, r.Born)

	if r.Bazi != nil {
		ew.write("## 八字信息\n\n")
		ew.write("| 项目 | 年柱 | 月柱 | 日柱 | 时柱 |\n")
		ew.write("|------|------|------|------|------|\n")
		ew.write("| 四柱 | %s | %s | %s | %s |\n", r.Bazi.Sizhu[0], r.Bazi.Sizhu[1], r.Bazi.Sizhu[2], r.Bazi.Sizhu[3])
		ew.write("| 五行 | %s | %s | %s | %s |\n", r.Bazi.Wuxing[0], r.Bazi.Wuxing[1], r.Bazi.Wuxing[2], r.Bazi.Wuxing[3])
		ew.write("| 纳音 | %s | %s | %s | %s |\n\n", r.Bazi.Nayin[0], r.Bazi.Nayin[1], r.Bazi.Nayin[2], r.Bazi.Nayin[3])
		ew.write("**生肖**: %s | **星座**: %s\n\n", r.Bazi.Zodiac, r.Bazi.Constellation)
	}

	if r.WuXing != nil {
		ew.write("## 五行喜忌分析\n\n")
		ew.write("| 项目 | 内容 |\n|------|------|\n")
		ew.write("| 算法 | %s |\n", r.WuXing.MethodName)
		if r.WuXing.GeJuName != "" {
			ew.write("| 格局 | %s |\n", r.WuXing.GeJuName)
		}
		ew.write("| 日主 | %s（%s）|\n", r.WuXing.DayGan, r.WuXing.DayWuxing)
		ew.write("| 强弱 | %s |\n", r.WuXing.QiangRuo)
		ew.write("| 用神 | %s |\n", r.WuXing.YongWuxing)
		ew.write("| 喜神 | %s |\n", strings.Join(r.WuXing.XiWuxing, "、"))
		ew.write("| 忌神 | %s |\n", strings.Join(r.WuXing.JiWuxing, "、"))
		ew.write("| 仇神 | %s |\n", strings.Join(r.WuXing.ChouWuxing, "、"))
		ew.write("| 闲神 | %s |\n\n", strings.Join(r.WuXing.XianWuxing, "、"))
		ew.write("> %s\n\n", r.WuXing.Analysis)
	}

	if len(r.TopNames) > 0 {
		ew.write("## 推荐名字 TOP %d（共 %d 个吉名可选）\n\n", len(r.TopNames), r.TotalNames)

		for _, nm := range r.TopNames {
			ew.write("### %d. %s <sup>%s %.1f分</sup>\n\n", nm.Rank, nm.FullName, nm.Grade, nm.Score)
			ew.write("**笔画**: %s | **三才**: %s（%s）| **读音**: %s %s\n\n", nm.Strokes, nm.SanCai, nm.SanCaiLuck, nm.Char1.Pinyin, nm.Char2.Pinyin)

			ew.write("#### 字基本信息\n\n")
			ew.write("| 字 | 繁体 | 简体笔画 | 繁体笔画 | 姓名学笔画 | 五行 | 偏旁 | 拼音 | 喜用 |\n")
			ew.write("|----|------|----------|----------|------------|------|------|------|------|\n")
			xiYong1 := ""
			if nm.Char1.IsXiYong {
				xiYong1 = "✓"
			}
			xiYong2 := ""
			if nm.Char2.IsXiYong {
				xiYong2 = "✓"
			}
			ew.write("| %s | %s | %d | %d | %d | %s | %s | %s | %s |\n", nm.Char1.Char, nm.Char1.TraditionalChar, nm.Char1.SimplifiedStroke, nm.Char1.TraditionalStroke, nm.Char1.ScienceStroke, nm.Char1.WuXing, nm.Char1.Radical, nm.Char1.Pinyin, xiYong1)
			ew.write("| %s | %s | %d | %d | %d | %s | %s | %s | %s |\n\n", nm.Char2.Char, nm.Char2.TraditionalChar, nm.Char2.SimplifiedStroke, nm.Char2.TraditionalStroke, nm.Char2.ScienceStroke, nm.Char2.WuXing, nm.Char2.Radical, nm.Char2.Pinyin, xiYong2)

			if nm.WuGe != nil {
				ew.write("#### 五格图\n\n")
				ew.write("| 格 | 笔画 | 阴阳五行 | 吉凶 | 大衍 | 九星解析 |\n")
				ew.write("|-----|------|----------|------|------|----------|\n")
				items := []GeItem{nm.WuGe.TianGe, nm.WuGe.RenGe, nm.WuGe.DiGe, nm.WuGe.WaiGe, nm.WuGe.ZongGe}
				for _, g := range items {
					ew.write("| %s | %d | %s | %s | %s | %s |\n", g.Name, g.Stroke, g.YinYangWuXing, g.Lucky, g.DaYan, g.SkyNine)
				}
				ew.write("\n")
			}

			ew.write("#### 运势解析\n\n")
			ew.write("| 项目 | 解析 |\n|------|------|\n")
			ew.write("| 三才解析 | %s（%s）|\n", nm.SanCai, nm.SanCaiLuck)
			ew.write("| 三才详解 | %s |\n", nm.SanCaiDetail)
			ew.write("| 基础运(人地) | %s |\n", nm.JiChuYun)
			ew.write("| 成功运(人天) | %s |\n", nm.ChengGongYun)
			ew.write("| 人际关系(人外) | %s |\n\n", nm.RenJiGuanXi)

			if nm.ZhouYi != nil {
				ew.write("#### 周易卦象\n\n")
				ew.write("| 项目 | 内容 |\n|------|------|\n")
				ew.write("| 本卦 | %s（%s）|\n", nm.ZhouYi.BenGuaName, nm.ZhouYi.BenGuaJiXiong)
				ew.write("| 动爻 | %s |\n", nm.ZhouYi.DongYaoDesc)
				ew.write("| 变卦 | %s |\n", nm.ZhouYi.BianGuaName)
				ew.write("| 大象 | %s |\n", nm.ZhouYi.DaXiang)
				ew.write("| 卦象评分 | %d分 |\n\n", nm.ZhouYi.Score)

				if nm.ZhouYi.ShiYe != "" || nm.ZhouYi.HunLian != "" {
					ew.write("| 方面 | 解读 |\n")
					ew.write("|------|------|\n")
					if nm.ZhouYi.ShiYe != "" {
						ew.write("| 事业 | %s |\n", nm.ZhouYi.ShiYe)
					}
					if nm.ZhouYi.JingShang != "" {
						ew.write("| 经商 | %s |\n", nm.ZhouYi.JingShang)
					}
					if nm.ZhouYi.QiuMing != "" {
						ew.write("| 求名 | %s |\n", nm.ZhouYi.QiuMing)
					}
					if nm.ZhouYi.HunLian != "" {
						ew.write("| 婚恋 | %s |\n", nm.ZhouYi.HunLian)
					}
					if nm.ZhouYi.JueCe != "" {
						ew.write("| 决策 | %s |\n", nm.ZhouYi.JueCe)
					}
					ew.write("\n")
				}
			}

			ew.write("#### 评分\n")
			ew.write("| 项目 | 分数 |\n|------|------|\n")
			ew.write("| 文化印象 | %.0f |\n", nm.ScoreDetail.WenHuaYinXiang)
			ew.write("| 五行八字 | %.0f |\n", nm.ScoreDetail.WuXingBaZi)
			ew.write("| 生肖 | %.0f |\n", nm.ScoreDetail.ShengXiao)
			ew.write("| 五格数理 | %.0f |\n", nm.ScoreDetail.WuGeShuLi)
			ew.write("| **综合** | **%.1f** |\n\n", nm.Score)

			if nm.Char1.Meaning != "" || nm.Char2.Meaning != "" {
				m1 := nm.Char1.Meaning
				m2 := nm.Char2.Meaning
				if len(m1) > 40 {
					m1 = m1[:40] + "…"
				}
				if len(m2) > 40 {
					m2 = m2[:40] + "…"
				}
				ew.write("> **字义**: %s—%s；%s—%s\n\n", nm.Char1.Char, m1, nm.Char2.Char, m2)
			}
			ew.write("> %s\n\n", nm.Interpret)
		}
	}

	return ew.err
}

// JSONFormatter 以 JSON 格式输出报告
type JSONFormatter struct {
	Indent bool
}

// Extension 返回 JSON 格式的文件扩展名
func (f *JSONFormatter) Extension() string { return ".json" }

// Format 将报告序列化为 JSON 并写入指定输出流
func (f *JSONFormatter) Format(w io.Writer, r *FateReport) error {
	enc := json.NewEncoder(w)
	if f.Indent {
		enc.SetIndent("", "  ")
	}
	enc.SetEscapeHTML(false)
	return enc.Encode(r)
}
