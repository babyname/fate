// Command quickstart demonstrates the primary FATE API in v4.3.0:
// create an engine with fate.New, then call Generate synchronously.
//
// Run from the repository root (the SQLite file "fate" lives there):
//
//	go run ./cmd/dbinit          # once, creates the database
//	go run ./examples/quickstart
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

	// config.DefaultConfig() points at the local SQLite file named "fate".
	engine, err := fate.New(config.DefaultConfig())
	if err != nil {
		log.Fatalf("init fate (run 'go run ./cmd/dbinit' first): %v", err)
	}

	// DefaultFilterOption() is the single source of truth for bounded scoring.
	// A nil Filter resolves to exactly this option; spelling it out here shows
	// how to customise it, e.g. .WithAvoidCharacters("病", "死").
	opt := fate.DefaultFilterOption()
	result, err := engine.Generate(context.Background(), fate.GenerateRequest{
		LastName: "张",
		Born:     born,
		Sex:      fate.SexBoy,
		Filter:   &opt,
	})
	if err != nil {
		log.Fatalf("generate: %v", err)
	}
	if len(result.TopNames) == 0 {
		log.Fatal("generate returned no names; check the database")
	}

	fmt.Printf("generated %d name results, %d entries in the excellent table\n",
		len(result.TopNames), result.ExcellentTable.Len())
	for i, n := range result.TopNames {
		if i >= 5 {
			break
		}
		fmt.Printf("  #%d %s  %.1f (%s)\n", n.Rank, n.FullName, n.Score, n.Grade)
	}
}
