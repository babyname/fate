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

	a := &FateAnalysis{
		Bazi:       fateData.Bazi,
		WuxingXiji: fateData.WuxingXiji,
		Names: []NameResult{
			{
				FullName:  "张伟安",
				Surname:   "张",
				FirstName: "伟安",
				Score:     85.5,
				WuGe: &WuGeResult{
					TianGe: GeItem{Name: "天格", Stroke: 12, Lucky: true},
					RenGe:  GeItem{Name: "人格", Stroke: 18, Lucky: true},
					DiGe:   GeItem{Name: "地格", Stroke: 14, Lucky: true},
					WaiGe:  GeItem{Name: "外格", Stroke: 8, Lucky: false},
					ZongGe: GeItem{Name: "总格", Stroke: 25, Lucky: true},
				},
				WuXing: &WuXingResult{
					SanCai:     "火土木",
					SanCaiLuck: "半吉",
					DayMaster:  "丁",
					XiYong:     []string{"土", "水"},
					JiXing:     []string{"火", "木"},
				},
				Interpret: "此名三才配置为火土木，人格为吉数，总格大吉，适合八字喜土水之人。",
			},
		},
	}

	var buf bytes.Buffer
	f := &TextFormatter{}
	err := f.Format(&buf, a)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "姓名分析报告") {
		t.Error("Output should contain report title")
	}
	if !strings.Contains(output, "八字信息") {
		t.Error("Output should contain bazi section")
	}
	if !strings.Contains(output, "五行喜忌分析") {
		t.Error("Output should contain wuxing section")
	}
	if !strings.Contains(output, "推荐名字") {
		t.Error("Output should contain name section")
	}
	if !strings.Contains(output, "张伟安") {
		t.Error("Output should contain the name")
	}
	if !strings.Contains(output, "85.5") {
		t.Error("Output should contain the score")
	}
}

func TestMarkdownFormatter(t *testing.T) {
	fateData, _ := v2.GetFateData(&v2.FateInput{
		BirthDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.Local),
		Gender:    1,
		Surname:   "张",
	})

	a := &FateAnalysis{
		Bazi:       fateData.Bazi,
		WuxingXiji: fateData.WuxingXiji,
		Names: []NameResult{
			{
				FullName:  "张伟安",
				Surname:   "张",
				FirstName: "伟安",
				Score:     85.5,
				WuGe: &WuGeResult{
					TianGe: GeItem{Name: "天格", Stroke: 12, Lucky: true},
					RenGe:  GeItem{Name: "人格", Stroke: 18, Lucky: true},
					DiGe:   GeItem{Name: "地格", Stroke: 14, Lucky: true},
					WaiGe:  GeItem{Name: "外格", Stroke: 8, Lucky: false},
					ZongGe: GeItem{Name: "总格", Stroke: 25, Lucky: true},
				},
				WuXing: &WuXingResult{
					SanCai:     "火土木",
					SanCaiLuck: "半吉",
					DayMaster:  "丁",
					XiYong:     []string{"土", "水"},
					JiXing:     []string{"火", "木"},
				},
				Interpret: "此名三才配置为火土木，人格为吉数，总格大吉。",
			},
		},
	}

	var buf bytes.Buffer
	f := &MarkdownFormatter{}
	err := f.Format(&buf, a)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "# 姓名分析报告") {
		t.Error("Output should contain markdown title")
	}
	if !strings.Contains(output, "## 八字信息") {
		t.Error("Output should contain bazi section")
	}
	if !strings.Contains(output, "| 四柱 |") {
		t.Error("Output should contain table header")
	}
	if !strings.Contains(output, "## 推荐名字") {
		t.Error("Output should contain name section")
	}
}
