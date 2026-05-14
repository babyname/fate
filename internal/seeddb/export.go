package seeddb

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/sqlite3ent/sqlite3"
)

type oldCharacter struct {
	Hash                     string
	PinYin                   string
	Ch                       string
	Radical                  string
	RadicalStroke            int
	Stroke                   int
	IsKangXi                 bool
	KangXi                   string
	KangXiStroke             int
	SimpleRadical            string
	SimpleRadicalStroke      int
	SimpleTotalStroke        int
	TraditionalRadical       string
	TraditionalRadicalStroke int
	TraditionalTotalStroke   int
	NameScience              bool
	WuXing                   string
	Lucky                    string
	Regular                  bool
	TraditionalCharacter     string
	VariantCharacter         string
	Comment                  string
	ScienceStroke            int
}

type oldNCharacter struct {
	ID             int
	PinYin         string
	Char           string
	CharStroke     int
	Radical        string
	RadicalStroke  int
	IsRegular      bool
	IsSimplified   bool
	SimplifiedID   string
	IsTraditional  bool
	TraditionalID  string
	IsKangXi       bool
	KangXiID       string
	KangXiStroke   int
	IsVariant      bool
	VariantID      string
	IsScience      bool
	ScienceStroke  int
	WuXing         string
	Lucky          string
	Explanation    string
	Comment        string
	NeedFix        bool
}

type oldWuGeLucky struct {
	LastStroke1  int
	LastStroke2  int
	FirstStroke1 int
	FirstStroke2 int
	ZongLucky    bool
}

type oldWuXing struct {
	ID      string
	First   string
	Second  string
	Third   string
	Fortune string
}



func (e *Exporter) Export() error {
	log.Println("Loading external reference data...")
	if err := e.loadUnihanData(); err != nil {
		log.Printf("Warning: failed to load Unihan data: %v", err)
	} else {
		log.Printf("  → %d pinyin entries, %d strokes, %d definitions",
			len(e.pinyinMap), len(e.totalStrokes), len(e.definitions))
	}

	db, err := sql.Open("sqlite3", e.dbPath+"?mode=ro")
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	log.Println("Building ID→Char lookup from n_character...")
	idToChar, err := e.buildIDLookup(db)
	if err != nil {
		return err
	}
	log.Printf("  → %d ID mappings", len(idToChar))

	var allChars []SeedCharacter

	log.Println("Exporting n_character table (primary)...")
	nChars, err := e.queryNCharacters(db)
	if err != nil {
		return err
	}
	nSeeds := e.transformNCharacters(nChars, idToChar)
	log.Printf("  → %d from n_character", len(nSeeds))
	allChars = append(allChars, nSeeds...)

	log.Println("Exporting character table (supplement)...")
	oldChars, err := e.queryCharacters(db)
	if err != nil {
		return err
	}
	cSeeds := e.transformCharacters(oldChars)
	log.Printf("  → %d from character", len(cSeeds))

	allChars = e.mergeCharacters(allChars, cSeeds)
	log.Printf("  → %d after merge (dedup by char)", len(allChars))

	log.Println("Exporting wu_ge_lucky table...")
	oldWuGes, err := e.queryWuGeLucky(db)
	if err != nil {
		return err
	}
	seedWuGes := transformWuGeLucky(oldWuGes)

	log.Println("Exporting wu_xing table...")
	oldWuXings, err := e.queryWuXing(db)
	if err != nil {
		return err
	}
	seedWuXings := transformWuXing(oldWuXings)

	charPath := filepath.Join(e.seedDir, "character.json")
	if err := writeJSON(charPath, allChars); err != nil {
		return err
	}
	log.Printf("  → %d characters → %s", len(allChars), charPath)

	wugePath := filepath.Join(e.seedDir, "wu_ge_lucky.json")
	if err := writeJSON(wugePath, seedWuGes); err != nil {
		return err
	}
	log.Printf("  → %d wu_ge_lucky → %s", len(seedWuGes), wugePath)

	wuxingPath := filepath.Join(e.seedDir, "wu_xing.json")
	if err := writeJSON(wuxingPath, seedWuXings); err != nil {
		return err
	}
	log.Printf("  → %d wu_xing → %s", len(seedWuXings), wuxingPath)

	if len(e.changes) > 0 {
		changelogPath := filepath.Join(e.seedDir, "changelog.json")
		if err := writeJSON(changelogPath, e.changes); err != nil {
			return err
		}
		log.Printf("  → %d field changes → %s", len(e.changes), changelogPath)
	}

	log.Println("Export completed!")
	return nil
}

