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

type oldWuGeLucky struct {
	LastStroke1  int
	LastStroke2  int
	FirstStroke1 int
	FirstStroke2 int
	ZongLucky    bool
}

type oldWuXing struct {
	ID      string
	SanCai  string
	Lucky   bool
	Comment string
}

func (e *Exporter) Export() error {
	db, err := sql.Open("sqlite3", e.dbPath+"?mode=ro")
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	log.Println("Exporting character table...")
	oldChars, err := e.queryCharacters(db)
	if err != nil {
		return err
	}

	log.Println("Exporting wu_ge_lucky table...")
	oldWuGes, err := e.queryWuGeLucky(db)
	if err != nil {
		return err
	}

	log.Println("Exporting wu_xing table...")
	oldWuXings, err := e.queryWuXing(db)
	if err != nil {
		return err
	}

	log.Println("Transforming data to seed format...")
	seedChars := transformCharacters(oldChars)
	seedWuGes := transformWuGeLucky(oldWuGes)
	seedWuXings := transformWuXing(oldWuXings)

	charPath := filepath.Join(e.seedDir, "character.json")
	if err := writeJSON(charPath, seedChars); err != nil {
		return err
	}
	log.Printf("  → %d characters → %s", len(seedChars), charPath)

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
		var first, second, third, fortune string
		err := rows.Scan(
			&w.ID, &first, &second, &third, &fortune,
		)
		if err != nil {
			return nil, fmt.Errorf("scan wu_xing row: %w", err)
		}
		w.SanCai = first + second + third
		w.Comment = fortune
		w.Lucky = fortune == "大吉" || fortune == "吉"
		results = append(results, w)
	}
	log.Printf("  Read %d rows from wu_xing", len(results))
	return results, nil
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
