// Command stream demonstrates GenerateStream: parallel scoring delivers
// incremental top-result batches, and exactly one final batch (IsFinal=true)
// carries the completed result.
//
//	go run ./examples/stream
package main

import (
	"context"
	"fmt"
	"log"
	"sync"
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

	// onBatch may be called from several scoring workers concurrently.
	// Keep all shared state inside the callback under one mutex.
	var (
		mu           sync.Mutex
		batches      int
		finalBatches int
		lastWasFinal bool
	)

	err = engine.GenerateStream(context.Background(), fate.GenerateRequest{
		LastName: "李",
		Born:     born,
		Sex:      fate.SexGirl,
		// Bounded filter: the permissive DefaultFilter makes the candidate
		// cross-product impractically large. See examples/quickstart.
		Filter: &fate.FilterOption{
			CharacterFilter:     true,
			CharacterFilterType: fate.CharacterFilterTypeDefault,
			MinStroke:           3,
			MaxStroke:           18,
			RegularFilter:       true,
			DaYanFilter:         true,
			WuXingFilter:        true,
		},
	}, func(batch fate.StreamBatch) {
		mu.Lock()
		defer mu.Unlock()

		batches++
		lastWasFinal = batch.IsFinal
		if batch.IsFinal {
			finalBatches++
		}
		fmt.Printf("batch %d: total=%d preview=%d final=%v\n",
			batches, batch.Total, len(batch.TopResults), batch.IsFinal)
	})
	if err != nil {
		log.Fatalf("generate stream: %v", err)
	}

	if batches == 0 || finalBatches != 1 || !lastWasFinal {
		log.Fatalf("expected exactly one final batch as the last one; got %d batches, %d final, lastWasFinal=%v",
			batches, finalBatches, lastWasFinal)
	}
	fmt.Printf("stream complete: %d batches, final batch delivered last\n", batches)
}
