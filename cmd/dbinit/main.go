package main

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"os"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/internal/database"
	"github.com/babyname/fate/internal/seeddb"
)

func main() {
	cfg := config.DefaultConfig()
	cfg.Database.Mode = "file"
	cfg.Database.Name = "fate"

	b := database.New(cfg.Database)
	client, err := b.Client()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	count, err := client.Character.Query().Count(ctx)
	if err != nil {
		log.Fatalf("Failed to count characters: %v", err)
	}

	if count > 0 {
		log.Printf("Database already has %d characters, skipping import", count)
	} else {
		seeds, err := seeddb.LoadEmbeddedCharacters()
		if err != nil {
			log.Fatalf("Failed to load embedded seeds: %v", err)
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
	}

	finalCount, _ := client.Character.Query().Count(ctx)
	log.Printf("Database has %d characters total", finalCount)

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

	testChars := []string{"西", "门", "王", "李", "张", "诸", "葛"}
	for _, ch := range testChars {
		c, err := client.Character.Query().Where(character.CharEQ(ch)).First(ctx)
		if err != nil {
			log.Printf("  %q: NOT FOUND!", ch)
		} else {
			log.Printf("  %q: stroke=%d wuxing=%s", ch, c.ScienceStroke, c.WuXing)
		}
	}

	client.Close()
	log.Println("Compressing fate.db -> fate.db.gz...")
	if err := compressFile("fate", "internal/database/data/fate.db.gz"); err != nil {
		log.Fatalf("Failed to compress: %v", err)
	}
	log.Println("Done! fate.db.gz is ready for embedding")
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
