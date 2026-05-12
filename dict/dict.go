package dict

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type CharEntry struct {
	Char              string   `json:"char"`
	Unicode           string   `json:"unicode"`
	IsSimplified      bool     `json:"is_simplified"`
	IsTraditional     bool     `json:"is_traditional"`
	IsKangxi          bool     `json:"is_kangxi"`
	IsVariant         bool     `json:"is_variant"`
	IsAncient         bool     `json:"is_ancient"`
	Pinyin            []string `json:"pinyin"`
	Radical           string   `json:"radical"`
	RadicalStroke     int      `json:"radical_stroke"`
	SimplifiedStroke  int      `json:"simplified_stroke"`
	TraditionalStroke int      `json:"traditional_stroke"`
	KangxiStroke      int      `json:"kangxi_stroke"`
	ScienceStroke     int      `json:"science_stroke"`
	WuXing            string   `json:"wu_xing"`
	Regular           bool     `json:"regular"`
	CommonLevel       int      `json:"common_level"`
	GenderHint        string   `json:"gender_hint"`
	Nameable          bool     `json:"nameable"`
	Meaning           string   `json:"meaning"`
	Source            string   `json:"source"`
	SourceConfidence  float64  `json:"source_confidence"`
	TraditionalChars  []string `json:"traditional_chars"`
	SimplifiedChars   []string `json:"simplified_chars"`
	VariantChars      []string `json:"variant_chars"`
}

type UnihanEntry struct {
	CodePoint string
	KMandarin string
	KDefinition string
	KTotalStrokes string
	KRSKangXi string
	KTraditionalVariant string
	KSimplifiedVariant string
}

func ParseUnihanIRGSource(path string) (map[rune]*UnihanEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open unihan file: %w", err)
	}
	defer f.Close()

	entries := make(map[rune]*UnihanEntry)
	scanner := bufio.NewScanner(f)
	
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 3 {
			continue
		}

		codeStr := fields[0]
		fieldName := fields[1]
		fieldValue := fields[2]

		codePoint := strings.TrimPrefix(codeStr, "U+")
		r, err := strconv.ParseInt(codePoint, 16, 32)
		if err != nil {
			continue
		}

		runeVal := rune(r)
		if _, ok := entries[runeVal]; !ok {
			entries[runeVal] = &UnihanEntry{CodePoint: codeStr}
		}

		switch fieldName {
		case "kMandarin":
			entries[runeVal].KMandarin = fieldValue
		case "kDefinition":
			entries[runeVal].KDefinition = fieldValue
		case "kTotalStrokes":
			entries[runeVal].KTotalStrokes = fieldValue
		case "kRSKangXi":
			entries[runeVal].KRSKangXi = fieldValue
		case "kTraditionalVariant":
			entries[runeVal].KTraditionalVariant = fieldValue
		case "kSimplifiedVariant":
			entries[runeVal].KSimplifiedVariant = fieldValue
		}
	}

	return entries, scanner.Err()
}

func UnihanToCharEntries(unihan map[rune]*UnihanEntry) []*CharEntry {
	var entries []*CharEntry
	
	for r, u := range unihan {
		if !unicode.Is(unicode.Han, r) {
			continue
		}

		entry := &CharEntry{
			Char:     string(r),
			Unicode:  u.CodePoint,
			IsSimplified: true,
			IsTraditional: true,
			IsKangxi: true,
			Source:   "unihan",
			SourceConfidence: 0.9,
			Nameable: true,
		}

		if u.KMandarin != "" {
			pinyins := strings.Split(u.KMandarin, " ")
			entry.Pinyin = pinyins
		}

		if u.KTotalStrokes != "" {
			strokes := strings.Split(u.KTotalStrokes, " ")
			if s, err := strconv.Atoi(strokes[0]); err == nil {
				entry.SimplifiedStroke = s
			}
		}

		if u.KRSKangXi != "" {
			parts := strings.Split(u.KRSKangXi, " ")
			first := parts[0]
			rsParts := strings.Split(first, ".")
			if len(rsParts) >= 2 {
				if rs, err := strconv.Atoi(rsParts[1]); err == nil {
					entry.KangxiStroke = rs
					entry.ScienceStroke = rs
				}
				entry.Radical = rsParts[0]
			}
		}

		if u.KTraditionalVariant != "" {
			vars := strings.Split(u.KTraditionalVariant, " ")
			for _, v := range vars {
				v = strings.TrimPrefix(v, "U+")
				if code, err := strconv.ParseInt(v, 16, 32); err == nil {
					entry.TraditionalChars = append(entry.TraditionalChars, string(rune(code)))
					entry.IsTraditional = false
				}
			}
		}

		if u.KSimplifiedVariant != "" {
			vars := strings.Split(u.KSimplifiedVariant, " ")
			for _, v := range vars {
				v = strings.TrimPrefix(v, "U+")
				if code, err := strconv.ParseInt(v, 16, 32); err == nil {
					entry.SimplifiedChars = append(entry.SimplifiedChars, string(rune(code)))
					entry.IsSimplified = false
				}
			}
		}

		entries = append(entries, entry)
	}

	return entries
}

