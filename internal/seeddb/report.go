package seeddb

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (r *Reporter) Generate() error {
	report := DataReport{}

	charFile := filepath.Join(r.seedDir, "character.json")
	if _, err := os.Stat(charFile); err == nil {
		var chars []SeedCharacter
		if err := readJSON(charFile, &chars); err != nil {
			return fmt.Errorf("read characters: %w", err)
		}
		report.Characters = analyzeCharacters(chars)
	}

	wugeFile := filepath.Join(r.seedDir, "wu_ge_lucky.json")
	if _, err := os.Stat(wugeFile); err == nil {
		var wuges []SeedWuGeLucky
		if err := readJSON(wugeFile, &wuges); err != nil {
			return fmt.Errorf("read wu_ge_lucky: %w", err)
		}
		report.WuGeLucky = analyzeWuGeLucky(wuges)
	}

	wuxingFile := filepath.Join(r.seedDir, "wu_xing.json")
	if _, err := os.Stat(wuxingFile); err == nil {
		var wuxings []SeedWuXing
		if err := readJSON(wuxingFile, &wuxings); err != nil {
			return fmt.Errorf("read wu_xing: %w", err)
		}
		report.WuXing = analyzeWuXing(wuxings)
	}

	printReport(report)

	reportFile := filepath.Join(r.seedDir, "report.json")
	if err := writeJSON(reportFile, report); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	log.Printf("Report saved to %s", reportFile)

	return nil
}

func analyzeCharacters(chars []SeedCharacter) CharacterReport {
	rpt := CharacterReport{
		Total:      len(chars),
		WuXingDist: make(map[string]int),
	}

	for _, c := range chars {
		if c.WuXing != "" {
			rpt.WithWuXing++
			rpt.WuXingDist[c.WuXing]++
		} else {
			rpt.WithoutWuXing++
		}

		if len(c.Pinyin) > 0 {
			rpt.WithPinyin++
		}

		if c.Regular {
			rpt.RegularCount++
		}
		if c.Nameable {
			rpt.NameableCount++
		}
		if c.IsSimplified {
			rpt.SimplifiedCount++
		}
		if c.IsTraditional {
			rpt.TraditionalCount++
		}
		if c.IsKangxi {
			rpt.KangxiCount++
		}
		if c.IsVariant {
			rpt.VariantCount++
		}

		if c.ScienceStroke <= 0 && c.Nameable {
			rpt.StrokeIssues = append(rpt.StrokeIssues, StrokeIssue{
				Char:    c.Char,
				Field:   "science_stroke",
				Value:   c.ScienceStroke,
				Message: "nameable character missing science_stroke",
			})
		}
	}

	if rpt.Total > 0 {
		rpt.WuXingCoverage = float64(rpt.WithWuXing) / float64(rpt.Total) * 100
		rpt.PinyinCoverage = float64(rpt.WithPinyin) / float64(rpt.Total) * 100
	}

	return rpt
}

func analyzeWuGeLucky(wuges []SeedWuGeLucky) WuGeLuckyReport {
	rpt := WuGeLuckyReport{Total: len(wuges)}

	for _, w := range wuges {
		if w.ZongLucky {
			rpt.LuckyCount++
		}
		if w.ZongMax {
			rpt.MaxCount++
		}
		if w.ZongSex {
			rpt.SexCount++
		}
	}

	if rpt.Total > 0 {
		rpt.LuckyRate = float64(rpt.LuckyCount) / float64(rpt.Total) * 100
	}

	return rpt
}

func analyzeWuXing(wuxings []SeedWuXing) WuXingReport {
	rpt := WuXingReport{
		Total:       len(wuxings),
		FortuneDist: make(map[string]int),
	}

	for _, w := range wuxings {
		if strings.HasPrefix(w.Fortune, "吉") {
			rpt.LuckyCount++
		} else if strings.HasPrefix(w.Fortune, "凶") {
			rpt.UnluckyCount++
		}
		key := w.Fortune
		if idx := strings.Index(key, "|"); idx >= 0 {
			key = key[:idx]
		}
		rpt.FortuneDist[key]++
	}

	return rpt
}

func printReport(rpt DataReport) {
	fmt.Println("=== Data Quality Report ===")
	fmt.Println()

	fmt.Println("--- Characters ---")
	fmt.Printf("  Total:            %d\n", rpt.Characters.Total)
	fmt.Printf("  With WuXing:      %d (%.1f%%)\n", rpt.Characters.WithWuXing, rpt.Characters.WuXingCoverage)
	fmt.Printf("  Without WuXing:   %d\n", rpt.Characters.WithoutWuXing)
	fmt.Printf("  With Pinyin:      %d (%.1f%%)\n", rpt.Characters.WithPinyin, rpt.Characters.PinyinCoverage)
	fmt.Printf("  Regular:          %d\n", rpt.Characters.RegularCount)
	fmt.Printf("  Nameable:         %d\n", rpt.Characters.NameableCount)
	fmt.Printf("  Simplified:       %d\n", rpt.Characters.SimplifiedCount)
	fmt.Printf("  Traditional:      %d\n", rpt.Characters.TraditionalCount)
	fmt.Printf("  Kangxi:           %d\n", rpt.Characters.KangxiCount)
	fmt.Printf("  Variant:          %d\n", rpt.Characters.VariantCount)

	if len(rpt.Characters.WuXingDist) > 0 {
		fmt.Println("  WuXing Distribution:")
		for wx, count := range rpt.Characters.WuXingDist {
			fmt.Printf("    %s: %d\n", wx, count)
		}
	}

	if len(rpt.Characters.StrokeIssues) > 0 {
		fmt.Printf("  Stroke Issues:    %d\n", len(rpt.Characters.StrokeIssues))
		if len(rpt.Characters.StrokeIssues) <= 10 {
			for _, si := range rpt.Characters.StrokeIssues {
				fmt.Printf("    [%s] %s = %d: %s\n", si.Char, si.Field, si.Value, si.Message)
			}
		}
	}

	fmt.Println()
	fmt.Println("--- WuGeLucky ---")
	fmt.Printf("  Total:            %d\n", rpt.WuGeLucky.Total)
	fmt.Printf("  Lucky:            %d (%.1f%%)\n", rpt.WuGeLucky.LuckyCount, rpt.WuGeLucky.LuckyRate)
	fmt.Printf("  Max:              %d\n", rpt.WuGeLucky.MaxCount)
	fmt.Printf("  Sex-specific:     %d\n", rpt.WuGeLucky.SexCount)

	fmt.Println()
	fmt.Println("--- WuXing ---")
	fmt.Printf("  Total:            %d\n", rpt.WuXing.Total)
	fmt.Printf("  Lucky:            %d\n", rpt.WuXing.LuckyCount)
	fmt.Printf("  Unlucky:          %d\n", rpt.WuXing.UnluckyCount)

	if len(rpt.WuXing.FortuneDist) > 0 {
		fmt.Println("  Fortune Distribution:")
		for f, count := range rpt.WuXing.FortuneDist {
			fmt.Printf("    %s: %d\n", f, count)
		}
	}

	fmt.Println()
	fmt.Println("=== End Report ===")
}
