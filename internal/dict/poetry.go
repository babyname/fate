package dict

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type TangShi struct {
	Author    string   `json:"author"`
	Paragraphs []string `json:"paragraphs"`
	Title     string   `json:"title"`
	ID        string   `json:"id"`
}

type SongCi struct {
	Author    string   `json:"author"`
	Paragraphs []string `json:"paragraphs"`
	Rhythmic  string   `json:"rhythmic"`
}

type ShiJing struct {
	Title   string   `json:"title"`
	Chapter string   `json:"chapter"`
	Section string   `json:"section"`
	Content []string `json:"content"`
}

type PoemEntry struct {
	Title    string
	Author   string
	Dynasty  string
	Content  string
	Type     string
	Source   string
}

type CharPoetryRef struct {
	Char     string
	Position int
	Sentence string
	Context  string
}

func LoadTangShi(path string) ([]*TangShi, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open tang shi: %w", err)
	}
	defer f.Close()

	var entries []*TangShi
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode tang shi: %w", err)
	}
	return entries, nil
}

func LoadSongCi(path string) ([]*SongCi, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open song ci: %w", err)
	}
	defer f.Close()

	var entries []*SongCi
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode song ci: %w", err)
	}
	return entries, nil
}

func LoadShiJing(path string) ([]*ShiJing, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open shi jing: %w", err)
	}
	defer f.Close()

	var entries []*ShiJing
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode shi jing: %w", err)
	}
	return entries, nil
}

func TangShiToPoemEntries(ts []*TangShi) []*PoemEntry {
	entries := make([]*PoemEntry, 0, len(ts))
	for _, t := range ts {
		content := strings.Join(t.Paragraphs, "\n")
		entries = append(entries, &PoemEntry{
			Title:   t.Title,
			Author:  t.Author,
			Dynasty: "唐",
			Content: content,
			Type:    "shi",
			Source:  "chinese-poetry",
		})
	}
	return entries
}

func SongCiToPoemEntries(ci []*SongCi) []*PoemEntry {
	entries := make([]*PoemEntry, 0, len(ci))
	for _, c := range ci {
		title := c.Rhythmic
		if title == "" {
			title = "无名"
		}
		content := strings.Join(c.Paragraphs, "\n")
		entries = append(entries, &PoemEntry{
			Title:   title,
			Author:  c.Author,
			Dynasty: "宋",
			Content: content,
			Type:    "ci",
			Source:  "chinese-poetry",
		})
	}
	return entries
}

func ShiJingToPoemEntries(sj []*ShiJing) []*PoemEntry {
	entries := make([]*PoemEntry, 0, len(sj))
	for _, s := range sj {
		content := strings.Join(s.Content, "\n")
		author := s.Section
		if s.Chapter != "" {
			author = s.Chapter + "·" + s.Section
		}
		entries = append(entries, &PoemEntry{
			Title:   s.Title,
			Author:  author,
			Dynasty: "先秦",
			Content: content,
			Type:    "jing",
			Source:  "chinese-poetry",
		})
	}
	return entries
}

func ExtractCharRefs(content string) []*CharPoetryRef {
	runes := []rune(content)
	fullText := strings.ReplaceAll(content, "\n", "")
	textRunes := []rune(fullText)

	var refs []*CharPoetryRef
	sentences := SplitSentences(fullText)
	charSentences := make(map[rune]string)
	for _, s := range sentences {
		for _, r := range s {
			if unicode.Is(unicode.Han, r) && charSentences[r] == "" {
				charSentences[r] = s
			}
		}
	}

	for i, r := range textRunes {
		if !unicode.Is(unicode.Han, r) {
			continue
		}
		sentence := charSentences[r]
		context := extractContextFromRunes(textRunes, i, 5)
		refs = append(refs, &CharPoetryRef{
			Char:     string(r),
			Position: i,
			Sentence: sentence,
			Context:  context,
		})
	}
	_ = runes
	return refs
}

func SplitSentences(text string) []string {
	var sentences []string
	var current []rune

	for _, r := range text {
		current = append(current, r)
		if r == '。' || r == '，' || r == '！' || r == '？' || r == '；' || r == '、' {
			s := strings.TrimSpace(string(current))
			if len(s) > 1 {
				sentences = append(sentences, s)
			}
			current = current[:0]
		}
	}

	if len(current) > 0 {
		s := strings.TrimSpace(string(current))
		if len(s) > 1 {
			sentences = append(sentences, s)
		}
	}

	return sentences
}

func extractContextFromRunes(runes []rune, pos int, radius int) string {
	start := pos - radius
	if start < 0 {
		start = 0
	}
	end := pos + radius + 1
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

func LoadPoetryFromDir(dir string) ([]*PoemEntry, error) {
	return loadPoetryFromDir(dir, false)
}

func LoadSelectedPoetryFromDir(dir string) ([]*PoemEntry, error) {
	return loadPoetryFromDir(dir, true)
}

func loadPoetryFromDir(dir string, selectedOnly bool) ([]*PoemEntry, error) {
	var allEntries []*PoemEntry

	tangDir := filepath.Join(dir, "全唐诗")
	if info, err := os.Stat(tangDir); err == nil && info.IsDir() {
		selectedFiles := []string{
			filepath.Join(tangDir, "唐诗三百首.json"),
		}
		for _, f := range selectedFiles {
			if _, err := os.Stat(f); err != nil {
				continue
			}
			ts, err := LoadTangShi(f)
			if err != nil {
				continue
			}
			allEntries = append(allEntries, TangShiToPoemEntries(ts)...)
		}

		if !selectedOnly {
			patterns := []string{
				filepath.Join(tangDir, "poet.tang.*.json"),
			}
			for _, pattern := range patterns {
				files, err := filepath.Glob(pattern)
				if err != nil {
					continue
				}
				for _, f := range files {
					ts, err := LoadTangShi(f)
					if err != nil {
						continue
					}
					allEntries = append(allEntries, TangShiToPoemEntries(ts)...)
				}
			}
		}
	}

	ciDir := filepath.Join(dir, "宋词")
	if info, err := os.Stat(ciDir); err == nil && info.IsDir() {
		selectedFiles := []string{
			filepath.Join(ciDir, "宋词三百首.json"),
		}
		for _, f := range selectedFiles {
			if _, err := os.Stat(f); err != nil {
				continue
			}
			ci, err := LoadSongCi(f)
			if err != nil {
				continue
			}
			allEntries = append(allEntries, SongCiToPoemEntries(ci)...)
		}

		if !selectedOnly {
			patterns := []string{
				filepath.Join(ciDir, "ci.song.*.json"),
			}
			for _, pattern := range patterns {
				files, err := filepath.Glob(pattern)
				if err != nil {
					continue
				}
				for _, f := range files {
					ci, err := LoadSongCi(f)
					if err != nil {
						continue
					}
					allEntries = append(allEntries, SongCiToPoemEntries(ci)...)
				}
			}
		}
	}

	sjPath := filepath.Join(dir, "诗经", "shijing.json")
	if _, err := os.Stat(sjPath); err == nil {
		sj, err := LoadShiJing(sjPath)
		if err == nil {
			allEntries = append(allEntries, ShiJingToPoemEntries(sj)...)
		}
	}

	return allEntries, nil
}
