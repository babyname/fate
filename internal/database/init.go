package database

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/babyname/fate/v4/config"
	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/ent/schema"
	"github.com/babyname/fate/v4/resources"
)

type seedChar struct {
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

func ensureDBFile(cfg config.DBConfig) (string, error) {
	name := cfg.Name
	if name == "" {
		name = "fate"
	}

	dbFile := cfg.DBFile
	if dbFile == "" {
		dbFile = name + ".db"
	}

	if _, err := os.Stat(dbFile); err == nil {
		return dbFile, nil
	}

	gzFile := dbFile + ".gz"
	if _, err := os.Stat(gzFile); err == nil {
		if err := decompressGZ(gzFile, dbFile); err != nil {
			return "", fmt.Errorf("decompress %s: %w", gzFile, err)
		}
		return dbFile, nil
	}

	return "", nil
}

func decompressGZ(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer gr.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, gr)
	return err
}

func needsInit(client *ent.Client) (bool, error) {
	ctx := context.Background()
	first, err := client.Version.Query().First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return false, fmt.Errorf("query version: %w", err)
	}
	if first == nil {
		return true, nil
	}
	return first.CurrentVersion != schema.CurrentDataVersion, nil
}

func initializeFromJSON(ctx context.Context, client *ent.Client) error {
	var seeds []seedChar
	if err := json.Unmarshal(resources.CharacterJSON, &seeds); err != nil {
		return fmt.Errorf("parse character.json: %w", err)
	}

	charMap := make(map[string]*ent.Character, len(seeds))
	for _, s := range seeds {
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
			continue
		}
		charMap[s.Char] = created
	}

	for _, s := range seeds {
		thisChar, ok := charMap[s.Char]
		if !ok {
			continue
		}
		updater := thisChar.Update()
		updated := false

		if s.SimplifiedOfChar != "" {
			if target, ok := charMap[s.SimplifiedOfChar]; ok {
				updater.AddSimplifiedOf(target)
				updated = true
			}
		}
		if s.VariantOfChar != "" {
			if target, ok := charMap[s.VariantOfChar]; ok {
				updater.SetVariantOf(target)
				updated = true
			}
		}

		if updated {
			_, _ = updater.Save(ctx)
		}
	}

	_, err := client.Version.Create().
		SetCurrentVersion(schema.CurrentDataVersion).
		SetUpdatedUnix(int(time.Now().Unix())).
		Save(ctx)
	return err
}
