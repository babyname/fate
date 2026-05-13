package seeddb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	db, err := sql.Open("sqlite3", e.dbPath+"?mode=ro")
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	var allChars []SeedCharacter

	log.Println("Exporting n_character table (primary)...")
	nChars, err := e.queryNCharacters(db)
	if err != nil {
		return err
	}
	nSeeds := transformNCharacters(nChars)
	log.Printf("  → %d from n_character", len(nSeeds))
	allChars = append(allChars, nSeeds...)

	log.Println("Exporting character table (supplement)...")
	oldChars, err := e.queryCharacters(db)
	if err != nil {
		return err
	}
	cSeeds := transformCharacters(oldChars)
	log.Printf("  → %d from character", len(cSeeds))

	allChars = mergeCharacters(allChars, cSeeds)
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

	log.Println("Export completed!")
	return nil
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

func mergeCharacters(primary []SeedCharacter, supplement []SeedCharacter) []SeedCharacter {
	existing := make(map[string]*SeedCharacter, len(primary))
	for i := range primary {
		existing[primary[i].Char] = &primary[i]
	}

	added := 0
	enriched := 0
	for _, sup := range supplement {
		if ex, ok := existing[sup.Char]; ok {
			if enrichCharacter(ex, sup) {
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

func enrichCharacter(dst *SeedCharacter, src SeedCharacter) bool {
	changed := false

	if dst.WuXing == "" && src.WuXing != "" {
		dst.WuXing = src.WuXing
		changed = true
	}
	if len(dst.Pinyin) == 0 && len(src.Pinyin) > 0 {
		dst.Pinyin = src.Pinyin
		changed = true
	}
	if dst.SimplifiedStroke == 0 && src.SimplifiedStroke > 0 {
		dst.SimplifiedStroke = src.SimplifiedStroke
		changed = true
	}
	if dst.TraditionalStroke == 0 && src.TraditionalStroke > 0 {
		dst.TraditionalStroke = src.TraditionalStroke
		changed = true
	}
	if dst.KangxiStroke == 0 && src.KangxiStroke > 0 {
		dst.KangxiStroke = src.KangxiStroke
		changed = true
	}
	if dst.ScienceStroke == 0 && src.ScienceStroke > 0 {
		dst.ScienceStroke = src.ScienceStroke
		changed = true
	}
	if dst.Radical == "" && src.Radical != "" {
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
		dst.IsSimplified = src.IsSimplified
		if dst.SimplifiedOfChar == "" {
			dst.SimplifiedOfChar = src.SimplifiedOfChar
		}
		changed = true
	}
	if !dst.IsTraditional && src.IsTraditional {
		dst.IsTraditional = src.IsTraditional
		changed = true
	}
	if !dst.IsVariant && src.IsVariant {
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
