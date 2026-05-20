package fate

import (
	"context"
	"fmt"
	"log"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/internal/database"
	filterpkg "github.com/babyname/fate/internal/filter"
	"github.com/babyname/fate/internal/naming"
	"github.com/babyname/fate/internal/repository"
	"github.com/babyname/fate/internal/seeddb"
	"github.com/babyname/fate/internal/session"
)

const minCharacterCount = 10000

type Fate interface {
	NewSession() Session
	NewSessionWithFilter(f Filter) Session
}

type fateImpl struct {
	cfg *config.Config
	db  *repository.Repository
}

func (f *fateImpl) NewSessionWithFilter(flt Filter) Session {
	return session.NewSession(f.db, flt)
}

func (f *fateImpl) NewSession() Session {
	return f.NewSessionWithFilter(filterpkg.DefaultFilter())
}

func New(cfg *config.Config) (Fate, error) {
	b := database.New(cfg.Database)
	client, err := b.Client()
	if err != nil {
		return nil, err
	}
	repo := repository.New(client)
	if err := ensureDataSeeded(repo); err != nil {
		log.Printf("Warning: auto-seed failed: %v", err)
	}
	return &fateImpl{
		cfg: cfg,
		db:  repo,
	}, nil
}

func ensureDataSeeded(repo *repository.Repository) error {
	ctx := context.Background()
	count, err := repo.Character.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("check character count: %w", err)
	}

	if count >= minCharacterCount {
		log.Printf("[DATA-CHECK] Database OK: %d characters", count)
		return ensureMeaningUpdated(repo)
	}

	log.Printf("[DATA-WARN] Database has only %d characters (threshold: %d)", count, minCharacterCount)

	seeds, err := seeddb.LoadEmbeddedCharacters()
	if err != nil {
		log.Printf("[DATA-ERROR] Failed to load embedded seed data: %v", err)
		log.Println("[DATA-FALLBACK] Falling back to builtin seeds (~340 chars)")
		return seedBuiltinCharacters(ctx, repo)
	}

	log.Printf("[DATA-REPAIR] Importing %d characters from embedded seed data...", len(seeds))
	imported, err := importSeedCharacters(ctx, repo, seeds)
	if err != nil {
		log.Printf("[DATA-ERROR] Embedded seed import failed: %v", err)
		log.Println("[DATA-FALLBACK] Falling back to builtin seeds (~340 chars)")
		return seedBuiltinCharacters(ctx, repo)
	}

	newCount, _ := repo.Character.Query().Count(ctx)
	log.Printf("[DATA-REPAIR] Imported %d characters, database now has %d total", imported, newCount)

	if newCount < minCharacterCount {
		log.Printf("[DATA-WARN] After import only %d characters, still below threshold %d", newCount, minCharacterCount)
		return seedBuiltinCharacters(ctx, repo)
	}
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

func seedBuiltinCharacters(ctx context.Context, repo *repository.Repository) error {
	log.Println("[DATA-FALLBACK] Seeding with built-in character data (limited coverage, ~340 chars)...")
	seeds := seeddb.BuiltinCharacterSeeds()
	created := 0
	for i := range seeds {
		s := &seeds[i]
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
		if _, err := builder.Save(ctx); err != nil {
			log.Printf("  Warning: failed to seed char %q: %v", s.Char, err)
			continue
		}
		created++
	}
	log.Printf("[DATA-FALLBACK] Seeded %d characters into database", created)

	wxSeeds := seeddb.BuiltinWuXingSeeds()
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
	log.Printf("[DATA-FALLBACK] Seeded %d wu_xing records into database", wxCreated)
	return nil
}