func LoadCharEntriesFromJSON(path string) ([]*CharEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open json file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read json file: %w", err)
	}

	var entries []*CharEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	return entries, nil
}

func SaveCharEntriesToJSON(entries []*CharEntry, path string) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write json file: %w", err)
	}

	return nil
}

type KangxiEntry struct {
	Char     string `json:"char"`
	Radical  string `json:"radical"`
	Strokes  int    `json:"strokes"`
	WuXing   string `json:"wu_xing"`
	Pinyin   string `json:"pinyin"`
	Meaning  string `json:"meaning"`
}

func LoadKangxiDict(path string) ([]*KangxiEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open kangxi dict: %w", err)
	}
	defer f.Close()

	var entries []*KangxiEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode kangxi dict: %w", err)
	}

	return entries, nil
}

type WuxingEntry struct {
	Char   string `json:"char"`
	WuXing string `json:"wu_xing"`
	Source string `json:"source"`
}

func LoadWuxingDict(path string) ([]*WuxingEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open wuxing dict: %w", err)
	}
	defer f.Close()

	var entries []*WuxingEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode wuxing dict: %w", err)
	}

	return entries, nil
}

type PinyinEntry struct {
	Char   string   `json:"char"`
	Pinyin []string `json:"pinyin"`
	Source string   `json:"source"`
}

func LoadPinyinDict(path string) ([]*PinyinEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open pinyin dict: %w", err)
	}
	defer f.Close()

	var entries []*PinyinEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode pinyin dict: %w", err)
	}

	return entries, nil
}

type MergeResult struct {
	Total      int
	Updated    int
	Inserted   int
	Skipped    int
	Conflicts  []string
}

func MergeEntries(base []*CharEntry, updates []*CharEntry, source string) (*MergeResult, []*CharEntry) {
	result := &MergeResult{}
	index := make(map[string]*CharEntry)
	for _, e := range base {
		index[e.Char] = e
	}

	for _, u := range updates {
		result.Total++
		if existing, ok := index[u.Char]; ok {
			merged := mergeSingle(existing, u, source)
			if merged {
				result.Updated++
			} else {
				result.Skipped++
			}
		} else {
			u.Source = source
			base = append(base, u)
			index[u.Char] = u
			result.Inserted++
		}
	}

	return result, base
}

func mergeSingle(existing *CharEntry, update *CharEntry, source string) bool {
	changed := false

	if len(update.Pinyin) > 0 && len(existing.Pinyin) == 0 {
		existing.Pinyin = update.Pinyin
		existing.Source = source
		changed = true
	}

	if update.SimplifiedStroke > 0 && existing.SimplifiedStroke == 0 {
		existing.SimplifiedStroke = update.SimplifiedStroke
		changed = true
	}

	if update.TraditionalStroke > 0 && existing.TraditionalStroke == 0 {
		existing.TraditionalStroke = update.TraditionalStroke
		changed = true
	}

	if update.KangxiStroke > 0 && existing.KangxiStroke == 0 {
		existing.KangxiStroke = update.KangxiStroke
		changed = true
	}

	if update.ScienceStroke > 0 && existing.ScienceStroke == 0 {
		existing.ScienceStroke = update.ScienceStroke
		changed = true
	}

	if update.WuXing != "" && existing.WuXing == "" {
		existing.WuXing = update.WuXing
		changed = true
	}

	if update.Radical != "" && existing.Radical == "" {
		existing.Radical = update.Radical
		changed = true
	}

	if update.Meaning != "" && existing.Meaning == "" {
		existing.Meaning = update.Meaning
		changed = true
	}

	if len(update.TraditionalChars) > 0 && len(existing.TraditionalChars) == 0 {
		existing.TraditionalChars = update.TraditionalChars
		changed = true
	}

	if len(update.SimplifiedChars) > 0 && len(existing.SimplifiedChars) == 0 {
		existing.SimplifiedChars = update.SimplifiedChars
		changed = true
	}

	return changed
}

