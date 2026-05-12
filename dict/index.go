package dict

import (
	"fmt"
	"sort"
	"sync"
)

type DictIndex struct {
	mu sync.RWMutex

	chars         map[rune]*CharEntry
	byWuXing      map[string][]*CharEntry
	byStroke      map[int][]*CharEntry
	byScienceStroke map[int][]*CharEntry
	byKangxiStroke  map[int][]*CharEntry
	byRadical     map[string][]*CharEntry
	byCommonLevel map[int][]*CharEntry
	bySimplified  map[bool][]*CharEntry
	byTraditional map[bool][]*CharEntry
	byKangxi      map[bool][]*CharEntry
	byVariant     map[bool][]*CharEntry
	byAncient     map[bool][]*CharEntry
	byPinyin      map[string][]*CharEntry

	simplifiedToTraditional map[rune][]rune
	traditionalToSimplified map[rune][]rune
	variantToStandard       map[rune][]rune
	standardToVariant       map[rune][]rune
}

func NewDictIndex() *DictIndex {
	return &DictIndex{
		chars:                   make(map[rune]*CharEntry),
		byWuXing:                make(map[string][]*CharEntry),
		byStroke:                make(map[int][]*CharEntry),
		byScienceStroke:         make(map[int][]*CharEntry),
		byKangxiStroke:          make(map[int][]*CharEntry),
		byRadical:               make(map[string][]*CharEntry),
		byCommonLevel:           make(map[int][]*CharEntry),
		bySimplified:            make(map[bool][]*CharEntry),
		byTraditional:           make(map[bool][]*CharEntry),
		byKangxi:                make(map[bool][]*CharEntry),
		byVariant:               make(map[bool][]*CharEntry),
		byAncient:               make(map[bool][]*CharEntry),
		byPinyin:                make(map[string][]*CharEntry),
		simplifiedToTraditional: make(map[rune][]rune),
		traditionalToSimplified: make(map[rune][]rune),
		variantToStandard:       make(map[rune][]rune),
		standardToVariant:       make(map[rune][]rune),
	}
}

func (idx *DictIndex) Build(entries []*CharEntry) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	for _, e := range entries {
		r := []rune(e.Char)
		if len(r) == 0 {
			continue
		}
		charRune := r[0]

		idx.chars[charRune] = e

		if e.WuXing != "" {
			idx.byWuXing[e.WuXing] = append(idx.byWuXing[e.WuXing], e)
		}
		if e.SimplifiedStroke > 0 {
			idx.byStroke[e.SimplifiedStroke] = append(idx.byStroke[e.SimplifiedStroke], e)
		}
		if e.ScienceStroke > 0 {
			idx.byScienceStroke[e.ScienceStroke] = append(idx.byScienceStroke[e.ScienceStroke], e)
		}
		if e.KangxiStroke > 0 {
			idx.byKangxiStroke[e.KangxiStroke] = append(idx.byKangxiStroke[e.KangxiStroke], e)
		}
		if e.Radical != "" {
			idx.byRadical[e.Radical] = append(idx.byRadical[e.Radical], e)
		}
		if e.CommonLevel > 0 {
			idx.byCommonLevel[e.CommonLevel] = append(idx.byCommonLevel[e.CommonLevel], e)
		}
		idx.bySimplified[e.IsSimplified] = append(idx.bySimplified[e.IsSimplified], e)
		idx.byTraditional[e.IsTraditional] = append(idx.byTraditional[e.IsTraditional], e)
		idx.byKangxi[e.IsKangxi] = append(idx.byKangxi[e.IsKangxi], e)
		idx.byVariant[e.IsVariant] = append(idx.byVariant[e.IsVariant], e)
		idx.byAncient[e.IsAncient] = append(idx.byAncient[e.IsAncient], e)
		for _, py := range e.Pinyin {
			idx.byPinyin[py] = append(idx.byPinyin[py], e)
		}

		for _, tc := range e.TraditionalChars {
			tr := []rune(tc)
			if len(tr) > 0 {
				idx.simplifiedToTraditional[charRune] = appendUniqueRune(idx.simplifiedToTraditional[charRune], tr[0])
				idx.traditionalToSimplified[tr[0]] = appendUniqueRune(idx.traditionalToSimplified[tr[0]], charRune)
			}
		}

		for _, sc := range e.SimplifiedChars {
			sr := []rune(sc)
			if len(sr) > 0 {
				idx.traditionalToSimplified[charRune] = appendUniqueRune(idx.traditionalToSimplified[charRune], sr[0])
				idx.simplifiedToTraditional[sr[0]] = appendUniqueRune(idx.simplifiedToTraditional[sr[0]], charRune)
			}
		}

		for _, vc := range e.VariantChars {
			vr := []rune(vc)
			if len(vr) > 0 {
				idx.variantToStandard[vr[0]] = appendUniqueRune(idx.variantToStandard[vr[0]], charRune)
				idx.standardToVariant[charRune] = appendUniqueRune(idx.standardToVariant[charRune], vr[0])
			}
		}
	}
}

type QueryFilter struct {
	WuXing          string
	MinStroke       int
	MaxStroke       int
	StrokeField     string
	IsSimplified    *bool
	IsTraditional   *bool
	IsKangxi        *bool
	IsVariant       *bool
	IsAncient       *bool
	RegularOnly     bool
	CommonLevel     int
	NameableOnly    bool
	GenderHint      string
	Radical         string
	ExcludeChars    map[rune]bool
}

