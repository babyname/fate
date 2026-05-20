package dict

import (
	"testing"
)

func TestCharEntryBasic(t *testing.T) {
	entry := &CharEntry{
		Char:          "明",
		Unicode:       "U+660E",
		IsSimplified:  true,
		IsTraditional: true,
		IsKangxi:      true,
		Pinyin:        []string{"míng"},
		Radical:       "日",
		KangxiStroke:  8,
		ScienceStroke: 8,
		WuXing:        "火",
		Regular:       true,
		CommonLevel:   1,
		Nameable:      true,
	}

	if entry.Char != "明" {
		t.Errorf("expected 明, got %s", entry.Char)
	}
	if entry.KangxiStroke != 8 {
		t.Errorf("expected 8, got %d", entry.KangxiStroke)
	}
	if entry.WuXing != "火" {
		t.Errorf("expected 火, got %s", entry.WuXing)
	}
	if len(entry.Pinyin) != 1 || entry.Pinyin[0] != "míng" {
		t.Errorf("expected [míng], got %v", entry.Pinyin)
	}
	if !entry.IsSimplified || !entry.IsTraditional || !entry.IsKangxi {
		t.Errorf("expected 明 to be simplified+traditional+kangxi")
	}
}

func TestDictIndexBuild(t *testing.T) {
	entries := []*CharEntry{
		{Char: "明", WuXing: "火", ScienceStroke: 8, KangxiStroke: 8, Regular: true, CommonLevel: 1, Nameable: true, IsSimplified: true, IsTraditional: true, IsKangxi: true},
		{Char: "华", WuXing: "水", ScienceStroke: 14, KangxiStroke: 14, Regular: true, CommonLevel: 1, Nameable: true, IsSimplified: true, TraditionalChars: []string{"華"}},
		{Char: "華", WuXing: "木", ScienceStroke: 14, KangxiStroke: 14, Regular: false, Nameable: true, IsTraditional: true, IsKangxi: true, SimplifiedChars: []string{"华"}},
		{Char: "森", WuXing: "木", ScienceStroke: 12, KangxiStroke: 12, Regular: true, CommonLevel: 2, Nameable: true, IsSimplified: true, IsTraditional: true, IsKangxi: true},
		{Char: "龘", WuXing: "火", ScienceStroke: 48, KangxiStroke: 48, Regular: false, Nameable: false, IsTraditional: true},
	}

	idx := NewDictIndex()
	idx.Build(entries)

	char := idx.GetChar('明')
	if char == nil || char.Char != "明" {
		t.Error("expected to find 明")
	}

	char = idx.GetChar('龘')
	if char == nil || char.Nameable != false {
		t.Error("expected 龘 to be not nameable")
	}

	traditional := idx.GetTraditional('华')
	if len(traditional) != 1 || traditional[0] != '華' {
		t.Errorf("expected [華], got %v", traditional)
	}

	simplified := idx.GetSimplified('華')
	if len(simplified) != 1 || simplified[0] != '华' {
		t.Errorf("expected [华], got %v", simplified)
	}

	stroke := idx.GetScienceStroke('明')
	if stroke != 8 {
		t.Errorf("expected 8, got %d", stroke)
	}

	wuxing := idx.GetWuXing('森')
	if wuxing != "木" {
		t.Errorf("expected 木, got %s", wuxing)
	}
}

