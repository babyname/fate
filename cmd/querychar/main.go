package main

import (
	"context"
	"fmt"
	"os"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/ent/character"
	_ "github.com/sqlite3ent/sqlite3"
)

func main() {
	client, err := ent.Open("sqlite3", "file:fate?cache=shared&_journal=WAL&_fk=1")
	if err != nil {
		fmt.Println("Open error:", err)
		os.Exit(1)
	}
	defer client.Close()

	ctx := context.Background()

	queryChars := []string{"峰", "洛", "凌", "明", "瑞", "轩", "宇", "诗", "词", "文"}
	if len(os.Args) > 1 {
		queryChars = os.Args[1:]
	}

	for _, ch := range queryChars {
		c, err := client.Character.Query().Where(character.CharEQ(ch)).First(ctx)
		if err != nil {
			fmt.Printf("Char %q: NOT FOUND (%v)\n", ch, err)
		} else {
			fmt.Printf("Char %q: Char=%q IsSimplified=%v IsTraditional=%v ScienceStroke=%d WuXing=%q HasPoetry=%v Regular=%v Nameable=%v\n",
				ch, c.Char, c.IsSimplified, c.IsTraditional, c.ScienceStroke, c.WuXing, c.HasPoetry, c.Regular, c.Nameable)
		}
	}

	total, _ := client.Character.Query().Count(ctx)
	simplified, _ := client.Character.Query().Where(character.IsSimplifiedEQ(true)).Count(ctx)
	traditional, _ := client.Character.Query().Where(character.IsTraditionalEQ(true)).Count(ctx)
	poetryCount, _ := client.Character.Query().Where(character.HasPoetryEQ(true)).Count(ctx)
	fmt.Printf("\nTotal: %d, Simplified: %d, Traditional: %d, HasPoetry: %d\n", total, simplified, traditional, poetryCount)

	poetryChars, _ := client.Character.Query().Where(character.HasPoetryEQ(true)).Limit(10).All(ctx)
	for _, pc := range poetryChars {
		fmt.Printf("  Poetry char: %q WuXing=%q Regular=%v\n", pc.Char, pc.WuXing, pc.Regular)
	}

	poetryRegular, _ := client.Character.Query().Where(character.HasPoetryEQ(true), character.RegularEQ(true)).Count(ctx)
	fmt.Printf("\nHasPoetry AND Regular: %d\n", poetryRegular)
	poetryRegularChars, _ := client.Character.Query().Where(character.HasPoetryEQ(true), character.RegularEQ(true)).Limit(20).All(ctx)
	for _, pc := range poetryRegularChars {
		fmt.Printf("  Poetry+Regular: %q WuXing=%q\n", pc.Char, pc.WuXing)
	}

	nameable, _ := client.Character.Query().Where(character.NameableEQ(true)).Count(ctx)
	regular, _ := client.Character.Query().Where(character.RegularEQ(true)).Count(ctx)
	regularNameable, _ := client.Character.Query().Where(character.RegularEQ(true), character.NameableEQ(true)).Count(ctx)
	fmt.Printf("\nNameable: %d, Regular: %d, Regular+Nameable: %d\n", nameable, regular, regularNameable)
}
