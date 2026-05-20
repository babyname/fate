package main

import (
	"context"
	"fmt"
	"log"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/internal/database"
	"github.com/babyname/fate/internal/repository"
	"github.com/babyname/fate/internal/seeddb"
)

func main() {
	cfg := config.DefaultConfig()
	b := database.New(cfg.Database)
	client, err := b.Client()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	repo := repository.New(client)

	count, _ := repo.Character.Query().Count(context.Background())
	fmt.Printf("Total characters: %d\n", count)

	testChars := []string{"西", "门", "东", "南", "北", "上", "司", "马", "诸", "葛", "欧", "阳", "王", "李", "张"}
	for _, ch := range testChars {
		c, err := repo.Character.Query().Where(character.CharEQ(ch)).First(context.Background())
		if err != nil {
			fmt.Printf("  %q: NOT FOUND (%v)\n", ch, err)
		} else {
			fmt.Printf("  %q: stroke=%d wuxing=%s pinyin=%v source=%s\n", ch, c.ScienceStroke, c.WuXing, c.Pinyin, c.Source)
		}
	}

	seeds, _ := seeddb.LoadEmbeddedCharacters()
	fmt.Printf("\nEmbedded seeds: %d\n", len(seeds))

	seedMap := make(map[string]bool)
	for _, s := range seeds {
		seedMap[s.Char] = true
	}
	for _, ch := range testChars {
		fmt.Printf("  %q in seeds: %v\n", ch, seedMap[ch])
	}
}
