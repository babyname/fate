package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/babyname/fate/internal/seeddb"
	_ "github.com/sqlite3ent/sqlite3"
)

func parseJSONString(s string) []string {
	if s == "" || s == "null" {
		return nil
	}
	var result []string
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		return []string{s}
	}
	return result
}

func main() {
	db, err := sql.Open("sqlite3", "model/db/data.db?mode=ro")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, pin_yin, ch, science_stroke, radical, radical_stroke, stroke, is_kangxi, kangxi, kangxi_stroke, simple_radical, simple_radical_stroke, simple_total_stroke, traditional_radical, traditional_radical_stroke, traditional_total_stroke, is_name_science, wu_xing, lucky, is_regular, traditional_character, variant_character, comment FROM characters`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var seeds []seeddb.SeedCharacter
	for rows.Next() {
		var c struct {
			ID                      string
			PinYin                  string
			Ch                      string
			ScienceStroke           int
			Radical                 string
			RadicalStroke           int
			Stroke                  int
			IsKangxi                bool
			Kangxi                  string
			KangxiStroke            int
			SimpleRadical           string
			SimpleRadicalStroke     int
			SimpleTotalStroke       int
			TraditionalRadical      string
			TraditionalRadicalStroke int
			TraditionalTotalStroke  int
			IsNameScience           bool
			WuXing                  string
			Lucky                   string
			IsRegular               bool
			TraditionalCharacter    string
			VariantCharacter        string
			Comment                 string
		}
		err := rows.Scan(&c.ID, &c.PinYin, &c.Ch, &c.ScienceStroke, &c.Radical, &c.RadicalStroke, &c.Stroke, &c.IsKangxi, &c.Kangxi, &c.KangxiStroke, &c.SimpleRadical, &c.SimpleRadicalStroke, &c.SimpleTotalStroke, &c.TraditionalRadical, &c.TraditionalRadicalStroke, &c.TraditionalTotalStroke, &c.IsNameScience, &c.WuXing, &c.Lucky, &c.IsRegular, &c.TraditionalCharacter, &c.VariantCharacter, &c.Comment)
		if err != nil {
			log.Printf("scan error: %v", err)
			continue
		}

		wx := c.WuXing
		if wx == "岁" {
			wx = ""
		}

		nameable := c.IsNameScience || (c.ScienceStroke > 0 && c.ScienceStroke <= 30)

		sc := seeddb.SeedCharacter{
			Char:              c.Ch,
			IsSimplified:      c.SimpleTotalStroke > 0,
			IsKangxi:          c.IsKangxi,
			Regular:           c.IsRegular,
			Nameable:          nameable,
			WuXing:            wx,
			Radical:           c.Radical,
			RadicalStroke:     c.RadicalStroke,
			SimplifiedStroke:  c.SimpleTotalStroke,
			TraditionalStroke: c.TraditionalTotalStroke,
			KangxiStroke:      c.KangxiStroke,
			ScienceStroke:     c.ScienceStroke,
			Source:            "v4_characters",
		}

		sc.Pinyin = parseJSONString(c.PinYin)

		if c.TraditionalTotalStroke > 0 && c.SimpleTotalStroke != c.TraditionalTotalStroke {
			tradChars := parseJSONString(c.TraditionalCharacter)
			if len(tradChars) > 0 {
				sc.IsSimplified = true
				sc.SimplifiedOfChar = tradChars[0]
			} else {
				sc.IsTraditional = true
			}
		}

		if c.Kangxi != "" && c.Kangxi != c.Ch {
			sc.IsVariant = true
			sc.VariantOfChar = c.Kangxi
		}

		variantChars := parseJSONString(c.VariantCharacter)
		if len(variantChars) > 0 {
			sc.IsVariant = true
			if sc.VariantOfChar == "" {
				sc.VariantOfChar = variantChars[0]
			}
		}

		if c.Lucky != "" {
			sc.Comment = fmt.Sprintf("lucky=%s", c.Lucky)
		}
		commentParts := parseJSONString(c.Comment)
		if len(commentParts) > 0 {
			if sc.Comment != "" {
				sc.Comment += "; "
			}
			if len(commentParts[0]) > 200 {
				sc.Comment += commentParts[0][:200] + "..."
			} else {
				sc.Comment += commentParts[0]
			}
		}

		if sc.Nameable && len(sc.Pinyin) == 0 {
			sc.Nameable = false
		}

		seeds = append(seeds, sc)
	}

	log.Printf("Converted %d characters from v4 database", len(seeds))

	if err := os.MkdirAll("data/seed", 0755); err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("data/seed/character.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(seeds); err != nil {
		log.Fatal(err)
	}

	log.Printf("Written %d characters to data/seed/character.json", len(seeds))

	for _, ch := range []string{"西", "门", "东", "南", "李", "王"} {
		for _, s := range seeds {
			if s.Char == ch {
				fmt.Printf("  %s: pinyin=%v, wuxing=%s, science=%d, kangxi=%d, nameable=%v\n",
					ch, s.Pinyin, s.WuXing, s.ScienceStroke, s.KangxiStroke, s.Nameable)
				break
			}
		}
	}
}
