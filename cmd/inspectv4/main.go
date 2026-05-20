package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/sqlite3ent/sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "model/db/data.db?mode=ro")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tables, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer tables.Close()

	fmt.Println("=== Tables ===")
	for tables.Next() {
		var name string
		tables.Scan(&name)
		fmt.Println(" ", name)
	}

	for _, table := range []string{"characters", "n_character"} {
		fmt.Printf("\n=== %s columns ===\n", table)
		cols, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}
		for cols.Next() {
			var cid int
			var name, typ string
			var notnull int
			var dflt sql.NullString
			var pk int
			cols.Scan(&cid, &name, &typ, &notnull, &dflt, &pk)
			fmt.Printf("  %d: %s (%s) notnull=%d pk=%d\n", cid, name, typ, notnull, pk)
		}
		cols.Close()

		var count int
		db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		fmt.Printf("  Total rows: %d\n", count)
	}

	fmt.Println("\n=== n_character sample (first 3) ===")
	rows, err := db.Query("SELECT id, char, pin_yin, wu_xing, kang_xi_stroke, science_stroke, is_simplified, is_traditional, is_kang_xi, explanation FROM n_character LIMIT 3")
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		for rows.Next() {
			var id int
			var char, pinyin, wuxing, explanation string
			var kangxiStroke, scienceStroke int
			var isSimplified, isTraditional, isKangxi bool
			rows.Scan(&id, &char, &pinyin, &wuxing, &kangxiStroke, &scienceStroke, &isSimplified, &isTraditional, &isKangxi, &explanation)
			fmt.Printf("  id=%d char=%s pinyin=%s wuxing=%s kangxi=%d science=%d simp=%v trad=%v kangxi=%v expl=%s\n",
				id, char, pinyin, wuxing, kangxiStroke, scienceStroke, isSimplified, isTraditional, isKangxi, explanation)
		}
		rows.Close()
	}

	fmt.Println("\n=== n_character wu_xing coverage ===")
	var total, withWx int
	db.QueryRow("SELECT COUNT(*) FROM n_character").Scan(&total)
	db.QueryRow("SELECT COUNT(*) FROM n_character WHERE wu_xing IS NOT NULL AND wu_xing != ''").Scan(&withWx)
	fmt.Printf("  Total: %d, With wu_xing: %d (%.1f%%)\n", total, withWx, float64(withWx)/float64(total)*100)

	fmt.Println("\n=== n_character pinyin coverage ===")
	var withPy int
	db.QueryRow("SELECT COUNT(*) FROM n_character WHERE pin_yin IS NOT NULL AND pin_yin != ''").Scan(&withPy)
	fmt.Printf("  Total: %d, With pinyin: %d (%.1f%%)\n", total, withPy, float64(withPy)/float64(total)*100)

	fmt.Println("\n=== n_character kangxi_stroke coverage ===")
	var withKx int
	db.QueryRow("SELECT COUNT(*) FROM n_character WHERE kang_xi_stroke > 0").Scan(&withKx)
	fmt.Printf("  Total: %d, With kangxi_stroke: %d (%.1f%%)\n", total, withKx, float64(withKx)/float64(total)*100)

	fmt.Println("\n=== characters wu_xing coverage ===")
	var cTotal, cWithWx int
	db.QueryRow("SELECT COUNT(*) FROM characters").Scan(&cTotal)
	db.QueryRow("SELECT COUNT(*) FROM characters WHERE wu_xing IS NOT NULL AND wu_xing != ''").Scan(&cWithWx)
	fmt.Printf("  Total: %d, With wu_xing: %d (%.1f%%)\n", cTotal, cWithWx, float64(cWithWx)/float64(cTotal)*100)
}