func (e *Exporter) buildIDLookup(db *sql.DB) (map[int]string, error) {
	rows, err := db.Query("SELECT id, char FROM n_character")
	if err != nil {
		return nil, fmt.Errorf("query n_character for lookup: %w", err)
	}
	defer rows.Close()

	lookup := make(map[int]string)
	for rows.Next() {
		var id int
		var ch string
		if err := rows.Scan(&id, &ch); err != nil {
			return nil, fmt.Errorf("scan lookup row: %w", err)
		}
		lookup[id] = ch
	}
	return lookup, nil
}

func (e *Exporter) queryNCharacters(db *sql.DB) ([]oldNCharacter, error) {
	rows, err := db.Query("SELECT id, pin_yin, char, char_stroke, radical, radical_stroke, is_regular, is_simplified, simplified_id, is_traditional, traditional_id, is_kang_xi, kang_xi_id, kang_xi_stroke, is_variant, variant_id, is_science, science_stroke, wu_xing, lucky, explanation, comment, need_fix FROM n_character")
	if err != nil {
		return nil, fmt.Errorf("query n_character: %w", err)
	}
	defer rows.Close()

	var results []oldNCharacter
	for rows.Next() {
		var c oldNCharacter
		err := rows.Scan(
			&c.ID, &c.PinYin, &c.Char, &c.CharStroke, &c.Radical, &c.RadicalStroke,
			&c.IsRegular, &c.IsSimplified, &c.SimplifiedID,
			&c.IsTraditional, &c.TraditionalID,
			&c.IsKangXi, &c.KangXiID, &c.KangXiStroke,
			&c.IsVariant, &c.VariantID,
			&c.IsScience, &c.ScienceStroke,
			&c.WuXing, &c.Lucky, &c.Explanation, &c.Comment, &c.NeedFix,
		)
		if err != nil {
			return nil, fmt.Errorf("scan n_character row: %w", err)
		}
		results = append(results, c)
	}
	log.Printf("  Read %d rows from n_character", len(results))
	return results, nil
}

func (e *Exporter) queryCharacters(db *sql.DB) ([]oldCharacter, error) {
	rows, err := db.Query("SELECT hash, pin_yin, ch, radical, radical_stroke, stroke, is_kang_xi, kang_xi, kang_xi_stroke, simple_radical, simple_radical_stroke, simple_total_stroke, traditional_radical, traditional_radical_stroke, traditional_total_stroke, name_science, wu_xing, lucky, regular, traditional_character, variant_character, comment, science_stroke FROM character")
	if err != nil {
		return nil, fmt.Errorf("query character: %w", err)
	}
	defer rows.Close()

	var results []oldCharacter
	for rows.Next() {
		var c oldCharacter
		err := rows.Scan(
			&c.Hash, &c.PinYin, &c.Ch, &c.Radical, &c.RadicalStroke,
			&c.Stroke, &c.IsKangXi, &c.KangXi, &c.KangXiStroke,
			&c.SimpleRadical, &c.SimpleRadicalStroke, &c.SimpleTotalStroke,
			&c.TraditionalRadical, &c.TraditionalRadicalStroke, &c.TraditionalTotalStroke,
			&c.NameScience, &c.WuXing, &c.Lucky, &c.Regular,
			&c.TraditionalCharacter, &c.VariantCharacter, &c.Comment, &c.ScienceStroke,
		)
		if err != nil {
			return nil, fmt.Errorf("scan character row: %w", err)
		}
		results = append(results, c)
	}
	log.Printf("  Read %d rows from character", len(results))
	return results, nil
}

func (e *Exporter) queryWuGeLucky(db *sql.DB) ([]oldWuGeLucky, error) {
	rows, err := db.Query("SELECT last_stroke_1, last_stroke_2, first_stroke_1, first_stroke_2, tian_ge, tian_da_yan, ren_ge, ren_da_yan, di_ge, di_da_yan, wai_ge, wai_da_yan, zong_ge, zong_da_yan, zong_lucky, zong_sex, zong_max FROM wu_ge_lucky")
	if err != nil {
		return nil, fmt.Errorf("query wu_ge_lucky: %w", err)
	}
	defer rows.Close()

	var results []oldWuGeLucky
	for rows.Next() {
		var w oldWuGeLucky
		var discardTianGe, discardRenGe, discardDiGe, discardWaiGe, discardZongGe int
		var discardTianDaYan, discardRenDaYan, discardDiDaYan, discardWaiDaYan, discardZongDaYan string
		var discardZongSex, discardZongMax bool
		err := rows.Scan(
			&w.LastStroke1, &w.LastStroke2, &w.FirstStroke1, &w.FirstStroke2,
			&discardTianGe, &discardTianDaYan, &discardRenGe, &discardRenDaYan,
			&discardDiGe, &discardDiDaYan, &discardWaiGe, &discardWaiDaYan,
			&discardZongGe, &discardZongDaYan, &w.ZongLucky, &discardZongSex, &discardZongMax,
		)
		if err != nil {
			return nil, fmt.Errorf("scan wu_ge_lucky row: %w", err)
		}
		results = append(results, w)
	}
	log.Printf("  Read %d rows from wu_ge_lucky", len(results))
	return results, nil
}

