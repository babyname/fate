package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/ent/character"
	"github.com/babyname/fate/v4/ent/schema"
	"github.com/babyname/fate/v4/internal/dict"
	"github.com/babyname/fate/v4/internal/seeddb"
	_ "github.com/sqlite3ent/sqlite3"
)

func main() {
	dbName := "fate"
	outputPath := "fate.db.gz"

	for _, p := range []string{dbName, dbName + "-shm", dbName + "-wal"} {
		os.Remove(p)
	}
	log.Printf("Cleared existing database files")

	client, err := ent.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&_journal=WAL&_fk=1", dbName))
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
	if err := importWuXing(ctx, client); err != nil {
		log.Fatalf("Failed to import wuxing: %v", err)
	}

	poetryDir := "data/chinese-poetry"
	if len(os.Args) > 1 {
		poetryDir = os.Args[1]
	}
	if info, err := os.Stat(poetryDir); err == nil && info.IsDir() {
		if err := importPoetry(ctx, client, poetryDir); err != nil {
			log.Printf("Warning: Failed to import poetry: %v", err)
		}
	} else {
		log.Printf("Poetry directory %s not found, skipping poetry import", poetryDir)
	}

	verifyCharacters(ctx, client)
	client.Close()

	log.Printf("Compressing %s -> %s...", dbName, outputPath)
	if err := compressFile(dbName, outputPath); err != nil {
		log.Fatalf("Failed to compress: %v", err)
	}
	log.Println("Done! fate.db.gz is ready for embedding")
}

func importCharacters(ctx context.Context, client *ent.Client) error {
	seeds, err := seeddb.LoadEmbeddedCharacters()
	if err != nil {
		return fmt.Errorf("load embedded seeds: %w", err)
	}
	log.Printf("Importing %d characters...", len(seeds))
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
		if s.Comment != "" {
			builder.SetComment(s.Comment)
		}
		if _, err := builder.Save(ctx); err != nil {
			log.Printf("  Warning: failed to create char %q: %v", s.Char, err)
			continue
		}
		imported++
		if (i+1)%5000 == 0 {
			log.Printf("  Progress: %d/%d", i+1, len(seeds))
		}
	}
	log.Printf("Imported %d characters", imported)

	finalCount, _ := client.Character.Query().Count(ctx)
	log.Printf("Database has %d characters total", finalCount)
	return nil
}

func importWuXing(ctx context.Context, client *ent.Client) error {
	wxSeeds := seeddb.BuiltinWuXingSeeds()
	for i := range wxSeeds {
		w := &wxSeeds[i]
		builder := client.WuXing.Create().SetID(w.ID)
		if w.First != "" {
			builder.SetFirst(w.First)
		}
		if w.Second != "" {
			builder.SetSecond(w.Second)
		}
		if w.Third != "" {
			builder.SetThird(w.Third)
		}
		if w.Fortune != "" {
			builder.SetFortune(w.Fortune)
		}
		if _, err := builder.Save(ctx); err != nil {
			continue
		}
	}
	log.Printf("Imported wu_xing data")
	return nil
}

func importPoetry(ctx context.Context, client *ent.Client, poetryDir string) error {
	log.Printf("Importing poetry from %s...", poetryDir)

	entries, err := dict.LoadSelectedPoetryFromDir(poetryDir)
	if err != nil {
		return fmt.Errorf("load poetry: %w", err)
	}
	log.Printf("Loaded %d poems", len(entries))

	uniqueChars := make(map[string]bool)

	for i, e := range entries {
		refs := dict.ExtractCharRefs(e.Content)
		for _, ref := range refs {
			if !unicode.Is(unicode.Han, []rune(ref.Char)[0]) {
				continue
			}
			uniqueChars[ref.Char] = true
		}

		if (i+1)%500 == 0 {
			log.Printf("  Poetry progress: %d/%d", i+1, len(entries))
		}
	}

	log.Printf("Found %d unique chars with poetry source from %d poems", len(uniqueChars), len(entries))

	charWithPoetry := 0
	for ch := range uniqueChars {
		n, err := client.Character.Update().
			Where(character.CharEQ(ch)).
			SetHasPoetry(true).
			Save(ctx)
		if err == nil && n > 0 {
			charWithPoetry++
		}
	}
	log.Printf("Marked %d characters with has_poetry=true", charWithPoetry)

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
