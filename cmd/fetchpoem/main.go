package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/ent/poem"
	"github.com/babyname/fate/internal/database"
	"github.com/babyname/fate/internal/repository"
	"github.com/babyname/fate/internal/seeddb"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/sqlite3ent/sqlite3"
)

type PoetryData struct {
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Paragraphs []string `json:"paragraphs"`
	Content    []string `json:"content"`
	Strains    []string `json:"strains"`
}

func main() {
	cfg := config.DefaultConfig()
	builder := database.New(cfg.Database)
	client, err := builder.Client()
	if err != nil {
		panic(err)
	}
	repo := repository.New(client)
	ctx := context.Background()

	fmt.Println("Importing Tang poetry...")
	importFromURL(ctx, repo,
		"https://raw.githubusercontent.com/chinese-poetry/chinese-poetry/master/%E5%85%A8%E5%94%90%E8%AF%97/poet.tang.0.json",
		poem.TypeShi, "唐", "chinese-poetry")
	importFromURL(ctx, repo,
		"https://raw.githubusercontent.com/chinese-poetry/chinese-poetry/master/%E5%85%A8%E5%94%90%E8%AF%97/poet.tang.1000.json",
		poem.TypeShi, "唐", "chinese-poetry")
	importFromURL(ctx, repo,
		"https://raw.githubusercontent.com/chinese-poetry/chinese-poetry/master/%E5%85%A8%E5%94%90%E8%AF%97/poet.tang.2000.json",
		poem.TypeShi, "唐", "chinese-poetry")

	fmt.Println("Importing Song poetry...")
	importFromURL(ctx, repo,
		"https://raw.githubusercontent.com/chinese-poetry/chinese-poetry/master/%E5%85%A8%E5%AE%8B%E8%AF%8D/poet.song.0.json",
		poem.TypeCi, "宋", "chinese-poetry")
	importFromURL(ctx, repo,
		"https://raw.githubusercontent.com/chinese-poetry/chinese-poetry/master/%E5%85%A8%E5%AE%8B%E8%AF%8D/poet.song.1000.json",
		poem.TypeCi, "宋", "chinese-poetry")

	fmt.Println("Importing Song shi...")
	importFromURL(ctx, repo,
		"https://raw.githubusercontent.com/chinese-poetry/chinese-poetry/master/%E5%85%A8%E5%AE%8B%E8%AF%97/poet.song.0.json",
		poem.TypeShi, "宋", "chinese-poetry")

	fmt.Println("Done!")
}

func importFromURL(ctx context.Context, repo *repository.Repository, url string, poemType poem.Type, dynasty, source string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read %s: %v\n", url, err)
		return
	}

	var poems []PoetryData
	if err := json.Unmarshal(data, &poems); err != nil {
		fmt.Printf("Failed to parse %s: %v\n", url, err)
		return
	}

	var input []seeddb.ChinesePoetryJSON
	for _, p := range poems {
		content := strings.Join(p.Paragraphs, "\n")
		if content == "" {
			content = strings.Join(p.Content, "\n")
		}
		input = append(input, seeddb.ChinesePoetryJSON{
			Title:      p.Title,
			Author:     p.Author,
			Dynasty:    dynasty,
			Content:    content,
			Paragraphs: p.Paragraphs,
		})
	}

	if err := seeddb.ImportPoetryFromJSON(ctx, repo, input, poemType, source); err != nil {
		fmt.Printf("Failed to import %s: %v\n", url, err)
	}

	fmt.Printf("Imported %d poems from %s\n", len(poems), url)
}
