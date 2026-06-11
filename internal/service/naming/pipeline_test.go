package naming

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	v2 "github.com/godcong/chronos/v2"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/internal/chronosfate"
	"github.com/babyname/fate/v4/internal/filter"
	"github.com/babyname/fate/v4/internal/repository"
	_ "github.com/sqlite3ent/sqlite3"
)

func resolveTestDBPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename) // internal/service/naming/
	root := filepath.Join(dir, "..", "..", "..")
	return filepath.Join(root, "fate_test.db")
}

func openTestRepository(t *testing.T) *repository.Repository {
	t.Helper()
	dbPath := resolveTestDBPath()
	dsn := fmt.Sprintf("file:%s?cache=shared&_journal=WAL&_fk=0", dbPath)
	client, err := ent.Open("sqlite3", dsn)
	if err != nil {
		t.Fatalf("Open database error = %v", err)
	}
	t.Cleanup(func() { client.Close() })
	return repository.New(client)
}

func getZhang(t *testing.T, repo *repository.Repository) *ent.Character {
	t.Helper()
	zhang, err := repo.GetCharacter(context.Background(), repository.Char("张"))
	if err != nil {
		t.Fatalf("GetCharacter(张) error = %v", err)
	}
	return zhang
}

func TestNewPipeline(t *testing.T) {
	repo := openTestRepository(t)
	p := NewPipeline(repo)
	if p == nil {
		t.Fatal("NewPipeline returned nil")
	}
}

func TestPipelineGenerateSmoke(t *testing.T) {
	repo := openTestRepository(t)
	p := NewPipeline(repo)
	zhang := getZhang(t, repo)

	born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")

	flt := filter.NewFilter(filter.FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: filter.CharacterFilterTypeDefault,
		MinStroke:           3,
		MaxStroke:           18,
		RegularFilter:       true,
		DaYanFilter:         true,
		WuXingFilter:        true,
	})

	result, err := p.Generate(context.Background(), GenerateRequest{
		LastName: [2]*ent.Character{zhang, nil},
		Born:     born,
		Sex:      1,
		Filter:   flt,
	})

	if err != nil {
		t.Fatalf("Pipeline.Generate() error = %v", err)
	}
	if result == nil || result.ExcellentTable == nil || len(result.TopNames) == 0 {
		t.Fatal("Pipeline.Generate() produced empty result")
	}

	t.Logf("Generated %d names, Top1: %s (%.1f)", result.ExcellentTable.Len(), result.TopNames[0].FullName, result.TopNames[0].Score)
	if result.CharMap == nil {
		t.Fatal("CharMap is nil")
	}
	if result.FateData == nil {
		t.Fatal("FateData is nil")
	}
	if result.LastNameStrokes[0] != 11 {
		t.Errorf("Expected stroke 11 for 张, got %d", result.LastNameStrokes[0])
	}
}

func TestPipelineWithPrecalcFateData(t *testing.T) {
	repo := openTestRepository(t)
	p := NewPipeline(repo)
	zhang := getZhang(t, repo)

	fateData, err := chronosfate.GetFateData(chronosfate.FateInput{
		Calendar: v2.NewSolarCalendar(time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)),
		Gender:   1,
	})
	if err != nil {
		t.Fatalf("GetFateData() error = %v", err)
	}

	flt := filter.NewFilter(filter.FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: filter.CharacterFilterTypeDefault,
		MinStroke:           3,
		MaxStroke:           18,
		RegularFilter:       true,
		DaYanFilter:         true,
		WuXingFilter:        true,
	})

	result, err := p.Generate(context.Background(), GenerateRequest{
		LastName: [2]*ent.Character{zhang, nil},
		Born:     time.Time{},
		Sex:      1,
		Filter:   flt,
		FateData: fateData,
	})

	if err != nil {
		t.Fatalf("Pipeline.Generate() error = %v", err)
	}
	if result.ExcellentTable.Len() == 0 {
		t.Fatal("Expected non-empty results")
	}
	if result.FateData != fateData {
		t.Fatal("FateData should be the same pointer")
	}
	t.Logf("Generated %d names with pre-calculated FateData", result.ExcellentTable.Len())
}

