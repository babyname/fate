package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/analysis"
	v2 "github.com/godcong/chronos/v2"
)

func main() {
	born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")

	methods := []struct {
		method v2.XiYongMethod
		suffix string
	}{
		{v2.XiYongMethodBalance, "平衡用神法"},
		{v2.XiYongMethodGeJu, "格局用神法"},
	}

	names := []struct {
		c1Char, c1Trad, c1Pinyin, c1WuXing, c1Radical                    string
		c1SimpStroke, c1TradStroke, c1SciStroke, c1KangxiStroke          int
		c1Meaning                                                         string
		c2Char, c2Trad, c2Pinyin, c2WuXing, c2Radical                    string
		c2SimpStroke, c2TradStroke, c2SciStroke, c2KangxiStroke          int
		c2Meaning                                                         string
	}{
		{
			"驰", "馳", "chí", "火", "马",
			6, 13, 13, 13,
			"奔驰、疾驰，象征快速前进",
			"筎", "筎", "rú", "木", "竹",
			12, 12, 12, 12,
			"竹筎，竹子内层的皮，象征坚韧",
		},
		{
			"伟", "偉", "wěi", "土", "亻",
			6, 11, 11, 11,
			"伟大、壮美、卓越",
			"安", "安", "ān", "土", "宀",
			6, 6, 6, 6,
			"安定、平安、安宁",
		},
		{
			"俊", "俊", "jùn", "火", "亻",
			9, 9, 9, 9,
			"俊秀、才智出众",
			"宇", "宇", "yǔ", "土", "宀",
			6, 6, 6, 6,
			"宇宙、屋檐，象征广阔",
		},
		{
			"瑞", "瑞", "ruì", "金", "王",
			13, 14, 14, 14,
			"祥瑞、吉祥",
			"霖", "霖", "lín", "水", "雨",
			16, 16, 16, 16,
			"甘霖、久下不停的雨",
		},
		{
			"梓", "梓", "zǐ", "木", "木",
			11, 11, 11, 11,
			"梓树，象征故乡和生机",
			"萱", "萱", "xuān", "木", "艹",
			12, 15, 15, 15,
			"萱草，忘忧草",
		},
	}

	l1 := 11
	var l2 int

	for _, m := range methods {
		fateData, _ := v2.GetFateData(&v2.FateInput{
			BirthDate: born,
			Gender:    1,
			Surname:   "张",
			Method:    m.method,
		})

		report := analysis.NewReport("张", born.Format("2006年01月02日 15:04"), "男", fateData, 61730)

		for i, nm := range names {
			c1 := &ent.Character{
				Char:              nm.c1Char,
				Pinyin:            []string{nm.c1Pinyin},
				WuXing:            nm.c1WuXing,
				SimplifiedStroke:  nm.c1SimpStroke,
				TraditionalStroke: nm.c1TradStroke,
				ScienceStroke:     nm.c1SciStroke,
				KangxiStroke:      nm.c1KangxiStroke,
				Radical:           nm.c1Radical,
				Meaning:           nm.c1Meaning,
				Regular:           true,
				CommonLevel:       1,
			}
			c2 := &ent.Character{
				Char:              nm.c2Char,
				Pinyin:            []string{nm.c2Pinyin},
				WuXing:            nm.c2WuXing,
				SimplifiedStroke:  nm.c2SimpStroke,
				TraditionalStroke: nm.c2TradStroke,
				ScienceStroke:     nm.c2SciStroke,
				KangxiStroke:      nm.c2KangxiStroke,
				Radical:           nm.c2Radical,
				Meaning:           nm.c2Meaning,
				Regular:           true,
				CommonLevel:       1,
			}

			nr := analysis.BuildNameResult(i+1, "张", c1, c2, l1, l2, fateData)
			report.TopNames = append(report.TopNames, nr)
		}

		outputDir := filepath.Join(".", "output")
		os.MkdirAll(outputDir, os.ModePerm)

		formatters := []analysis.Formatter{
			&analysis.TextFormatter{},
			&analysis.MarkdownFormatter{},
			&analysis.JSONFormatter{Indent: true},
		}

		for _, fmt_ := range formatters {
			filename := "张_姓名分析报告_" + m.suffix + fmt_.Extension()
			fp := filepath.Join(outputDir, filename)

			f, err := os.Create(fp)
			if err != nil {
				panic(err)
			}

			err = fmt_.Format(f, report)
			f.Close()
			if err != nil {
				panic(err)
			}

			println("已生成:", fp)
		}
	}
}