type ValidateResult struct {
	Total     int
	Valid     int
	Warnings  []string
	Errors    []string
}

func ValidateEntries(entries []*CharEntry) *ValidateResult {
	result := &ValidateResult{}
	charSet := make(map[string]int)

	for i, e := range entries {
		result.Total++

		if e.Char == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("entry[%d]: empty char", i))
			continue
		}

		if firstIdx, exists := charSet[e.Char]; exists {
			result.Errors = append(result.Errors, fmt.Sprintf("entry[%d]: duplicate char '%s' (first at %d)", i, e.Char, firstIdx))
		} else {
			charSet[e.Char] = i
		}

		if len(e.Pinyin) == 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("'%s': missing pinyin", e.Char))
		}

		if e.KangxiStroke == 0 && e.ScienceStroke == 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("'%s': missing kangxi_stroke and science_stroke", e.Char))
		}

		if e.WuXing == "" {
			result.Warnings = append(result.Warnings, fmt.Sprintf("'%s': missing wu_xing", e.Char))
		}

		if e.SimplifiedStroke == 0 && e.TraditionalStroke == 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("'%s': missing both simplified and traditional stroke", e.Char))
		}

		if e.IsSimplified && len(e.TraditionalChars) == 0 && !e.IsTraditional {
			result.Warnings = append(result.Warnings, fmt.Sprintf("'%s': simplified char without traditional mapping", e.Char))
		}

		result.Valid++
	}

	return result
}

func ApplyKangxiStrokes(entries []*CharEntry, kangxiDict []*KangxiEntry) int {
	kangxiMap := make(map[string]*KangxiEntry)
	for _, k := range kangxiDict {
		kangxiMap[k.Char] = k
	}

	updated := 0
	for _, e := range entries {
		if k, ok := kangxiMap[e.Char]; ok {
			if e.KangxiStroke == 0 && k.Strokes > 0 {
				e.KangxiStroke = k.Strokes
				e.ScienceStroke = k.Strokes
				e.Source = "kangxi"
				updated++
			}
			if e.WuXing == "" && k.WuXing != "" {
				e.WuXing = k.WuXing
			}
			if e.Radical == "" && k.Radical != "" {
				e.Radical = k.Radical
			}
			if e.Meaning == "" && k.Meaning != "" {
				e.Meaning = k.Meaning
			}
		}
	}
	return updated
}

func ApplyWuxing(entries []*CharEntry, wuxingDict []*WuxingEntry) int {
	wxMap := make(map[string]*WuxingEntry)
	for _, w := range wuxingDict {
		wxMap[w.Char] = w
	}

	updated := 0
	for _, e := range entries {
		if w, ok := wxMap[e.Char]; ok {
			if e.WuXing == "" {
				e.WuXing = w.WuXing
				e.Source = w.Source
				updated++
			}
		}
	}
	return updated
}

func ApplyPinyin(entries []*CharEntry, pinyinDict []*PinyinEntry) int {
	pyMap := make(map[string]*PinyinEntry)
	for _, p := range pinyinDict {
		pyMap[p.Char] = p
	}

	updated := 0
	for _, e := range entries {
		if p, ok := pyMap[e.Char]; ok {
			if len(e.Pinyin) == 0 {
				e.Pinyin = p.Pinyin
				e.Source = p.Source
				updated++
			}
		}
	}
	return updated
}

func ApplyScienceStrokeFix(entries []*CharEntry) int {
	updated := 0
	for _, e := range entries {
		if e.ScienceStroke != 0 {
			continue
		}
		if e.KangxiStroke != 0 {
			e.ScienceStroke = e.KangxiStroke
			updated++
		} else if e.TraditionalStroke != 0 {
			e.ScienceStroke = e.TraditionalStroke
			updated++
		} else if e.SimplifiedStroke != 0 {
			e.ScienceStroke = e.SimplifiedStroke
			updated++
		}
	}
	return updated
}
