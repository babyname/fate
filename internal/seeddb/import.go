package seeddb

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/internal/database"
	"github.com/google/uuid"
	"golang.org/x/net/context"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/sqlite3ent/sqlite3"
)

const batchSize = 500

func (imp *Importer) Import() error {
	ctx := context.Background()

	cfg := config.DBConfig{
		Driver: imp.cfg.Driver,
		DSN:    imp.cfg.DSN,
		Host:   imp.cfg.Host,
		Port:   imp.cfg.Port,
		User:   imp.cfg.User,
		Pwd:    imp.cfg.Pwd,
		Name:   imp.cfg.Name,
	}

	builder := database.New(cfg)
	client, err := builder.Client()
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer client.Close()

	log.Printf("Connected to %s database", imp.cfg.Driver)

	charFile := filepath.Join(imp.seedDir, "character.json")
	if _, err := os.Stat(charFile); err == nil {
		log.Println("Importing characters...")
		if err := imp.importCharacters(ctx, client, charFile); err != nil {
			return fmt.Errorf("import characters: %w", err)
		}
	} else {
		log.Printf("Skipping characters: %s not found", charFile)
	}

	wugeFile := filepath.Join(imp.seedDir, "wu_ge_lucky.json")
	if _, err := os.Stat(wugeFile); err == nil {
		log.Println("Importing wu_ge_lucky...")
		if err := imp.importWuGeLucky(ctx, client, wugeFile); err != nil {
			return fmt.Errorf("import wu_ge_lucky: %w", err)
		}
	} else {
		log.Printf("Skipping wu_ge_lucky: %s not found", wugeFile)
	}

	wuxingFile := filepath.Join(imp.seedDir, "wu_xing.json")
	if _, err := os.Stat(wuxingFile); err == nil {
		log.Println("Importing wu_xing...")
		if err := imp.importWuXing(ctx, client, wuxingFile); err != nil {
			return fmt.Errorf("import wu_xing: %w", err)
		}
	} else {
		log.Printf("Skipping wu_xing: %s not found", wuxingFile)
	}

	log.Println("Import completed!")
	return nil
}

func (imp *Importer) importCharacters(ctx context.Context, client *ent.Client, filename string) error {
	var seeds []SeedCharacter
	if err := readJSON(filename, &seeds); err != nil {
		return err
	}

	charIDMap := make(map[string]int)
	total := len(seeds)
	imported := 0

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batch := seeds[i:end]

		builders := make([]*ent.CharacterCreate, 0, len(batch))
		for _, sc := range batch {
			builder := client.Character.Create().
				SetChar(sc.Char).
				SetIsSimplified(sc.IsSimplified).
				SetIsTraditional(sc.IsTraditional).
				SetIsKangxi(sc.IsKangxi).
				SetIsVariant(sc.IsVariant).
				SetIsAncient(sc.IsAncient).
				SetRegular(sc.Regular).
				SetNameable(sc.Nameable)

			if sc.Unicode != "" {
				builder.SetUnicode(sc.Unicode)
			}
			if len(sc.Pinyin) > 0 {
				builder.SetPinyin(sc.Pinyin)
			}
			if sc.Radical != "" {
				builder.SetRadical(sc.Radical)
			}
			if sc.RadicalStroke > 0 {
				builder.SetRadicalStroke(sc.RadicalStroke)
			}
			if sc.SimplifiedStroke > 0 {
				builder.SetSimplifiedStroke(sc.SimplifiedStroke)
			}
			if sc.TraditionalStroke > 0 {
				builder.SetTraditionalStroke(sc.TraditionalStroke)
			}
			if sc.KangxiStroke > 0 {
				builder.SetKangxiStroke(sc.KangxiStroke)
			}
			if sc.ScienceStroke > 0 {
				builder.SetScienceStroke(sc.ScienceStroke)
			}
			if sc.WuXing != "" {
				builder.SetWuXing(sc.WuXing)
			}
			if sc.CommonLevel > 0 {
				builder.SetCommonLevel(sc.CommonLevel)
			}
			if sc.GenderHint != "" {
				builder.SetGenderHint(sc.GenderHint)
			}
			if sc.Meaning != "" {
				builder.SetMeaning(sc.Meaning)
			}
			if sc.Source != "" {
				builder.SetSource(sc.Source)
			}
			if sc.SourceConfidence > 0 {
				builder.SetSourceConfidence(sc.SourceConfidence)
			}
			if sc.Comment != "" {
				builder.SetComment(sc.Comment)
			}

			builders = append(builders, builder)
		}

		created, err := client.Character.CreateBulk(builders...).Save(ctx)
		if err != nil {
			return fmt.Errorf("batch %d-%d: %w", i, end, err)
		}

		for idx, c := range created {
			charIDMap[batch[idx].Char] = c.ID
		}

		imported += len(created)
		log.Printf("  Characters: %d/%d", imported, total)
	}

	if err := imp.linkCharacterEdges(ctx, client, seeds, charIDMap); err != nil {
		log.Printf("  Warning: failed to link character edges: %v", err)
	}

	return nil
}

