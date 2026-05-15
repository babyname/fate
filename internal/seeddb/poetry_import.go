package seeddb

import (
	"context"
	"strings"
	"unicode"

	"github.com/babyname/fate/ent/poem"
	"github.com/babyname/fate/internal/repository"
)

type ChinesePoetryJSON struct {
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Dynasty    string   `json:"dynasty"`
	Content    string   `json:"content"`
	Paragraphs []string `json:"paragraphs"`
}

func ImportPoetryFromJSON(ctx context.Context, repo *repository.Repository, poems []ChinesePoetryJSON, poemType poem.Type, source string) error {
	for _, p := range poems {
		content := p.Content
		if content == "" && len(p.Paragraphs) > 0 {
			content = strings.Join(p.Paragraphs, "\n")
		}
		if content == "" {
			continue
		}

		keywords := extractNameableChars(content)

		pm, err := repo.InsertPoem(ctx, p.Title, p.Author, p.Dynasty, content, "", keywords, nil, poemType, source)
		if err != nil {
			continue
		}

		position := 0
		sentences := splitSentences(content)
		charSentenceMap := buildCharSentenceMap(sentences)

		for _, ch := range content {
			if !unicode.Is(unicode.Han, ch) {
				position++
				continue
			}
			charStr := string(ch)
			sentence := charSentenceMap[charStr]
			context5 := extractContext(content, position, 5)

			_, err := repo.InsertPoemChar(ctx, pm.ID, charStr, position, sentence, context5)
			if err != nil {
				position++
				continue
			}
			position++
		}
	}
	return nil
}

func extractNameableChars(content string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, ch := range content {
		if unicode.Is(unicode.Han, ch) {
			s := string(ch)
			if !seen[s] {
				seen[s] = true
				result = append(result, s)
			}
		}
	}
	return result
}

func splitSentences(content string) []string {
	separators := []string{"。", "！", "？", "；", ".", "!", "?"}
	result := []string{content}
	for _, sep := range separators {
		var newResult []string
		for _, s := range result {
			parts := strings.Split(s, sep)
			for i, part := range parts {
				if i < len(parts)-1 {
					newResult = append(newResult, part+sep)
				} else if part != "" {
					newResult = append(newResult, part)
				}
			}
		}
		result = newResult
	}
	return result
}

func buildCharSentenceMap(sentences []string) map[string]string {
	m := make(map[string]string)
	for _, s := range sentences {
		for _, ch := range s {
			if unicode.Is(unicode.Han, ch) {
				charStr := string(ch)
				if _, exists := m[charStr]; !exists {
					m[charStr] = strings.TrimSpace(s)
				}
			}
		}
	}
	return m
}

func extractContext(content string, position, radius int) string {
	runes := []rune(content)
	start := position - radius
	if start < 0 {
		start = 0
	}
	end := position + radius + 1
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}
