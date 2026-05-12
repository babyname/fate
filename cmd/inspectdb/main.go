package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/sqlite3ent/sqlite3"
)

func main() {
	dbPath := "fate"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	db, err := sql.Open("sqlite3", dbPath+"?mode=ro")
	if err != nil {
		fmt.Println("open error:", err)
		return
	}
	defer db.Close()

	tables, err := getTables(db)
	if err != nil {
		fmt.Println("get tables error:", err)
		return
	}

	fmt.Println("=== 数据库表列表 ===")
	for _, t := range tables {
		fmt.Printf("\n--- 表: %s ---\n", t)

		schema, err := getTableSchema(db, t)
		if err != nil {
			fmt.Println("  schema error:", err)
			continue
		}
		fmt.Println(schema)

		count, err := getRowCount(db, t)
		if err != nil {
			fmt.Println("  count error:", err)
			continue
		}
		fmt.Printf("  行数: %d\n", count)

		if count > 0 && count < 500000 {
			sample, err := getSampleRows(db, t, 3)
			if err != nil {
				fmt.Println("  sample error:", err)
				continue
			}
			fmt.Println("  样本数据:")
			for _, row := range sample {
				fmt.Printf("    %s\n", row)
			}
		}
	}

	// Character table statistics
	fmt.Println("\n\n=== Character 表统计 ===")
	printCharStats(db)
}

func printCharStats(db *sql.DB) {
	stats := []struct {
		name  string
		query string
	}{
		{"总字数", "SELECT COUNT(*) FROM character"},
		{"有拼音", "SELECT COUNT(*) FROM character WHERE pin_yin != '' AND pin_yin != 'null'"},
		{"有五行", "SELECT COUNT(*) FROM character WHERE wu_xing != ''"},
		{"有康熙笔画", "SELECT COUNT(*) FROM character WHERE kang_xi_stroke > 0"},
		{"有简体笔画", "SELECT COUNT(*) FROM character WHERE simple_total_stroke > 0"},
		{"有繁体笔画", "SELECT COUNT(*) FROM character WHERE traditional_total_stroke > 0"},
		{"有姓名学笔画", "SELECT COUNT(*) FROM character WHERE science_stroke > 0"},
		{"常用字(regular)", "SELECT COUNT(*) FROM character WHERE regular = 1"},
		{"姓名学可用", "SELECT COUNT(*) FROM character WHERE name_science = 1"},
		{"有康熙字", "SELECT COUNT(*) FROM character WHERE kang_xi != '' AND kang_xi != 'null'"},
		{"有繁体字", "SELECT COUNT(*) FROM character WHERE traditional_character != '' AND traditional_character != 'null'"},
		{"有异体字", "SELECT COUNT(*) FROM character WHERE variant_character != '' AND variant_character != 'null'"},
		{"五行分布", "SELECT wu_xing, COUNT(*) as cnt FROM character WHERE wu_xing != '' GROUP BY wu_xing ORDER BY cnt DESC"},
		{"吉凶分布", "SELECT lucky, COUNT(*) as cnt FROM character WHERE lucky != '' GROUP BY lucky ORDER BY cnt DESC"},
		{"笔画范围(简体)", "SELECT MIN(simple_total_stroke), MAX(simple_total_stroke), AVG(simple_total_stroke) FROM character WHERE simple_total_stroke > 0"},
		{"笔画范围(康熙)", "SELECT MIN(kang_xi_stroke), MAX(kang_xi_stroke), AVG(kang_xi_stroke) FROM character WHERE kang_xi_stroke > 0"},
	}

	for _, s := range stats {
		rows, err := db.Query(s.query)
		if err != nil {
			fmt.Printf("  %s: error %v\n", s.name, err)
			continue
		}
		if rows.Next() {
			cols, _ := rows.Columns()
			vals := make([]interface{}, len(cols))
			ptrs := make([]interface{}, len(cols))
			for i := range vals {
				ptrs[i] = &vals[i]
			}
			rows.Scan(ptrs...)
			parts := ""
			for i, c := range cols {
				if i > 0 {
					parts += ", "
				}
				parts += fmt.Sprintf("%s=%v", c, vals[i])
			}
			fmt.Printf("  %s: %s\n", s.name, parts)
		}
		rows.Close()
	}
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, nil
}

func getTableSchema(db *sql.DB, table string) (string, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	result := "  字段:\n"
	for rows.Next() {
		var cid int
		var name, typ string
		var notNull int
		var dfltValue interface{}
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notNull, &dfltValue, &pk); err != nil {
			return "", err
		}
		pkStr := ""
		if pk > 0 {
			pkStr = " [PK]"
		}
		result += fmt.Sprintf("    %s %s%s\n", name, typ, pkStr)
	}
	return result, nil
}

func getRowCount(db *sql.DB, table string) (int, error) {
	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
	return count, err
}

func getSampleRows(db *sql.DB, table string, limit int) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT %d", table, limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []string
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := ""
		for i, col := range cols {
			if i > 0 {
				row += ", "
			}
			v := vals[i]
			switch tv := v.(type) {
			case []byte:
				s := string(tv)
				if len(s) > 50 {
					s = s[:50] + "..."
				}
				row += fmt.Sprintf("%s=%s", col, s)
			case nil:
				row += fmt.Sprintf("%s=NULL", col)
			default:
				row += fmt.Sprintf("%s=%v", col, tv)
			}
		}
		results = append(results, row)
	}
	return results, nil
}