func (idx *DictIndex) Query(filter *QueryFilter) []*CharEntry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var candidates []*CharEntry

	switch {
	case filter.WuXing != "" && filter.StrokeField != "":
		candidates = idx.queryByWuXingAndStroke(filter)
	case filter.WuXing != "":
		candidates = idx.byWuXing[filter.WuXing]
	case filter.StrokeField != "":
		candidates = idx.queryByStrokeRange(filter)
	default:
		candidates = make([]*CharEntry, 0, len(idx.chars))
		for _, e := range idx.chars {
			candidates = append(candidates, e)
		}
	}

	return idx.applyFilters(candidates, filter)
}

func (idx *DictIndex) queryByWuXingAndStroke(filter *QueryFilter) []*CharEntry {
	wxEntries, ok := idx.byWuXing[filter.WuXing]
	if !ok {
		return nil
	}

	var result []*CharEntry
	for _, e := range wxEntries {
		stroke := idx.getStrokeField(e, filter.StrokeField)
		if stroke >= filter.MinStroke && stroke <= filter.MaxStroke {
			result = append(result, e)
		}
	}
	return result
}

func (idx *DictIndex) queryByStrokeRange(filter *QueryFilter) []*CharEntry {
	var result []*CharEntry
	for stroke := filter.MinStroke; stroke <= filter.MaxStroke; stroke++ {
		var entries []*CharEntry
		switch filter.StrokeField {
		case "science":
			entries = idx.byScienceStroke[stroke]
		case "kangxi":
			entries = idx.byKangxiStroke[stroke]
		default:
			entries = idx.byStroke[stroke]
		}
		result = append(result, entries...)
	}
	return result
}

func (idx *DictIndex) getStrokeField(e *CharEntry, field string) int {
	switch field {
	case "science":
		return e.ScienceStroke
	case "kangxi":
		return e.KangxiStroke
	case "traditional":
		return e.TraditionalStroke
	default:
		return e.SimplifiedStroke
	}
}

func (idx *DictIndex) applyFilters(candidates []*CharEntry, filter *QueryFilter) []*CharEntry {
	var result []*CharEntry
	for _, e := range candidates {
		if filter.RegularOnly && !e.Regular {
			continue
		}
		if filter.NameableOnly && !e.Nameable {
			continue
		}
		if filter.IsSimplified != nil && e.IsSimplified != *filter.IsSimplified {
			continue
		}
		if filter.IsTraditional != nil && e.IsTraditional != *filter.IsTraditional {
			continue
		}
		if filter.IsKangxi != nil && e.IsKangxi != *filter.IsKangxi {
			continue
		}
		if filter.IsVariant != nil && e.IsVariant != *filter.IsVariant {
			continue
		}
		if filter.IsAncient != nil && e.IsAncient != *filter.IsAncient {
			continue
		}
		if filter.CommonLevel > 0 && (e.CommonLevel == 0 || e.CommonLevel > filter.CommonLevel) {
			continue
		}
		if filter.GenderHint != "" && e.GenderHint != "" && e.GenderHint != filter.GenderHint && e.GenderHint != "neutral" {
			continue
		}
		if filter.Radical != "" && e.Radical != filter.Radical {
			continue
		}
		if len(filter.ExcludeChars) > 0 {
			r := []rune(e.Char)
			if len(r) > 0 && filter.ExcludeChars[r[0]] {
				continue
			}
		}
		result = append(result, e)
	}
	return result
}

func (idx *DictIndex) GetChar(c rune) *CharEntry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.chars[c]
}

func (idx *DictIndex) GetTraditional(c rune) []rune {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.simplifiedToTraditional[c]
}

func (idx *DictIndex) GetSimplified(c rune) []rune {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.traditionalToSimplified[c]
}

func (idx *DictIndex) GetVariants(c rune) []rune {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.standardToVariant[c]
}

func (idx *DictIndex) GetScienceStroke(c rune) int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	if e, ok := idx.chars[c]; ok {
		return e.ScienceStroke
	}
	return 0
}

func (idx *DictIndex) GetKangxiStroke(c rune) int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	if e, ok := idx.chars[c]; ok {
		return e.KangxiStroke
	}
	return 0
}

func (idx *DictIndex) GetWuXing(c rune) string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	if e, ok := idx.chars[c]; ok {
		return e.WuXing
	}
	return ""
}

func (idx *DictIndex) Stats() map[string]interface{} {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	stats := map[string]interface{}{
		"total_chars": len(idx.chars),
	}

	wxStats := make(map[string]int)
	for wx, entries := range idx.byWuXing {
		wxStats[wx] = len(entries)
	}
	stats["wuxing_distribution"] = wxStats

	strokeStats := make(map[string]int)
	for s, entries := range idx.byScienceStroke {
		strokeStats[fmt.Sprintf("stroke_%d", s)] = len(entries)
	}
	stats["stroke_distribution"] = strokeStats

	return stats
}

func (idx *DictIndex) SortedByStroke(entries []*CharEntry, strokeField string) []*CharEntry {
	sorted := make([]*CharEntry, len(entries))
	copy(sorted, entries)

	sort.Slice(sorted, func(i, j int) bool {
		si := idx.getStrokeField(sorted[i], strokeField)
		sj := idx.getStrokeField(sorted[j], strokeField)
		return si < sj
	})

	return sorted
}

func appendUniqueRune(slice []rune, r rune) []rune {
	for _, existing := range slice {
		if existing == r {
			return slice
		}
	}
	return append(slice, r)
}
