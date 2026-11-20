//go:build integration

package fate

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/babyname/fate/v4/config"
	"github.com/babyname/fate/v4/ranking"
)

// TestGenerateNilFilterCompletes proves at runtime that omitting Filter no
// longer expands the unbounded candidate cross-product. It needs the local
// SQLite file "fate" (see examples/README.md), so it is behind the
// integration build tag and is run explicitly by the acceptance scripts.
func TestGenerateNilFilterCompletes(t *testing.T) {
	engine, err := New(config.DefaultConfig())
	if err != nil {
		t.Fatalf("New() error = %v (run 'go run ./cmd/dbinit' first)", err)
	}

	born, err := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	start := time.Now()
	result, err := engine.Generate(ctx, GenerateRequest{
		LastName: "张",
		Born:     born,
		Sex:      SexBoy,
		// Filter deliberately omitted.
	})
	if err != nil {
		t.Fatalf("Generate with nil Filter failed after %s: %v", time.Since(start), err)
	}
	if result == nil || len(result.TopNames) == 0 {
		t.Fatalf("Generate with nil Filter returned no names after %s", time.Since(start))
	}
	if elapsed := time.Since(start); elapsed > 60*time.Second {
		t.Fatalf("Generate with nil Filter took %s, want < 60s", elapsed)
	}
	t.Logf("nil Filter Generate completed in %s with %d name results", time.Since(start), len(result.TopNames))
}

func TestGenerateStreamRecommendationWeight(t *testing.T) {
	engine, err := New(config.DefaultConfig())
	if err != nil {
		t.Fatalf("New() error = %v (run 'go run ./cmd/dbinit' first)", err)
	}
	weight, err := ranking.NewWeight(23)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	var mu sync.Mutex
	previewBatches, finalBatches := 0, 0
	err = engine.GenerateStream(ctx, GenerateRequest{
		LastName:             "李",
		Born:                 time.Date(2024, time.June, 15, 10, 30, 0, 0, time.UTC),
		Sex:                  SexGirl,
		RecommendationWeight: &weight,
	}, func(batch StreamBatch) {
		mu.Lock()
		defer mu.Unlock()

		if batch.IsFinal {
			finalBatches++
		} else {
			previewBatches++
		}
		if len(batch.TopResults) == 0 {
			t.Errorf("stream batch has no candidates (final=%t)", batch.IsFinal)
		}
		for index, name := range batch.TopResults {
			if name.Ranking == nil {
				t.Errorf("%s: missing ranking (final=%t)", name.FullName, batch.IsFinal)
				continue
			}
			want := weight.Combine(name.Ranking.Traditional, name.Ranking.SoundAndMeaning)
			if *name.Ranking != want || name.Score != want.Total || name.Rank != index+1 {
				t.Errorf("%s: got ranking=%+v score=%g rank=%d, want ranking=%+v score=%g rank=%d (final=%t)",
					name.FullName, *name.Ranking, name.Score, name.Rank, want, want.Total, index+1, batch.IsFinal)
			}
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if previewBatches == 0 || finalBatches != 1 {
		t.Fatalf("got %d preview batches and %d final batches, want at least one preview and exactly one final batch", previewBatches, finalBatches)
	}
	t.Logf("verified recommendation weight 23 in %d preview batches and the final batch", previewBatches)
}