func (e *Exporter) queryWuXing(db *sql.DB) ([]oldWuXing, error) {
	rows, err := db.Query("SELECT id, first, second, third, fortune FROM wu_xing")
	if err != nil {
		return nil, fmt.Errorf("query wu_xing: %w", err)
	}
	defer rows.Close()

	var results []oldWuXing
	for rows.Next() {
		var w oldWuXing
		err := rows.Scan(&w.ID, &w.First, &w.Second, &w.Third, &w.Fortune)
		if err != nil {
			return nil, fmt.Errorf("scan wu_xing row: %w", err)
		}
		results = append(results, w)
	}
	log.Printf("  Read %d rows from wu_xing", len(results))
	return results, nil
}

func (e *Exporter) mergeCharacters(primary []SeedCharacter, supplement []SeedCharacter) []SeedCharacter {
	existing := make(map[string]*SeedCharacter, len(primary))
	for i := range primary {
		existing[primary[i].Char] = &primary[i]
	}

	added := 0
	enriched := 0
	for _, sup := range supplement {
		if ex, ok := existing[sup.Char]; ok {
			if e.enrichCharacter(ex, sup) {
				enriched++
			}
		} else {
			primary = append(primary, sup)
			existing[sup.Char] = &primary[len(primary)-1]
			added++
		}
	}

	log.Printf("  Merge: %d added, %d enriched from supplement", added, enriched)
	return primary
}

func (e *Exporter) enrichCharacter(dst *SeedCharacter, src SeedCharacter) bool {
	changed := false

	if dst.WuXing == "" && src.WuXing != "" {
		e.recordChange(dst.Char, "wu_xing", "", src.WuXing, "enrich_from_character", "character")
		dst.WuXing = src.WuXing
		changed = true
	}
	if len(dst.Pinyin) == 0 && len(src.Pinyin) > 0 {
		e.recordChange(dst.Char, "pinyin", "", fmt.Sprintf("%v", src.Pinyin), "enrich_from_character", "character")
		dst.Pinyin = src.Pinyin
		changed = true
	}
	if dst.SimplifiedStroke == 0 && src.SimplifiedStroke > 0 {
		e.recordChange(dst.Char, "simplified_stroke", "0", fmt.Sprintf("%d", src.SimplifiedStroke), "enrich_from_character", "character")
		dst.SimplifiedStroke = src.SimplifiedStroke
		changed = true
	}
	if dst.TraditionalStroke == 0 && src.TraditionalStroke > 0 {
		e.recordChange(dst.Char, "traditional_stroke", "0", fmt.Sprintf("%d", src.TraditionalStroke), "enrich_from_character", "character")
		dst.TraditionalStroke = src.TraditionalStroke
		changed = true
	}
	if dst.KangxiStroke == 0 && src.KangxiStroke > 0 {
		e.recordChange(dst.Char, "kangxi_stroke", "0", fmt.Sprintf("%d", src.KangxiStroke), "enrich_from_character", "character")
		dst.KangxiStroke = src.KangxiStroke
		changed = true
	}
	if dst.ScienceStroke == 0 && src.ScienceStroke > 0 {
		e.recordChange(dst.Char, "science_stroke", "0", fmt.Sprintf("%d", src.ScienceStroke), "enrich_from_character", "character")
		dst.ScienceStroke = src.ScienceStroke
		changed = true
	}
	if dst.Radical == "" && src.Radical != "" {
		e.recordChange(dst.Char, "radical", "", src.Radical, "enrich_from_character", "character")
		dst.Radical = src.Radical
		dst.RadicalStroke = src.RadicalStroke
		changed = true
	}
	if dst.Meaning == "" && src.Meaning != "" {
		dst.Meaning = src.Meaning
		changed = true
	}
	if dst.Comment == "" && src.Comment != "" {
		dst.Comment = src.Comment
		changed = true
	}
	if !dst.IsSimplified && src.IsSimplified {
		e.recordChange(dst.Char, "is_simplified", "false", "true", "enrich_from_character", "character")
		dst.IsSimplified = src.IsSimplified
		if dst.SimplifiedOfChar == "" {
			dst.SimplifiedOfChar = src.SimplifiedOfChar
		}
		changed = true
	}
	if !dst.IsTraditional && src.IsTraditional {
		e.recordChange(dst.Char, "is_traditional", "false", "true", "enrich_from_character", "character")
		dst.IsTraditional = src.IsTraditional
		changed = true
	}
	if !dst.IsVariant && src.IsVariant {
		e.recordChange(dst.Char, "is_variant", "false", "true", "enrich_from_character", "character")
		dst.IsVariant = src.IsVariant
		if dst.VariantOfChar == "" {
			dst.VariantOfChar = src.VariantOfChar
		}
		changed = true
	}

	return changed
}

