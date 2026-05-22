package main

import (
	"context"
	"fmt"
	"os"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
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

	for _, ch := range []string{"张", "張", "李", "王"} {
		c, err := client.Character.Query().Where(character.CharEQ(ch)).First(ctx)
		if err != nil {
			fmt.Printf("Char %q: NOT FOUND (%v)\n", ch, err)
		} else {
			fmt.Printf("Char %q: Char=%q IsSimplified=%v IsTraditional=%v ScienceStroke=%d WuXing=%q\n",
				ch, c.Char, c.IsSimplified, c.IsTraditional, c.ScienceStroke, c.WuXing)
		}
	}

	total, _ := client.Character.Query().Count(ctx)
	simplified, _ := client.Character.Query().Where(character.IsSimplifiedEQ(true)).Count(ctx)
	traditional, _ := client.Character.Query().Where(character.IsTraditionalEQ(true)).Count(ctx)
	fmt.Printf("\nTotal: %d, Simplified: %d, Traditional: %d\n", total, simplified, traditional)
}
