package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/sqlite3ent/sqlite3"
)

type OldCharacter struct {
	Hash                    string `json:"hash"`
	PinYin                  string `json:"pin_yin"`
	Ch                      string `json:"ch"`
	Radical                 string `json:"radical"`
	RadicalStroke           int    `json:"radical_stroke"`
	Stroke                  int    `json:"stroke"`
	IsKangXi                bool   `json:"is_kang_xi"`
	KangXi                  string `json:"kang_xi"`
	KangXiStroke            int    `json:"kang_xi_stroke"`
	SimpleRadical           string `json:"simple_radical"`
	SimpleRadicalStroke     int    `json:"simple_radical_stroke"`
	SimpleTotalStroke       int    `json:"simple_total_stroke"`
	TraditionalRadical      string `json:"traditional_radical"`
	TraditionalRadicalStroke int   `json:"traditional_radical_stroke"`
	TraditionalTotalStroke  int    `json:"traditional_total_stroke"`
	NameScience             bool   `json:"name_science"`
	WuXing                  string `json:"wu_xing"`
	Lucky                   string `json:"lucky"`
	Regular                 bool   `json:"regular"`
	TraditionalCharacter    string `json:"traditional_character"`
	VariantCharacter        string `json:"variant_character"`
	Comment                 string `json:"comment"`
	ScienceStroke           int    `json:"science_stroke"`
}

type OldWuGeLucky struct {
	LastStroke1  int    `json:"last_stroke_1"`
	LastStroke2  int    `json:"last_stroke_2"`
	FirstStroke1 int    `json:"first_stroke_1"`
	FirstStroke2 int    `json:"first_stroke_2"`
	ZongLucky    bool   `json:"zong_lucky"`
}

type OldWuXing struct {
	ID      string `json:"id"`
	SanCai  string `json:"san_cai"`
	Lucky   bool   `json:"lucky"`
	Comment string `json:"comment"`
}

type Exporter struct {
	dbPath  string
	seedDir string
}

func NewExporter(dbPath, seedDir string) *Exporter {
	return &Exporter{dbPath: dbPath, seedDir: seedDir}
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
	seedChars := TransformCharacters(oldChars)
	seedWuGes := TransformWuGeLucky(oldWuGes)
	seedWuXings := TransformWuXing(oldWuXings)

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

func (e *Exporter) queryCharacters(db *sql.DB) ([]OldCharacter, error) {
	rows, err := db.Query("SELECT hash, pin_yin, ch, radical, radical_stroke, stroke, is_kang_xi, kang_xi, kang_xi_stroke, simple_radical, simple_radical_stroke, simple_total_stroke, traditional_radical, traditional_radical_stroke, traditional_total_stroke, name_science, wu_xing, lucky, regular, traditional_character, variant_character, comment, science_stroke FROM character")
	if err != nil {
		return nil, fmt.Errorf("query character: %w", err)
	}
	defer rows.Close()

	var results []OldCharacter
	for rows.Next() {
		var c OldCharacter
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

func (e *Exporter) queryWuGeLucky(db *sql.DB) ([]interface{}, error) {
	return e.queryTable(db, "wu_ge_lucky", e.scanWuGeLucky)
}

func (e *Exporter) queryWuXing(db *sql.DB) ([]interface{}, error) {
	return e.queryTable(db, "wu_xing", e.scanWuXing)
}

type scanFunc func(*sql.Rows) (interface{}, error)

func (e *Exporter) queryTable(db *sql.DB, table string, scan scanFunc) ([]interface{}, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", table, err)
	}
	defer rows.Close()

	var results []interface{}
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scan %s row: %w", table, err)
		}
		results = append(results, item)
	}
	log.Printf("  Read %d rows from %s", len(results), table)
	return results, nil
}

func (e *Exporter) scanWuGeLucky(rows *sql.Rows) (interface{}, error) {
	var w OldWuGeLucky
	var discardID string
	var discardTianGe, discardRenGe, discardDiGe, discardWaiGe, discardZongGe int
	var discardTianDaYan, discardRenDaYan, discardDiDaYan, discardWaiDaYan, discardZongDaYan string
	var discardZongSex, discardZongMax bool
	err := rows.Scan(
		&discardID, &w.LastStroke1, &w.LastStroke2, &w.FirstStroke1, &w.FirstStroke2,
		&discardTianGe, &discardTianDaYan, &discardRenGe, &discardRenDaYan,
		&discardDiGe, &discardDiDaYan, &discardWaiGe, &discardWaiDaYan,
		&discardZongGe, &discardZongDaYan, &w.ZongLucky, &discardZongSex, &discardZongMax,
	)
	_ = discardID
	return w, err
}

func (e *Exporter) scanWuXing(rows *sql.Rows) (interface{}, error) {
	var w OldWuXing
	var discardCreated, discardUpdated, discardDeleted string
	var discardVersion int
	var first, second, third, fortune string
	err := rows.Scan(
		&w.ID, &discardCreated, &discardUpdated, &discardDeleted, &discardVersion,
		&first, &second, &third, &fortune,
	)
	_ = discardCreated
	_ = discardUpdated
	_ = discardDeleted
	_ = discardVersion
	w.SanCai = first + second + third
	w.Comment = fortune
	w.Lucky = fortune == "大吉" || fortune == "吉"
	return w, err
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
