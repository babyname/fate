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

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Tables:")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Printf("  %s\n", name)
	}

	for _, table := range []string{"character", "characters", "n_character", "wu_ge_lucky", "wu_xing"} {
		rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
		if err != nil {
			fmt.Printf("\n%s: not found\n", table)
			continue
		}
		fmt.Printf("\n%s columns:\n", table)
		for rows.Next() {
			var cid int
			var name, typ string
			var notnull int
			var dfltValue interface{}
			var pk int
			rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk)
			fmt.Printf("  %d: %s (%s) notnull=%d pk=%d\n", cid, name, typ, notnull, pk)
		}
		rows.Close()
	}

	var charCount int
	db.QueryRow("SELECT COUNT(*) FROM characters").Scan(&charCount)
	fmt.Printf("\ncharacters count: %d\n", charCount)

	rows2, err := db.Query("SELECT id, ch, pin_yin, wu_xing, science_stroke, kangxi_stroke FROM characters LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nSample data:")
	for rows2.Next() {
		var id string
		var ch, pinyin, wuxing string
		var stroke, kangxiStroke int
		rows2.Scan(&id, &ch, &pinyin, &wuxing, &stroke, &kangxiStroke)
		fmt.Printf("  id=%s %s: pinyin=%s, wuxing=%s, science=%d, kangxi=%d\n", id, ch, pinyin, wuxing, stroke, kangxiStroke)
	}
	rows2.Close()

	rows3, err := db.Query("SELECT id, ch, pin_yin, wu_xing, science_stroke, kangxi_stroke FROM characters WHERE ch IN ('西', '门', '东', '南', '李', '王')")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nSpecific chars:")
	for rows3.Next() {
		var id string
		var ch, pinyin, wuxing string
		var stroke, kangxiStroke int
		rows3.Scan(&id, &ch, &pinyin, &wuxing, &stroke, &kangxiStroke)
		fmt.Printf("  id=%s %s: pinyin=%s, wuxing=%s, science=%d, kangxi=%d\n", id, ch, pinyin, wuxing, stroke, kangxiStroke)
	}
	rows3.Close()
}
