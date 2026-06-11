package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/ent/character"
	"github.com/babyname/fate/v4/ent/schema"
	"github.com/babyname/fate/v4/resources"
	_ "github.com/sqlite3ent/sqlite3"
)

// SeedChar matches the JSON schema in resources/character.json.
type SeedChar struct {
	Char              string   `json:"char"`
	Unicode           string   `json:"unicode,omitempty"`
	IsSimplified      bool     `json:"is_simplified"`
	IsTraditional     bool     `json:"is_traditional"`
	IsKangxi          bool     `json:"is_kangxi"`
	IsVariant         bool     `json:"is_variant"`
	IsAncient         bool     `json:"is_ancient"`
	Pinyin            []string `json:"pinyin,omitempty"`
	Radical           string   `json:"radical,omitempty"`
	RadicalStroke     int      `json:"radical_stroke,omitempty"`
	SimplifiedStroke  int      `json:"simplified_stroke,omitempty"`
	TraditionalStroke int      `json:"traditional_stroke,omitempty"`
	KangxiStroke      int      `json:"kangxi_stroke,omitempty"`
	ScienceStroke     int      `json:"science_stroke,omitempty"`
	WuXing            string   `json:"wu_xing,omitempty"`
	Regular           bool     `json:"regular"`
	CommonLevel       int      `json:"common_level,omitempty"`
	GenderHint        string   `json:"gender_hint,omitempty"`
	Nameable          bool     `json:"nameable"`
	Meaning           string   `json:"meaning,omitempty"`
	Source            string   `json:"source,omitempty"`
	SourceConfidence  float64  `json:"source_confidence,omitempty"`
	Comment           string   `json:"comment,omitempty"`
	SimplifiedOfChar  string   `json:"simplified_of_char,omitempty"`
	VariantOfChar     string   `json:"variant_of_char,omitempty"`
}

func main() {
	dbName := flag.String("db", "fate", "database file name (without extension)")
	outputPath := *dbName + ".db.gz"
	flag.Parse()

	for _, p := range []string{*dbName, *dbName + "-shm", *dbName + "-wal"} {
		os.Remove(p)
	}
	log.Printf("Cleared existing database files")

	client, err := ent.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&_journal=WAL&_fk=1", *dbName))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	_, err = client.Version.Query().First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		log.Fatalf("Failed to query version: %v", err)
	}
	if ent.IsNotFound(err) {
		_, err := client.Version.Create().
			SetCurrentVersion(schema.CurrentDataVersion).
			SetUpdatedUnix(0).
			Save(ctx)
		if err != nil {
			log.Fatalf("Failed to create version: %v", err)
		}
	}

	if err := importCharacters(ctx, client); err != nil {
		log.Fatalf("Failed to import characters: %v", err)
	}

	verifyCharacters(ctx, client)
	client.Close()

	log.Printf("Compressing %s -> %s...", *dbName, outputPath)
	if err := compressFile(*dbName, outputPath); err != nil {
		log.Fatalf("Failed to compress: %v", err)
	}
	log.Println("Done! fate.db.gz is ready for embedding")
}

func loadEmbeddedCharacters() ([]SeedChar, error) {
	var seeds []SeedChar
	if err := json.Unmarshal(resources.CharacterJSON, &seeds); err != nil {
		return nil, fmt.Errorf("parse character.json: %w", err)
	}
	log.Printf("Loaded %d characters from embedded resources", len(seeds))
	return seeds, nil
}

