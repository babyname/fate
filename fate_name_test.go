package fate

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/babyname/fate/v4/config"
	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/ent/character"
	_ "github.com/sqlite3ent/sqlite3"
)

const testDBPath = "fate_test.db"

func openTestDB(t *testing.T) *ent.Client {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?cache=shared&_journal=WAL&_fk=0", testDBPath)
	client, err := ent.Open("sqlite3", dsn)
	if err != nil {
		t.Fatalf("Open database error = %v", err)
	}
	return client
}

func TestDatabaseHasCharacters(t *testing.T) {
	client := openTestDB(t)
	defer client.Close()

	ctx := context.Background()
	count, err := client.Character.Query().Count(ctx)
	if err != nil {
		t.Fatalf("Count characters error = %v", err)
	}
	t.Logf("Total characters in database: %d", count)
	if count == 0 {
		t.Fatal("Database should have characters")
	}

	regularCount, err := client.Character.Query().Where(character.RegularEQ(true)).Count(ctx)
	if err != nil {
		t.Fatalf("Count regular characters error = %v", err)
	}
	t.Logf("Regular characters: %d", regularCount)

	nameableCount, err := client.Character.Query().Where(character.NameableEQ(true)).Count(ctx)
	if err != nil {
		t.Fatalf("Count nameable characters error = %v", err)
	}
	t.Logf("Nameable characters: %d", nameableCount)

	nameableRegular, err := client.Character.Query().
		Where(character.NameableEQ(true), character.RegularEQ(true)).
		Count(ctx)
	if err != nil {
		t.Fatalf("Count nameable regular characters error = %v", err)
	}
	t.Logf("Nameable + Regular characters: %d", nameableRegular)
}

func TestDatabaseCharacterDataQuality(t *testing.T) {
	client := openTestDB(t)
	defer client.Close()

	ctx := context.Background()

	chars, err := client.Character.Query().
		Where(character.RegularEQ(true), character.NameableEQ(true)).
		Limit(100).
		All(ctx)
	if err != nil {
		t.Fatalf("Query characters error = %v", err)
	}

	noPinyin := 0
	noWuXing := 0
	noStroke := 0
	for _, c := range chars {
		if len(c.Pinyin) == 0 {
			noPinyin++
		}
		if c.WuXing == "" {
			noWuXing++
		}
		if c.ScienceStroke == 0 {
			noStroke++
		}
	}
	t.Logf("Sample 100 regular+nameable chars: no_pinyin=%d, no_wuxing=%d, no_stroke=%d",
		noPinyin, noWuXing, noStroke)

	if noPinyin > 0 {
		t.Errorf("Found %d regular+nameable characters without pinyin (should be 0)", noPinyin)
	}
}

func TestQueryLastName(t *testing.T) {
	client := openTestDB(t)
	defer client.Close()

	ctx := context.Background()

	commonLastNames := []string{"张", "李", "王", "刘", "陈", "杨", "赵", "黄", "周", "吴"}
	for _, name := range commonLastNames {
		ch, err := client.Character.Query().Where(character.CharEQ(name)).First(ctx)
		if err != nil {
			t.Errorf("Query last name %q error = %v", name, err)
			continue
		}
		if ch.WuXing == "" {
			t.Errorf("Last name %q has no WuXing", name)
		}
		if len(ch.Pinyin) == 0 {
			t.Errorf("Last name %q has no pinyin", name)
		}
		if ch.ScienceStroke == 0 {
			t.Errorf("Last name %q has no science_stroke", name)
		}
		t.Logf("  %s: pinyin=%v, wuxing=%s, stroke=%d", name, ch.Pinyin, ch.WuXing, ch.ScienceStroke)
	}
}

func TestNameGenerationE2E(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Database = config.DBConfig{
		Driver: "sqlite3",
		Name:   testDBPath,
	}

	f, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	s := f.NewSessionWithFilter(NewFilter(FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: CharacterFilterTypeDefault,
		MinStroke:           3,
		MaxStroke:           18,
		RegularFilter:       true,
		DaYanFilter:         true,
		WuXingFilter:        true,
		SexFilter:           false,
	}))

	born, err := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")
	if err != nil {
		t.Fatalf("Parse time error = %v", err)
	}

	input := &Input{
		Last: [2]string{"张", ""},
		Born: born,
		Sex:  Sex(1),
	}

	err = s.Start(input)
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	s.Wait()

	output := input.Output()
	total := output.Total()
	t.Logf("Total names generated: %d", total)

	if total == 0 {
		t.Fatal("Expected at least some names to be generated, got 0")
	}

	count := 0
	for _, sn := range output.AllNames() {
		count++
		if count <= 10 {
			t.Logf("  Name: %s%s (Score: %.1f)", sn.Name[0].Char, sn.Name[1].Char, sn.Score)
		}
	}
	t.Logf("Iterated %d names", count)
}