func (imp *Importer) linkCharacterEdges(ctx context.Context, client *ent.Client, seeds []SeedCharacter, charIDMap map[string]int) error {
	traditionalToSimplified := make(map[int]int)
	variantOf := make(map[int]int)

	for _, sc := range seeds {
		if sc.SimplifiedOfChar != "" {
			if simplifiedID, ok := charIDMap[sc.SimplifiedOfChar]; ok {
				currentID := charIDMap[sc.Char]
				traditionalToSimplified[currentID] = simplifiedID
			}
		}
		if sc.VariantOfChar != "" {
			if standardID, ok := charIDMap[sc.VariantOfChar]; ok {
				currentID := charIDMap[sc.Char]
				variantOf[currentID] = standardID
			}
		}
	}

	linked := 0
	for tradID, simpID := range traditionalToSimplified {
		err := client.Character.UpdateOneID(tradID).
			SetTraditionalToSimplifiedID(simpID).
			Exec(ctx)
		if err != nil {
			log.Printf("  Warning: failed to link traditional_to_simplified %d→%d: %v", tradID, simpID, err)
			continue
		}
		linked++
	}

	for variantID, standardID := range variantOf {
		err := client.Character.UpdateOneID(variantID).
			SetVariantOfID(standardID).
			Exec(ctx)
		if err != nil {
			log.Printf("  Warning: failed to link variant_of %d→%d: %v", variantID, standardID, err)
			continue
		}
		linked++
	}

	log.Printf("  Linked %d character edges (trad→simp: %d, variant→std: %d)", linked, len(traditionalToSimplified), len(variantOf))
	return nil
}

func (imp *Importer) importWuGeLucky(ctx context.Context, client *ent.Client, filename string) error {
	var seeds []SeedWuGeLucky
	if err := readJSON(filename, &seeds); err != nil {
		return err
	}

	total := len(seeds)
	imported := 0

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batch := seeds[i:end]

		builders := make([]*ent.WuGeLuckyCreate, 0, len(batch))
		for _, sw := range batch {
			uid := uuid.New()
			builder := client.WuGeLucky.Create().
				SetID(uid).
				SetLastStroke1(sw.LastStroke1).
				SetLastStroke2(sw.LastStroke2).
				SetFirstStroke1(sw.FirstStroke1).
				SetFirstStroke2(sw.FirstStroke2).
				SetTianGe(sw.TianGe).
				SetTianDaYan(sw.TianDaYan).
				SetRenGe(sw.RenGe).
				SetRenDaYan(sw.RenDaYan).
				SetDiGe(sw.DiGe).
				SetDiDaYan(sw.DiDaYan).
				SetWaiGe(sw.WaiGe).
				SetWaiDaYan(sw.WaiDaYan).
				SetZongGe(sw.ZongGe).
				SetZongDaYan(sw.ZongDaYan).
				SetZongLucky(sw.ZongLucky).
				SetZongSex(sw.ZongSex).
				SetZongMax(sw.ZongMax)

			builders = append(builders, builder)
		}

		created, err := client.WuGeLucky.CreateBulk(builders...).Save(ctx)
		if err != nil {
			return fmt.Errorf("batch %d-%d: %w", i, end, err)
		}
		imported += len(created)
		log.Printf("  WuGeLucky: %d/%d", imported, total)
	}

	return nil
}

func (imp *Importer) importWuXing(ctx context.Context, client *ent.Client, filename string) error {
	var seeds []SeedWuXing
	if err := readJSON(filename, &seeds); err != nil {
		return err
	}

	total := len(seeds)
	imported := 0

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batch := seeds[i:end]

		builders := make([]*ent.WuXingCreate, 0, len(batch))
		for _, sw := range batch {
			builder := client.WuXing.Create().
				SetID(sw.ID)

			if sw.First != "" {
				builder.SetFirst(sw.First)
			}
			if sw.Second != "" {
				builder.SetSecond(sw.Second)
			}
			if sw.Third != "" {
				builder.SetThird(sw.Third)
			}
			if sw.Fortune != "" {
				builder.SetFortune(sw.Fortune)
			}

			builders = append(builders, builder)
		}

		created, err := client.WuXing.CreateBulk(builders...).Save(ctx)
		if err != nil {
			return fmt.Errorf("batch %d-%d: %w", i, end, err)
		}
		imported += len(created)
		log.Printf("  WuXing: %d/%d", imported, total)
	}

	return nil
}