func importSeedCharacters(ctx context.Context, repo *repository.Repository, seeds []seeddb.SeedCharacter) (int, error) {
	total := len(seeds)
	imported := 0

	existingChars, err := repo.Character.Query().All(ctx)
	if err != nil {
		return 0, fmt.Errorf("query existing characters: %w", err)
	}
	existingMap := make(map[string]*ent.Character, len(existingChars))
	for _, c := range existingChars {
		existingMap[c.Char] = c
	}
	log.Printf("  Existing characters in DB: %d, seeds to import: %d", len(existingMap), total)

	for i, s := range seeds {
		if existing, ok := existingMap[s.Char]; ok {
			needsUpdate := false
			builder := repo.Character.UpdateOneID(existing.ID)

			if (existing.WuXing == "" || !isValidWuXing(existing.WuXing)) && s.WuXing != "" && isValidWuXing(s.WuXing) {
				builder.SetWuXing(s.WuXing)
				needsUpdate = true
			}
			if !existing.IsSimplified && s.IsSimplified {
				builder.SetIsSimplified(s.IsSimplified)
				needsUpdate = true
			}
			if !existing.IsTraditional && s.IsTraditional {
				builder.SetIsTraditional(s.IsTraditional)
				needsUpdate = true
			}
			if existing.KangxiStroke == 0 && s.KangxiStroke > 0 {
				builder.SetKangxiStroke(s.KangxiStroke)
				needsUpdate = true
			}
			if existing.SimplifiedStroke == 0 && s.SimplifiedStroke > 0 {
				builder.SetSimplifiedStroke(s.SimplifiedStroke)
				needsUpdate = true
			}
			if existing.TraditionalStroke == 0 && s.TraditionalStroke > 0 {
				builder.SetTraditionalStroke(s.TraditionalStroke)
				needsUpdate = true
			}
			if existing.ScienceStroke == 0 && s.ScienceStroke > 0 {
				builder.SetScienceStroke(s.ScienceStroke)
				needsUpdate = true
			}
			if len(existing.Pinyin) == 0 && len(s.Pinyin) > 0 {
				builder.SetPinyin(s.Pinyin)
				needsUpdate = true
			}
			if existing.Radical == "" && s.Radical != "" {
				builder.SetRadical(s.Radical)
				needsUpdate = true
			}
			if !existing.Nameable && s.Nameable {
				builder.SetNameable(s.Nameable)
				needsUpdate = true
			}
			if !existing.Regular && s.Regular {
				builder.SetRegular(s.Regular)
				needsUpdate = true
			}
			if existing.CommonLevel == 0 && s.CommonLevel > 0 {
				builder.SetCommonLevel(s.CommonLevel)
				needsUpdate = true
			}
			if existing.Meaning == "" && s.Meaning != "" {
				builder.SetMeaning(s.Meaning)
				needsUpdate = true
			}

			if needsUpdate {
				if err := builder.Exec(ctx); err != nil {
					log.Printf("  Warning: failed to update char %q: %v", s.Char, err)
				}
			}
			imported++
		} else {
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
		}

		if (i+1)%5000 == 0 {
			log.Printf("  Characters: %d/%d", i+1, total)
		}
	}

	return imported, nil
}

func ensureMeaningUpdated(repo *repository.Repository) error {
	ctx := context.Background()
	chars, err := repo.Character.Query().
		Where(
			character.Or(
				character.MeaningEQ(""),
				character.MeaningIsNil(),
			),
		).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query chars without meaning: %w", err)
	}
	if len(chars) == 0 {
		return nil
	}

	seedMap := make(map[string]seeddb.SeedCharacter)
	for _, s := range seeddb.BuiltinCharacterSeeds() {
		seedMap[s.Char] = s
	}

	updated := 0
	for _, c := range chars {
		s, ok := seedMap[c.Char]
		if !ok || s.Meaning == "" {
			continue
		}
		err := repo.Character.UpdateOneID(c.ID).
			SetMeaning(s.Meaning).
			Exec(ctx)
		if err != nil {
			log.Printf("  Warning: failed to update meaning for %q: %v", c.Char, err)
			continue
		}
		updated++
	}
	if updated > 0 {
		log.Printf("Updated meaning for %d characters", updated)
	}
	return nil
}

type Session = session.Session
type Input = session.Input
type Output = session.Output
type SessionState = session.SessionState

type Filter = filterpkg.Filter
type FilterOption = filterpkg.FilterOption
type CharacterFilterType = filterpkg.CharacterFilterType

type Sex = naming.Sex
type Name = naming.Name
type NameBasic = naming.NameBasic
type FirstName = naming.FirstName

type ScoredName = session.ScoredName

const (
	SexBoy  Sex = naming.SexBoy
	SexGirl Sex = naming.SexGirl
)

var (
	NewFilter     = filterpkg.NewFilter
	DefaultFilter = filterpkg.DefaultFilter
)

const (
	SessionStateWaiting    SessionState = session.SessionStateWaiting
	SessionStateGenerating SessionState = session.SessionStateGenerating
	SessionStateFinish     SessionState = session.SessionStateFinish
	SessionStateCanceled   SessionState = session.SessionStateCanceled
	SessionStateFailed     SessionState = session.SessionStateFailed
)

const (
	CharacterFilterTypeDefault CharacterFilterType = filterpkg.CharacterFilterTypeDefault
	CharacterFilterTypeChs     CharacterFilterType = filterpkg.CharacterFilterTypeChs
	CharacterFilterTypeCht     CharacterFilterType = filterpkg.CharacterFilterTypeCht
	CharacterFilterTypeKangxi  CharacterFilterType = filterpkg.CharacterFilterTypeKangxi
)

var _ Fate = (*fateImpl)(nil)