func writeJSON(filename string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create %s: %w", filename, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func readJSON(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open %s: %w", filename, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(v); err != nil {
		return fmt.Errorf("decode %s: %w", filename, err)
	}
	return nil
}

func parseJSONString(s string) []string {
	if s == "" || s == "null" {
		return nil
	}
	var result []string
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return []string{s}
	}
	return result
}

func parseJSONInts(s string) []int {
	if s == "" || s == "null" || s == "[]" {
		return nil
	}
	var result []int
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return nil
	}
	return result
}

func resolveIDs(ids []int, lookup map[int]string) []string {
	var result []string
	for _, id := range ids {
		if ch, ok := lookup[id]; ok {
			result = append(result, ch)
		} else {
			result = append(result, fmt.Sprintf("id:%d", id))
		}
	}
	return result
}

func (e *Exporter) loadUnihanData() error {
	unihanDir := filepath.Join(e.rawDataDir, "unihan")

	if err := e.loadUnihanReadings(filepath.Join(unihanDir, "Unihan_Readings.txt")); err != nil {
		log.Printf("Warning: failed to load Unihan_Readings: %v", err)
	}

	if err := e.loadUnihanIRG(filepath.Join(unihanDir, "Unihan_IRGSources.txt")); err != nil {
		log.Printf("Warning: failed to load Unihan_IRGSources: %v", err)
	}

	return nil
}

func (e *Exporter) loadUnihanReadings(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		codePoint := parts[0]
		field := parts[1]
		value := strings.Join(parts[2:], " ")

		if !strings.HasPrefix(codePoint, "U+") {
			continue
		}

		r, err := strconv.ParseUint(codePoint[2:], 16, 32)
		if err != nil {
			continue
		}
		char := string(rune(r))

		switch field {
		case "kMandarin":
			pinyins := strings.Split(value, " ")
			var cleaned []string
			for _, p := range pinyins {
				p = strings.TrimSpace(p)
				if p != "" {
					cleaned = append(cleaned, p)
				}
			}
			if len(cleaned) > 0 && e.pinyinMap[char] == nil {
				e.pinyinMap[char] = cleaned
			}
		case "kHanyuPinyin":
			for _, part := range strings.Split(value, " ") {
				if idx := strings.Index(part, ":"); idx != -1 {
					pinyinStr := part[idx+1:]
					pinyins := strings.Split(pinyinStr, ",")
					var cleaned []string
					for _, p := range pinyins {
						p = strings.TrimSpace(p)
						if p != "" {
							cleaned = append(cleaned, p)
						}
					}
					if len(cleaned) > 0 && e.pinyinMap[char] == nil {
						e.pinyinMap[char] = cleaned
						break
					}
				}
			}
		case "kDefinition":
			if e.definitions[char] == "" {
				e.definitions[char] = value
			}
		}
	}

	return scanner.Err()
}

func (e *Exporter) loadUnihanIRG(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		codePoint := parts[0]
		field := parts[1]
		value := parts[2]

		if !strings.HasPrefix(codePoint, "U+") {
			continue
		}

		r, err := strconv.ParseUint(codePoint[2:], 16, 32)
		if err != nil {
			continue
		}
		char := string(rune(r))

		if field == "kTotalStrokes" {
			if strokes, err := strconv.Atoi(value); err == nil && e.totalStrokes[char] == 0 {
				e.totalStrokes[char] = strokes
			}
		}
	}

	return scanner.Err()
}