func TestPipelineGetLucky(t *testing.T) {
	repo := openTestRepository(t)
	p := NewPipeline(repo)
	li, err := repo.GetCharacter(context.Background(), repository.Char("李"))
	if err != nil {
		t.Fatalf("GetCharacter(李) error = %v", err)
	}

	born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")
	flt := filter.NewFilter(filter.FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: filter.CharacterFilterTypeDefault,
		MinStroke:           3,
		MaxStroke:           18,
		RegularFilter:       true,
		DaYanFilter:         true,
		WuXingFilter:        true,
	})

	result, err := p.Generate(context.Background(), GenerateRequest{
		LastName: [2]*ent.Character{li, nil},
		Born:     born,
		Sex:      1,
		Filter:   flt,
	})

	if err != nil {
		t.Fatalf("Pipeline.Generate() for 李 error = %v", err)
	}
	if result.ExcellentTable.Len() == 0 {
		t.Fatal("Expected at least some names for 李")
	}
	t.Logf("李 generated %d names, top1: %s (%.1f)", result.ExcellentTable.Len(), result.TopNames[0].FullName, result.TopNames[0].Score)
}

func TestPipelineStability(t *testing.T) {
	// Heavy test: run twice, verify deterministic output.
	// Use the same parameters as TestNameAnalysisOutput to ensure correctness.
	repo := openTestRepository(t)
	p := NewPipeline(repo)
	zhang := getZhang(t, repo)

	born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")
	flt := filter.NewFilter(filter.FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: filter.CharacterFilterTypeDefault,
		MinStroke:           3,
		MaxStroke:           18,
		RegularFilter:       true,
		DaYanFilter:         true,
		WuXingFilter:        true,
	})

	req := GenerateRequest{
		LastName: [2]*ent.Character{zhang, nil},
		Born:     born,
		Sex:      1,
		Filter:   flt,
	}

	result1, err := p.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("First Generate() error = %v", err)
	}

	result2, err := p.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("Second Generate() error = %v", err)
	}

	if len(result1.TopNames) != len(result2.TopNames) {
		t.Fatalf("TopNames count mismatch: %d vs %d", len(result1.TopNames), len(result2.TopNames))
	}
	for i := 0; i < len(result1.TopNames) && i < 10; i++ {
		if result1.TopNames[i].FullName != result2.TopNames[i].FullName {
			t.Errorf("Top %d name mismatch: %s vs %s", i+1, result1.TopNames[i].FullName, result2.TopNames[i].FullName)
			break
		}
		if result1.TopNames[i].Score != result2.TopNames[i].Score {
			t.Errorf("Top %d score mismatch: %.1f vs %.1f", i+1, result1.TopNames[i].Score, result2.TopNames[i].Score)
			break
		}
	}

	t.Logf("Stability: both runs produced %d names, top1: %s (%.1f)", result1.ExcellentTable.Len(), result1.TopNames[0].FullName, result1.TopNames[0].Score)
}

func TestPipelineContextCancellation(t *testing.T) {
	repo := openTestRepository(t)
	p := NewPipeline(repo)
	zhang := getZhang(t, repo)

	born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")
	flt := filter.NewFilter(filter.FilterOption{
		CharacterFilter: true,
		MinStroke:       3,
		MaxStroke:       18,
		RegularFilter:   true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := p.Generate(ctx, GenerateRequest{
		LastName: [2]*ent.Character{zhang, nil},
		Born:     born,
		Sex:      1,
		Filter:   flt,
	})

	// Must not panic.
	if err != nil {
		t.Logf("Cancelled context returned error: %v", err)
	} else if result != nil {
		t.Logf("Cancelled context produced %d names", result.ExcellentTable.Len())
	}
}