func importCharacters(ctx context.Context, client *ent.Client) error {
	seeds, err := loadEmbeddedCharacters()
	if err != nil {
		return fmt.Errorf("load embedded characters: %w", err)
	}
	log.Printf("Importing %d characters...", len(seeds))

	charMap := make(map[string]*ent.Character, len(seeds))
	imported := 0
	for i, s := range seeds {
		builder := client.Character.Create().
			SetChar(s.Char).
			SetIsSimplified(s.IsSimplified).
			SetIsTraditional(s.IsTraditional).
			SetIsKangxi(s.IsKangxi).
			SetIsVariant(s.IsVariant).
			SetIsAncient(s.IsAncient).
			SetRegular(s.Regular).
			SetNameable(s.Nameable)
		if s.Unicode != "" {
			builder.SetUnicode(s.Unicode)
		}
		if len(s.Pinyin) > 0 {
			builder.SetPinyin(s.Pinyin)
		}
		if s.Radical != "" {
			builder.SetRadical(s.Radical)
		}
		if s.RadicalStroke > 0 {
			builder.SetRadicalStroke(s.RadicalStroke)
		}
		if s.SimplifiedStroke > 0 {
			builder.SetSimplifiedStroke(s.SimplifiedStroke)
		}
		if s.TraditionalStroke > 0 {
			builder.SetTraditionalStroke(s.TraditionalStroke)
		}
		if s.KangxiStroke > 0 {
			builder.SetKangxiStroke(s.KangxiStroke)
		}
		if s.ScienceStroke > 0 {
			builder.SetScienceStroke(s.ScienceStroke)
		}
		if s.WuXing != "" {
			builder.SetWuXing(s.WuXing)
		}
		if s.CommonLevel > 0 {
			builder.SetCommonLevel(s.CommonLevel)
		}
		if s.GenderHint != "" {
			builder.SetGenderHint(s.GenderHint)
		}
		if s.Meaning != "" {
			builder.SetMeaning(s.Meaning)
		}
		if s.Source != "" {
			builder.SetSource(s.Source)
		}
		if s.SourceConfidence > 0 {
			builder.SetSourceConfidence(s.SourceConfidence)
		}
		if s.Comment != "" {
			builder.SetComment(s.Comment)
		}
		created, err := builder.Save(ctx)
		if err != nil {
			log.Printf("  Warning: failed to create char %q: %v", s.Char, err)
			continue
		}
		charMap[s.Char] = created
		imported++
		if (i+1)%5000 == 0 {
			log.Printf("  Progress (insert): %d/%d", i+1, len(seeds))
		}
	}
	log.Printf("Imported %d characters", imported)

	linkedSimp := 0
	linkedVar := 0
	for i, s := range seeds {
		thisChar, ok := charMap[s.Char]
		if !ok {
			continue
		}
		updater := thisChar.Update()
		updated := false

		if s.SimplifiedOfChar != "" {
			if target, ok := charMap[s.SimplifiedOfChar]; ok {
				updater.AddSimplifiedOf(target)
				linkedSimp++
				updated = true
			}
		}
		if s.VariantOfChar != "" {
			if target, ok := charMap[s.VariantOfChar]; ok {
				updater.SetVariantOf(target)
				linkedVar++
				updated = true
			}
		}
		if updated {
			if _, err := updater.Save(ctx); err != nil {
				log.Printf("  Warning: failed to update char %q: %v", s.Char, err)
			}
		}
		if (i+1)%5000 == 0 {
			log.Printf("  Progress (link): %d/%d", i+1, len(seeds))
		}
	}
	log.Printf("Linked %d simplified-of and %d variant-of relationships", linkedSimp, linkedVar)

	finalCount, _ := client.Character.Query().Count(ctx)
	log.Printf("Database has %d characters total", finalCount)
	return nil
}

func verifyCharacters(ctx context.Context, client *ent.Client) {
	testChars := []string{"西", "门", "明", "瑞", "轩"}
	for _, ch := range testChars {
		c, err := client.Character.Query().Where(character.CharEQ(ch)).First(ctx)
		if err != nil {
			log.Printf("  %q: NOT FOUND!", ch)
		} else {
			log.Printf("  %q: stroke=%d wuxing=%s meaning=%s has_poetry=%v", ch, c.ScienceStroke, c.WuXing, truncate(c.Meaning, 30), c.HasPoetry)
		}
	}
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}

func compressFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	gz.Name = "fate.db"
	defer gz.Close()

	info, _ := in.Stat()
	log.Printf("  Source: %s (%d bytes)", src, info.Size())

	if _, err := io.Copy(gz, in); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	outInfo, _ := out.Stat()
	log.Printf("  Output: %s (%d bytes, ratio: %.1f%%)", dst, outInfo.Size(), float64(outInfo.Size())/float64(info.Size())*100)
	return nil
}