func TestDictIndexQuery(t *testing.T) {
	entries := []*CharEntry{
		{Char: "明", WuXing: "火", ScienceStroke: 8, KangxiStroke: 8, Regular: true, CommonLevel: 1, Nameable: true, IsSimplified: true, IsTraditional: true, IsKangxi: true, Pinyin: []string{"míng"}},
		{Char: "华", WuXing: "水", ScienceStroke: 14, KangxiStroke: 14, Regular: true, CommonLevel: 1, Nameable: true, IsSimplified: true, Pinyin: []string{"huá"}},
		{Char: "森", WuXing: "木", ScienceStroke: 12, KangxiStroke: 12, Regular: true, CommonLevel: 2, Nameable: true, IsSimplified: true, IsTraditional: true, IsKangxi: true, Pinyin: []string{"sēn"}},
		{Char: "林", WuXing: "木", ScienceStroke: 8, KangxiStroke: 8, Regular: true, CommonLevel: 1, Nameable: true, IsSimplified: true, IsTraditional: true, IsKangxi: true, Pinyin: []string{"lín"}},
		{Char: "龘", WuXing: "火", ScienceStroke: 48, KangxiStroke: 48, Regular: false, Nameable: false, IsTraditional: true},
	}

	idx := NewDictIndex()
	idx.Build(entries)

	result := idx.Query(&QueryFilter{WuXing: "木"})
	if len(result) != 2 {
		t.Errorf("expected 2 木 chars, got %d", len(result))
	}

	result = idx.Query(&QueryFilter{
		WuXing:      "木",
		MinStroke:   5,
		MaxStroke:   10,
		StrokeField: "science",
	})
	if len(result) != 1 {
		t.Errorf("expected 1 木 char with stroke 5-10, got %d", len(result))
	}
	if len(result) > 0 && result[0].Char != "林" {
		t.Errorf("expected 林, got %s", result[0].Char)
	}

	result = idx.Query(&QueryFilter{
		WuXing:      "火",
		RegularOnly: true,
	})
	if len(result) != 1 {
		t.Errorf("expected 1 regular 火 char, got %d", len(result))
	}

	result = idx.Query(&QueryFilter{
		WuXing:       "火",
		NameableOnly: true,
	})
	if len(result) != 1 {
		t.Errorf("expected 1 nameable 火 char, got %d", len(result))
	}

	result = idx.Query(&QueryFilter{
		WuXing:       "木",
		ExcludeChars: map[rune]bool{'林': true},
	})
	if len(result) != 1 {
		t.Errorf("expected 1 木 char (excluding 林), got %d", len(result))
	}

	simplifiedTrue := true
	result = idx.Query(&QueryFilter{
		IsSimplified: &simplifiedTrue,
	})
	if len(result) != 4 {
		t.Errorf("expected 4 simplified chars, got %d", len(result))
	}

	traditionalTrue := true
	result = idx.Query(&QueryFilter{
		IsTraditional: &traditionalTrue,
	})
	if len(result) != 4 {
		t.Errorf("expected 4 traditional chars, got %d", len(result))
	}
}

func TestValidateEntries(t *testing.T) {
	entries := []*CharEntry{
		{Char: "明", Pinyin: []string{"míng"}, KangxiStroke: 8, ScienceStroke: 8, WuXing: "火"},
		{Char: "", Pinyin: []string{"huá"}},
		{Char: "森", WuXing: "木"},
		{Char: "明", WuXing: "水"},
	}

	result := ValidateEntries(entries)

	if result.Total != 4 {
		t.Errorf("expected 4 total, got %d", result.Total)
	}

	if len(result.Errors) == 0 {
		t.Error("expected errors for empty char and duplicate")
	}

	if len(result.Warnings) == 0 {
		t.Error("expected warnings for missing fields")
	}
}

