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

var criticalChars = []string{"西", "门", "东", "南", "北", "上", "诸葛", "司马", "欧阳"}

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
		log.Printf("[DATA-CHECK] Database has %d characters (>= %d threshold), OK", count, minCharacterCount)
		if err := verifyCriticalChars(ctx, repo); err != nil {
			log.Printf("[DATA-WARN] Critical character verification failed: %v", err)
			log.Println("[DATA-REPAIR] Re-importing from embedded seed data to fix missing characters...")
			seeds, loadErr := seeddb.LoadEmbeddedCharacters()
			if loadErr != nil {
				log.Printf("[DATA-ERROR] Failed to load embedded seed data for repair: %v", loadErr)
			} else {
				imported, importErr := importSeedCharacters(ctx, repo, seeds)
				if importErr != nil {
					log.Printf("[DATA-ERROR] Repair import failed: %v", importErr)
				} else {
					log.Printf("[DATA-REPAIR] Imported %d characters for repair", imported)
				}
			}
		}
		return ensureMeaningUpdated(repo)
	}

	if count > 0 {
		log.Printf("[DATA-WARN] Database has only %d characters (threshold: %d), data is incomplete!", count, minCharacterCount)
	}

	log.Println("[DATA-REPAIR] Loading embedded seed data (compiled into binary, always available)...")
	seeds, err := seeddb.LoadEmbeddedCharacters()
	if err != nil {
		log.Printf("[DATA-ERROR] Failed to load embedded seed data: %v", err)
		log.Println("[DATA-FALLBACK] Falling back to builtin seeds (limited: ~340 chars)")
		return seedBuiltinCharacters(ctx, repo)
	}

	log.Printf("[DATA-REPAIR] Importing %d characters from embedded seed data...", len(seeds))
	imported, err := importSeedCharacters(ctx, repo, seeds)
	if err != nil {
		log.Printf("[DATA-ERROR] Embedded seed import failed: %v", err)
		log.Println("[DATA-FALLBACK] Falling back to builtin seeds (limited: ~340 chars)")
		return seedBuiltinCharacters(ctx, repo)
	}

	newCount, _ := repo.Character.Query().Count(ctx)
	log.Printf("[DATA-REPAIR] Imported %d characters, database now has %d total", imported, newCount)

	if newCount >= minCharacterCount {
		log.Printf("[DATA-CHECK] Database integrity verified: %d characters >= %d threshold", newCount, minCharacterCount)
		return nil
	}

	log.Printf("[DATA-WARN] After import only %d characters, still below threshold %d", newCount, minCharacterCount)
	return seedBuiltinCharacters(ctx, repo)
}

func verifyCriticalChars(ctx context.Context, repo *repository.Repository) error {
	var missing []string
	for _, ch := range criticalChars {
		exists, err := repo.Character.Query().
			Where(character.CharEQ(ch)).
			Exist(ctx)
		if err != nil {
			return fmt.Errorf("verify char %q: %w", ch, err)
		}
		if !exists {
			missing = append(missing, ch)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing %d critical characters: %v (database data is corrupted)", len(missing), missing)
	}
	log.Printf("[DATA-CHECK] All %d critical characters verified present", len(criticalChars))
	return nil
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
	batchSize := 500

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batch := seeds[i:end]

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
			builders = append(builders, builder)
		}

		created, err := repo.Character.CreateBulk(builders...).Save(ctx)
		if err != nil {
			log.Printf("  Warning: batch %d-%d import error: %v", i, end, err)
			continue
		}
		imported += len(created)
		log.Printf("  Characters: %d/%d", imported, total)
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
