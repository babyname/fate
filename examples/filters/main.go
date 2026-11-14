// Command filters demonstrates how to constrain generation with
// fate.FilterOption, and verifies that avoided characters never show up.
//
//	go run ./examples/filters
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

	avoid := []string{"病", "死", "穷"}

	// FilterOption is a plain struct with chainable With* helpers. Everything
	// the Web form can configure is available here: strokes, character sets,
	// wuxing, name style, xiyong method, poetry mode and strictness.
	opt := fate.FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: fate.CharacterFilterTypeDefault,
		RegularFilter:       true,
		DaYanFilter:         true,
		WuXingFilter:        true,
		MinStroke:           3,
		MaxStroke:           18,
		StrokeMode:          fate.StrokeModeScience,
		NameStyle:           fate.NameStyleSimplified,
	}.WithAvoidCharacters(avoid...).
		WithXiYongMethod("balance").
		WithFilterStrictness("moderate").
		WithPoetryMode(0)

	engine, err := fate.New(config.DefaultConfig())
	if err != nil {
		log.Fatalf("init fate (run 'go run ./cmd/dbinit' first): %v", err)
	}

	result, err := engine.Generate(context.Background(), fate.GenerateRequest{
		LastName: "王",
		Born:     born,
		Sex:      fate.SexBoy,
		Filter:   &opt,
	})
	if err != nil {
		log.Fatalf("generate: %v", err)
	}
	if len(result.TopNames) == 0 {
		log.Fatal("generate returned no names; check the database or loosen the filter")
	}

	fmt.Printf("filter: avoid=%v strokes=%d-%d style=simplified strictness=moderate\n",
		avoid, opt.MinStroke, opt.MaxStroke)

	for _, n := range result.TopNames {
		for _, ch := range []string{n.Char1.Char, n.Char2.Char} {
			for _, a := range avoid {
				if ch == a {
					log.Fatalf("avoided character %q appeared in %s", ch, n.FullName)
				}
			}
		}
	}

	fmt.Printf("verified %d results, none contains %v; first name: %s (%.1f)\n",
		len(result.TopNames), avoid, result.TopNames[0].FullName, result.TopNames[0].Score)
}
