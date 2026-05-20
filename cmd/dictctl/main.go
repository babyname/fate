package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/babyname/fate/internal/dict"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dictctl <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  import-unihan <path>     Import Unihan IRG Sources")
		fmt.Println("  import-kangxi <path>     Import Kangxi Dictionary")
		fmt.Println("  import-wuxing <path>     Import Wuxing Dictionary")
		fmt.Println("  import-pinyin <path>     Import Pinyin Dictionary")
		fmt.Println("  merge <base.json> <update.json>  Merge two datasets")
		fmt.Println("  validate <path>          Validate dataset")
		fmt.Println("  fix-strokes <path>       Fix science strokes")
		fmt.Println("  stats <path>             Show statistics")
		fmt.Println("  export <path>            Export to JSON")
		fmt.Println("  update <base.json> <update.json>  Update only existing chars (no insert)")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "import-unihan":
		if len(os.Args) < 3 {
			fmt.Println("Usage: dictctl import-unihan <path>")
			os.Exit(1)
		}
		runImportUnihan(os.Args[2])
	case "import-kangxi":
		if len(os.Args) < 3 {
			fmt.Println("Usage: dictctl import-kangxi <path>")
			os.Exit(1)
		}
		runImportKangxi(os.Args[2])
	case "import-wuxing":
		if len(os.Args) < 3 {
			fmt.Println("Usage: dictctl import-wuxing <path>")
			os.Exit(1)
		}
		runImportWuxing(os.Args[2])
	case "import-pinyin":
		if len(os.Args) < 3 {
			fmt.Println("Usage: dictctl import-pinyin <path>")
			os.Exit(1)
		}
		runImportPinyin(os.Args[2])
	case "merge":
		if len(os.Args) < 4 {
			fmt.Println("Usage: dictctl merge <base.json> <update.json>")
			os.Exit(1)
		}
		runMerge(os.Args[2], os.Args[3])
	case "validate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: dictctl validate <path>")
			os.Exit(1)
		}
		runValidate(os.Args[2])
	case "fix-strokes":
		if len(os.Args) < 3 {
			fmt.Println("Usage: dictctl fix-strokes <path>")
			os.Exit(1)
		}
		runFixStrokes(os.Args[2])
	case "stats":
		if len(os.Args) < 3 {
			fmt.Println("Usage: dictctl stats <path>")
			os.Exit(1)
		}
		runStats(os.Args[2])
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: dictctl update <base.json> <update.json>")
			os.Exit(1)
		}
		runUpdate(os.Args[2], os.Args[3])
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func runImportUnihan(path string) {
	entries, err := dict.ParseUnihanIRGSource(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	charEntries := dict.UnihanToCharEntries(entries)
	fmt.Printf("Parsed %d entries from Unihan\n", len(charEntries))

	outputPath := strings.TrimSuffix(path, ".txt") + "_converted.json"
	if err := dict.SaveCharEntriesToJSON(charEntries, outputPath); err != nil {
		fmt.Printf("Error saving: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved to %s\n", outputPath)
}

func runImportKangxi(path string) {
	entries, err := dict.LoadKangxiDict(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded %d Kangxi entries\n", len(entries))

	charEntries := make([]*dict.CharEntry, 0, len(entries))
	for _, k := range entries {
		charEntries = append(charEntries, &dict.CharEntry{
			Char:          k.Char,
			KangxiStroke:  k.Strokes,
			ScienceStroke: k.Strokes,
			Radical:       k.Radical,
			WuXing:        k.WuXing,
			Pinyin:        []string{k.Pinyin},
			Meaning:       k.Meaning,
			Source:        "kangxi",
			SourceConfidence: 0.95,
			Nameable:      true,
		})
	}

	outputPath := strings.TrimSuffix(path, ".json") + "_converted.json"
	if err := dict.SaveCharEntriesToJSON(charEntries, outputPath); err != nil {
		fmt.Printf("Error saving: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved to %s\n", outputPath)
}

func runImportWuxing(path string) {
	entries, err := dict.LoadWuxingDict(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded %d Wuxing entries\n", len(entries))

	charEntries := make([]*dict.CharEntry, 0, len(entries))
	for _, w := range entries {
		charEntries = append(charEntries, &dict.CharEntry{
			Char:     w.Char,
			WuXing:   w.WuXing,
			Source:   w.Source,
			Nameable: true,
		})
	}

	outputPath := strings.TrimSuffix(path, ".json") + "_converted.json"
	if err := dict.SaveCharEntriesToJSON(charEntries, outputPath); err != nil {
		fmt.Printf("Error saving: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved to %s\n", outputPath)
}

func runImportPinyin(path string) {
	entries, err := dict.LoadPinyinDict(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded %d Pinyin entries\n", len(entries))

	charEntries := make([]*dict.CharEntry, 0, len(entries))
	for _, p := range entries {
		charEntries = append(charEntries, &dict.CharEntry{
			Char:     p.Char,
			Pinyin:   p.Pinyin,
			Source:   p.Source,
			Nameable: true,
		})
	}

	outputPath := strings.TrimSuffix(path, ".json") + "_converted.json"
	if err := dict.SaveCharEntriesToJSON(charEntries, outputPath); err != nil {
		fmt.Printf("Error saving: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved to %s\n", outputPath)
}

func runMerge(basePath, updatePath string) {
	base, err := dict.LoadCharEntriesFromJSON(basePath)
	if err != nil {
		fmt.Printf("Error loading base: %v\n", err)
		os.Exit(1)
	}

	update, err := dict.LoadCharEntriesFromJSON(updatePath)
	if err != nil {
		fmt.Printf("Error loading update: %v\n", err)
		os.Exit(1)
	}

	result, merged := dict.MergeEntries(base, update, "merge")
	fmt.Printf("Merge result: Total=%d Updated=%d Inserted=%d Skipped=%d\n",
		result.Total, result.Updated, result.Inserted, result.Skipped)

	if len(result.Conflicts) > 0 {
		fmt.Printf("Conflicts (%d):\n", len(result.Conflicts))
		for _, c := range result.Conflicts {
			fmt.Printf("  %s\n", c)
		}
	}

	if err := dict.SaveCharEntriesToJSON(merged, basePath); err != nil {
		fmt.Printf("Error saving: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved merged result to %s\n", basePath)
}

func runValidate(path string) {
	entries, err := dict.LoadCharEntriesFromJSON(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	result := dict.ValidateEntries(entries)
	fmt.Printf("Validation result: Total=%d Valid=%d\n", result.Total, result.Valid)
	fmt.Printf("Warnings (%d):\n", len(result.Warnings))
	for i, w := range result.Warnings {
		if i >= 20 {
			fmt.Printf("  ... and %d more\n", len(result.Warnings)-20)
			break
		}
		fmt.Printf("  %s\n", w)
	}
	fmt.Printf("Errors (%d):\n", len(result.Errors))
	for _, e := range result.Errors {
		fmt.Printf("  %s\n", e)
	}
}

func runFixStrokes(path string) {
	entries, err := dict.LoadCharEntriesFromJSON(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	updated := dict.ApplyScienceStrokeFix(entries)
	fmt.Printf("Fixed %d science strokes\n", updated)

	if err := dict.SaveCharEntriesToJSON(entries, path); err != nil {
		fmt.Printf("Error saving: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved to %s\n", path)
}

func runStats(path string) {
	entries, err := dict.LoadCharEntriesFromJSON(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	stats := struct {
		Total            int            `json:"total"`
		WithPinyin       int            `json:"with_pinyin"`
		WithKangxiStroke int            `json:"with_kangxi_stroke"`
		WithScienceStroke int           `json:"with_science_stroke"`
		WithWuXing       int            `json:"with_wu_xing"`
		WithRadical      int            `json:"with_radical"`
		WithMeaning      int            `json:"with_meaning"`
		Regular          int            `json:"regular"`
		SimplifiedCount  int            `json:"simplified_count"`
		TraditionalCount int            `json:"traditional_count"`
		KangxiCount      int            `json:"kangxi_count"`
		VariantCount     int            `json:"variant_count"`
		AncientCount     int            `json:"ancient_count"`
		WuXingDist       map[string]int `json:"wuxing_distribution"`
		SourceDist       map[string]int `json:"source_distribution"`
	}{
		WuXingDist: make(map[string]int),
		SourceDist: make(map[string]int),
	}

	for _, e := range entries {
		stats.Total++
		if len(e.Pinyin) > 0 {
			stats.WithPinyin++
		}
		if e.KangxiStroke > 0 {
			stats.WithKangxiStroke++
		}
		if e.ScienceStroke > 0 {
			stats.WithScienceStroke++
		}
		if e.WuXing != "" {
			stats.WithWuXing++
			stats.WuXingDist[e.WuXing]++
		}
		if e.Radical != "" {
			stats.WithRadical++
		}
		if e.Meaning != "" {
			stats.WithMeaning++
		}
		if e.Regular {
			stats.Regular++
		}
		if e.IsSimplified {
			stats.SimplifiedCount++
		}
		if e.IsTraditional {
			stats.TraditionalCount++
		}
		if e.IsKangxi {
			stats.KangxiCount++
		}
		if e.IsVariant {
			stats.VariantCount++
		}
		if e.IsAncient {
			stats.AncientCount++
		}
		if e.Source != "" {
			stats.SourceDist[e.Source]++
		}
	}

	data, _ := json.MarshalIndent(stats, "", "  ")
	fmt.Println(string(data))
}

func runUpdate(basePath, updatePath string) {
	base, err := dict.LoadCharEntriesFromJSON(basePath)
	if err != nil {
		fmt.Printf("Error loading base: %v\n", err)
		os.Exit(1)
	}

	update, err := dict.LoadCharEntriesFromJSON(updatePath)
	if err != nil {
		fmt.Printf("Error loading update: %v\n", err)
		os.Exit(1)
	}

	updateMap := make(map[string]*dict.CharEntry, len(update))
	for _, u := range update {
		updateMap[u.Char] = u
	}

	updated := 0
	meaningUpdated := 0
	pinyinUpdated := 0
	wuxingUpdated := 0
	radicalUpdated := 0
	strokeUpdated := 0

	for _, e := range base {
		u, ok := updateMap[e.Char]
		if !ok {
			continue
		}
		changed := false
		if e.Meaning == "" && u.Meaning != "" {
			e.Meaning = u.Meaning
			meaningUpdated++
			changed = true
		}
		if len(e.Pinyin) == 0 && len(u.Pinyin) > 0 {
			e.Pinyin = u.Pinyin
			pinyinUpdated++
			changed = true
		}
		if e.WuXing == "" && u.WuXing != "" {
			e.WuXing = u.WuXing
			wuxingUpdated++
			changed = true
		}
		if e.Radical == "" && u.Radical != "" {
			e.Radical = u.Radical
			radicalUpdated++
			changed = true
		} else if e.Radical != "" && u.Radical != "" {
			if isNumericRadical(e.Radical) && !isNumericRadical(u.Radical) {
				e.Radical = u.Radical
				radicalUpdated++
				changed = true
			}
		}
		if e.ScienceStroke == 0 && u.ScienceStroke > 0 {
			e.ScienceStroke = u.ScienceStroke
			strokeUpdated++
			changed = true
		}
		if e.KangxiStroke == 0 && u.KangxiStroke > 0 {
			e.KangxiStroke = u.KangxiStroke
			if e.ScienceStroke == 0 {
				e.ScienceStroke = u.KangxiStroke
			}
			strokeUpdated++
			changed = true
		}
		if e.SimplifiedStroke == 0 && u.SimplifiedStroke > 0 {
			e.SimplifiedStroke = u.SimplifiedStroke
			changed = true
		}
		if e.TraditionalStroke == 0 && u.TraditionalStroke > 0 {
			e.TraditionalStroke = u.TraditionalStroke
			changed = true
		}
		if changed {
			updated++
		}
	}

	fmt.Printf("Update result: %d chars updated (meaning=%d pinyin=%d wuxing=%d radical=%d stroke=%d)\n",
		updated, meaningUpdated, pinyinUpdated, wuxingUpdated, radicalUpdated, strokeUpdated)

	if err := dict.SaveCharEntriesToJSON(base, basePath); err != nil {
		fmt.Printf("Error saving: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved updated result to %s\n", basePath)
}

func isNumericRadical(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}
