//go:build integration

package fate

import (
	"context"
	"testing"
	"time"

	"github.com/babyname/fate/v4/config"
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
