package main

import (
	"context"
	"fmt"
	"log"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	client, err := ent.Open("sqlite3", "file:fate?cache=shared&_journal=WAL&_fk=1")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	ctx := context.Background()

	total, err := client.Character.Query().Count(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total characters in DB: %d\n", total)

	for _, ch := range []string{"西", "门", "李", "王", "张", "赵", "东", "南"} {
		n, err := client.Character.Query().Where(character.CharEQ(ch)).Count(ctx)
		if err != nil {
			fmt.Printf("  %s: error %v\n", ch, err)
			continue
		}
		if n > 0 {
			c, _ := client.Character.Query().Where(character.CharEQ(ch)).First(ctx)
			fmt.Printf("  %s: found (stroke=%d, kangxi=%d, science=%d, wuxing=%s, regular=%v, nameable=%v)\n",
				ch, c.SimplifiedStroke, c.KangxiStroke, c.ScienceStroke, c.WuXing, c.Regular, c.Nameable)
		} else {
			fmt.Printf("  %s: NOT FOUND\n", ch)
		}
	}

	regularCount, _ := client.Character.Query().Where(character.RegularEQ(true)).Count(ctx)
	fmt.Printf("Regular characters: %d\n", regularCount)

	nameableCount, _ := client.Character.Query().Where(character.NameableEQ(true)).Count(ctx)
	fmt.Printf("Nameable characters: %d\n", nameableCount)
}