func TestMergeEntries(t *testing.T) {
	base := []*CharEntry{
		{Char: "明", Pinyin: []string{"míng"}, KangxiStroke: 8, WuXing: "火"},
		{Char: "华", WuXing: "水"},
	}

	update := []*CharEntry{
		{Char: "明", Meaning: "光明"},
		{Char: "华", Pinyin: []string{"huá"}, KangxiStroke: 14},
		{Char: "森", WuXing: "木", KangxiStroke: 12},
	}

	result, merged := MergeEntries(base, update, "test")

	if result.Updated != 2 {
		t.Errorf("expected 2 updated, got %d", result.Updated)
	}
	if result.Inserted != 1 {
		t.Errorf("expected 1 inserted, got %d", result.Inserted)
	}

	mingEntry := findEntry(merged, "明")
	if mingEntry == nil {
		t.Fatal("expected to find 明")
	}
	if mingEntry.Meaning != "光明" {
		t.Errorf("expected meaning 光明, got %s", mingEntry.Meaning)
	}

	huaEntry := findEntry(merged, "华")
	if huaEntry == nil {
		t.Fatal("expected to find 华")
	}
	if len(huaEntry.Pinyin) != 1 || huaEntry.Pinyin[0] != "huá" {
		t.Errorf("expected pinyin [huá], got %v", huaEntry.Pinyin)
	}
	if huaEntry.KangxiStroke != 14 {
		t.Errorf("expected kangxi_stroke 14, got %d", huaEntry.KangxiStroke)
	}
}

func TestKangxiStrokeCorrections(t *testing.T) {
	tests := []struct {
		char     rune
		expected int
	}{
		{'东', 8},
		{'华', 14},
		{'国', 11},
		{'乐', 15},
		{'明', 0},
	}

	for _, tt := range tests {
		stroke, ok := GetScienceStrokeCorrection(tt.char)
		if tt.expected == 0 {
			if ok {
				t.Errorf("expected %c to have no correction, got %d", tt.char, stroke)
			}
		} else {
			if !ok {
				t.Errorf("expected %c to have correction", tt.char)
			}
			if stroke != tt.expected {
				t.Errorf("expected %c correction to be %d, got %d", tt.char, tt.expected, stroke)
			}
		}
	}
}

func TestApplyKangxiCorrections(t *testing.T) {
	entries := []*CharEntry{
		{Char: "东", ScienceStroke: 5, KangxiStroke: 5},
		{Char: "华", ScienceStroke: 6, KangxiStroke: 6},
		{Char: "明", ScienceStroke: 8, KangxiStroke: 8},
	}

	updated := ApplyKangxiCorrections(entries)
	if updated != 2 {
		t.Errorf("expected 2 updated, got %d", updated)
	}

	if entries[0].ScienceStroke != 8 {
		t.Errorf("expected 东 science_stroke=8, got %d", entries[0].ScienceStroke)
	}
	if entries[1].ScienceStroke != 14 {
		t.Errorf("expected 华 science_stroke=14, got %d", entries[1].ScienceStroke)
	}
	if entries[2].ScienceStroke != 8 {
		t.Errorf("expected 明 science_stroke=8 (unchanged), got %d", entries[2].ScienceStroke)
	}
}

func TestApplyScienceStrokeFix(t *testing.T) {
	entries := []*CharEntry{
		{Char: "明", KangxiStroke: 8, ScienceStroke: 0},
		{Char: "华", TraditionalStroke: 14, ScienceStroke: 0},
		{Char: "森", SimplifiedStroke: 12, ScienceStroke: 0},
		{Char: "林", ScienceStroke: 8},
	}

	updated := ApplyScienceStrokeFix(entries)
	if updated != 3 {
		t.Errorf("expected 3 updated, got %d", updated)
	}

	if entries[0].ScienceStroke != 8 {
		t.Errorf("expected 明 science_stroke=8 (from kangxi), got %d", entries[0].ScienceStroke)
	}
	if entries[1].ScienceStroke != 14 {
		t.Errorf("expected 华 science_stroke=14 (from traditional), got %d", entries[1].ScienceStroke)
	}
	if entries[2].ScienceStroke != 12 {
		t.Errorf("expected 森 science_stroke=12 (from simplified), got %d", entries[2].ScienceStroke)
	}
	if entries[3].ScienceStroke != 8 {
		t.Errorf("expected 林 science_stroke=8 (unchanged), got %d", entries[3].ScienceStroke)
	}
}

func findEntry(entries []*CharEntry, char string) *CharEntry {
	for _, e := range entries {
		if e.Char == char {
			return e
		}
	}
	return nil
}
