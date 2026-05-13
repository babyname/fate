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
	ID                     int    `json:"id"`
	Ch                     string `json:"ch"`
	PinYin                string `json:"pin_yin"`
	WuXing                 string `json:"wu_xing"`
	Stroke                 int    `json:"stroke"`
	KangXiStroke           int    `json:"kang_xi_stroke"`
	ScienceStroke          int    `json:"science_stroke"`
	IsKangXi               bool   `json:"is_kang_xi"`
	SimpleTotalStroke      int    `json:"simple_total_stroke"`
	TraditionalTotalStroke int    `json:"traditional_total_stroke"`
	Regular                bool   `json:"regular"`
	NameScience            bool   `json:"name_science"`
	Lucky                  string `json:"lucky"`
	TraditionalCharacter   string `json:"traditional_character"`
	VariantCharacter       string `json:"variant_character"`
	KangXi                 string `json:"kang_xi"`
}

type OldWuGeLucky struct {
	ID           int  `json:"id"`
	LastStroke1  int  `json:"last_stroke_1"`
	LastStroke2  int  `json:"last_stroke_2"`
	FirstStroke1 int  `json:"first_stroke_1"`
	FirstStroke2 int  `json:"first_stroke_2"`
	TianGeLucky  bool `json:"tian_ge_lucky"`
	RenGeLucky   bool `json:"ren_ge_lucky"`
	DiGeLucky    bool `json:"di_ge_lucky"`
	WaiGeLucky   bool `json:"wai_ge_lucky"`
	ZongGeLucky  bool `json:"zong_ge_lucky"`
	ZongLucky    bool `json:"zong_lucky"`
}

type OldWuXing struct {
	ID      string `json:"id"`
	TianGe  int    `json:"tian_ge"`
	RenGe   int    `json:"ren_ge"`
	DiGe    int    `json:"di_ge"`
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
	oldChars, err := e.queryTable(db, "character", e.scanCharacters)
	if err != nil {
		return err
	}

	log.Println("Exporting wu_ge_lucky table...")
	oldWuGes, err := e.queryTable(db, "wu_ge_lucky", e.scanWuGeLucky)
	if err != nil {
		return err
	}

	log.Println("Exporting wu_xing table...")
	oldWuXings, err := e.queryTable(db, "wu_xing", e.scanWuXing)
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

func (e *Exporter) scanCharacters(rows *sql.Rows) (interface{}, error) {
	var c OldCharacter
	err := rows.Scan(
		&c.ID, &c.Ch, &c.PinYin, &c.WuXing, &c.Stroke, &c.KangXiStroke,
		&c.ScienceStroke, &c.IsKangXi, &c.SimpleTotalStroke,
		&c.TraditionalTotalStroke, &c.Regular, &c.NameScience,
		&c.Lucky, &c.TraditionalCharacter, &c.VariantCharacter, &c.KangXi,
	)
	return c, err
}

func (e *Exporter) scanWuGeLucky(rows *sql.Rows) (interface{}, error) {
	var w OldWuGeLucky
	err := rows.Scan(
		&w.ID, &w.LastStroke1, &w.LastStroke2, &w.FirstStroke1, &w.FirstStroke2,
		&w.TianGeLucky, &w.RenGeLucky, &w.DiGeLucky, &w.WaiGeLucky,
		&w.ZongGeLucky, &w.ZongLucky,
	)
	return w, err
}

func (e *Exporter) scanWuXing(rows *sql.Rows) (interface{}, error) {
	var w OldWuXing
	err := rows.Scan(
		&w.ID, &w.TianGe, &w.RenGe, &w.DiGe, &w.SanCai, &w.Lucky, &w.Comment,
	)
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
