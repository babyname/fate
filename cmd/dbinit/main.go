package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/internal/database"
	"github.com/babyname/fate/internal/seeddb"
)

func main() {
	force := false
	for _, arg := range os.Args[1:] {
		if arg == "--force" || arg == "-f" {
			force = true
		}
	}

	dbName := "fate"
	outputPath := "resources/fate.db.gz"

	if force {
		for _, p := range []string{dbName, dbName + "-shm", dbName + "-wal"} {
			os.Remove(p)
		}
		log.Printf("Force mode: removed existing database files")
	}

	cfg := config.DefaultConfig()
	cfg.Database.Mode = "file"
	cfg.Database.Name = dbName

	b := database.New(cfg.Database)
	client, err := b.Client()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	ctx := context.Background()
	count, err := client.Character.Query().Count(ctx)
	if err != nil {
		log.Fatalf("Failed to count characters: %v", err)
	}

	if count > 0 && !force {
		log.Printf("Database already has %d characters, skipping import", count)
	} else {
		if count > 0 && force {
			log.Printf("Force mode: clearing %d existing characters", count)
			_, err := client.Character.Delete().Exec(ctx)
			if err != nil {
				log.Fatalf("Failed to clear characters: %v", err)
			}
		}
		if err := importCharacters(ctx, client); err != nil {
			log.Fatalf("Failed to import characters: %v", err)
		}
		if err := importWuXing(ctx, client); err != nil {
			log.Fatalf("Failed to import wuxing: %v", err)
		}
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

func verifyCharacters(ctx context.Context, client *ent.Client) {
	testChars := []string{"西", "门", "王", "李", "张", "诸", "葛"}
	for _, ch := range testChars {
		c, err := client.Character.Query().Where(character.CharEQ(ch)).First(ctx)
		if err != nil {
			log.Printf("  %q: NOT FOUND!", ch)
		} else {
			log.Printf("  %q: stroke=%d wuxing=%s radical=%s meaning=%s", ch, c.ScienceStroke, c.WuXing, c.Radical, truncate(c.Meaning, 30))
		}
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
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
