package seeddb

import (
	"context"
	"fmt"
	"log"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/internal/repository"
)

const (
	MinCharacterCount = 10000
)

func EnsureSeeded(ctx context.Context, repo *repository.Repository) error {
	count, err := repo.Character.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("check character count: %w", err)
	}

	if count >= MinCharacterCount {
		log.Printf("[SEED] Database OK: %d characters", count)
		return nil
	}

	log.Printf("[SEED] Database has only %d characters (threshold: %d)", count, MinCharacterCount)

	seeds, err := LoadEmbeddedCharacters()
	if err != nil {
		log.Printf("[SEED] Failed to load embedded seed data: %v", err)
		log.Println("[SEED] Falling back to builtin seeds (~340 chars)")
		return seedBuiltinCharacters(ctx, repo)
	}

	log.Printf("[SEED] Importing %d characters from embedded seed data...", len(seeds))
	imported, err := importSeedCharacters(ctx, repo, seeds)
	if err != nil {
		log.Printf("[SEED] Embedded seed import failed: %v", err)
		log.Println("[SEED] Falling back to builtin seeds (~340 chars)")
		return seedBuiltinCharacters(ctx, repo)
	}

	newCount, _ := repo.Character.Query().Count(ctx)
	log.Printf("[SEED] Imported %d characters, database now has %d total", imported, newCount)

	if newCount < MinCharacterCount {
		log.Printf("[SEED] After import only %d characters, still below threshold %d", newCount, MinCharacterCount)
		return seedBuiltinCharacters(ctx, repo)
	}
	return nil
}

func importSeedCharacters(ctx context.Context, repo *repository.Repository, seeds []SeedCharacter) (int, error) {
	existingChars, err := repo.Character.Query().
		Select(character.FieldChar).
		All(ctx)
	if err != nil {
		return 0, fmt.Errorf("query existing characters: %w", err)
	}
	existingSet := make(map[string]struct{}, len(existingChars))
	for _, c := range existingChars {
		existingSet[c.Char] = struct{}{}
	}

	var newSeeds []SeedCharacter
	for _, s := range seeds {
		if _, ok := existingSet[s.Char]; !ok {
			newSeeds = append(newSeeds, s)
		}
	}

	if len(newSeeds) == 0 {
		log.Printf("[SEED] All %d seed characters already exist, nothing to import", len(seeds))
		return 0, nil
	}

	log.Printf("[SEED] %d existing, %d new to import out of %d seeds", len(existingSet), len(newSeeds), len(seeds))

	imported := 0
	for i := 0; i < len(newSeeds); i += batchSize {
		end := i + batchSize
		if end > len(newSeeds) {
			end = len(newSeeds)
		}
		batch := newSeeds[i:end]

		builders := make([]*ent.CharacterCreate, 0, len(batch))
		for _, s := range batch {
			builder := repo.Character.Create().
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
			if s.WuXing != "" && isValidWuXing(s.WuXing) {
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
			builders = append(builders, builder)
		}

		created, err := repo.Character.CreateBulk(builders...).Save(ctx)
		if err != nil {
			log.Printf("[SEED] Batch %d-%d failed: %v, falling back to individual creates", i, end, err)
			for _, s := range batch {
				if err := createSingleChar(ctx, repo, &s); err != nil {
					log.Printf("[SEED] Warning: failed to create char %q: %v", s.Char, err)
				}
			}
			imported += len(batch)
		} else {
			imported += len(created)
		}

		log.Printf("[SEED] Progress: %d/%d", imported, len(newSeeds))
	}

	return imported, nil
}

func createSingleChar(ctx context.Context, repo *repository.Repository, s *SeedCharacter) error {
	builder := repo.Character.Create().
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
	if s.WuXing != "" && isValidWuXing(s.WuXing) {
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
	_, err := builder.Save(ctx)
	return err
}

func seedBuiltinCharacters(ctx context.Context, repo *repository.Repository) error {
	log.Println("[SEED] Seeding with built-in character data (limited coverage, ~340 chars)...")
	seeds := BuiltinCharacterSeeds()
	created := 0
	for i := range seeds {
		if err := createSingleChar(ctx, repo, &seeds[i]); err != nil {
			log.Printf("[SEED] Warning: failed to seed char %q: %v", seeds[i].Char, err)
			continue
		}
		created++
	}
	log.Printf("[SEED] Seeded %d characters into database", created)

	wxSeeds := BuiltinWuXingSeeds()
	wxCreated := 0
	for i := range wxSeeds {
		w := &wxSeeds[i]
		builder := repo.WuXing.Create().SetID(w.ID)
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
		wxCreated++
	}
	log.Printf("[SEED] Seeded %d wu_xing records into database", wxCreated)
	return nil
}

func isValidWuXing(wx string) bool {
	switch wx {
	case "金", "木", "水", "火", "土":
		return true
	default:
		return false
	}
}
