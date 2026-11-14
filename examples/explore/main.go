// Command explore demonstrates the two ways to read the excellent table:
// TopN (a read-only ranking) and Explore (a randomized, unshown-only batch
// that marks entries as shown). This is the backend behind the Web UI's
// /api/v1/names and /api/v1/explore endpoints.
//
//	go run ./examples/explore
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/babyname/fate/v4"
	"github.com/babyname/fate/v4/config"
)

func main() {
	born, err := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")
	if err != nil {
		log.Fatal(err)
	}

	engine, err := fate.New(config.DefaultConfig())
	if err != nil {
		log.Fatalf("init fate (run 'go run ./cmd/dbinit' first): %v", err)
	}

	result, err := engine.Generate(context.Background(), fate.GenerateRequest{
		LastName: "张",
		Born:     born,
		Sex:      fate.SexBoy,
		Filter: &fate.FilterOption{
			CharacterFilter:     true,
			CharacterFilterType: fate.CharacterFilterTypeDefault,
			MinStroke:           3,
			MaxStroke:           18,
			RegularFilter:       true,
			DaYanFilter:         true,
			WuXingFilter:        true,
		},
	})
	if err != nil {
		log.Fatalf("generate: %v", err)
	}

	table := result.ExcellentTable
	fmt.Printf("table: total=%d, name results=%d\n", table.Len(), len(result.TopNames))

	// TopN does not change any state; callers decide what has been "shown".
	top := table.TopN(10)
	for _, e := range top {
		fmt.Printf("  top %s%s  %.1f (%s) poetry=%v\n",
			e.Char1, e.Char2, e.Score, e.Grade, e.HasPoetry)
		table.MarkShown(e.Char1, e.Char2)
	}
	fmt.Printf("after MarkShown(TopN(10)): shown=%d\n", table.ShownCount())

	// A nil filter accepts every unshown entry. A nil return means the explore
	// pool is exhausted; the HTTP layer normalizes that to a JSON array [].
	first := table.Explore(5, nil)
	if first == nil {
		log.Fatal("explore pool exhausted unexpectedly after only 10 shown entries")
	}
	fmt.Printf("explore(nil filter): %d entries, shown=%d\n", len(first), table.ShownCount())

	// A predicate narrows the sample, e.g. only names backed by poetry.
	poetryOnly := func(e fate.ExcellentEntry) bool { return e.HasPoetry }
	for round := 1; round <= 3; round++ {
		batch := table.Explore(5, poetryOnly)
		if batch == nil {
			fmt.Printf("round %d: exhausted (nil -> [] in the HTTP layer)\n", round)
			break
		}
		fmt.Printf("round %d: %d poetry-backed entries, shown=%d\n",
			round, len(batch), table.ShownCount())
	}

	fmt.Printf("explore complete: shown=%d total=%d\n", table.ShownCount(), table.Len())
}
