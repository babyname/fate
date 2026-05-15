package poetry

import (
	"context"
	"sort"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/ent/poemchar"
)

type PoetryMode int

const (
	PoetryModeNone   PoetryMode = iota
	PoetryModePrefer
	PoetryModeOnly
)

type PoetrySource struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Dynasty  string `json:"dynasty"`
	Sentence string `json:"sentence"`
	Type     string `json:"type"`
}

type CharResult struct {
	Char    string
	WuXing  string
	Sources []PoetrySource
}

type PairResult struct {
	Char1   *ent.Character
	Char2   *ent.Character
	Sources []PoetrySource
	Score   float64
}

type PoemCharWithSource struct {
	PoemChar *ent.PoemChar
	Poem     *ent.Poem
}

type Namer struct {
	client *ent.Client
}

func NewNamer(client *ent.Client) *Namer {
	return &Namer{client: client}
}

func (n *Namer) FindCharsByWuXing(ctx context.Context, wuxing string, limit int) ([]CharResult, error) {
	chars, err := n.client.Character.Query().
		Where(
			character.WuXingEQ(wuxing),
			character.RegularEQ(true),
			character.NameableEQ(true),
		).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var results []CharResult
	for _, c := range chars {
		poemChars, err := n.client.PoemChar.Query().
			Where(poemchar.CharEQ(c.Char)).
			WithPoem().
			Limit(3).
			All(ctx)
		if err != nil || len(poemChars) == 0 {
			continue
		}

		var sources []PoetrySource
		for _, pc := range poemChars {
			p := pc.Edges.Poem
			if p != nil {
				sources = append(sources, PoetrySource{
					Title:    p.Title,
					Author:   p.Author,
					Dynasty:  p.Dynasty,
					Sentence: pc.Sentence,
					Type:     string(p.Type),
				})
			}
		}

		results = append(results, CharResult{
			Char:    c.Char,
			WuXing:  c.WuXing,
			Sources: sources,
		})
	}

	return results, nil
}

func (n *Namer) FindPairsByWuXing(ctx context.Context, wuxing1, wuxing2 string, limit int) ([]PairResult, error) {
	chars1, err := n.client.Character.Query().
		Where(
			character.WuXingEQ(wuxing1),
			character.RegularEQ(true),
			character.NameableEQ(true),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	chars2, err := n.client.Character.Query().
		Where(
			character.WuXingEQ(wuxing2),
			character.RegularEQ(true),
			character.NameableEQ(true),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	char1Map := make(map[string]*ent.Character)
	for _, c := range chars1 {
		char1Map[c.Char] = c
	}
	char2Map := make(map[string]*ent.Character)
	for _, c := range chars2 {
		char2Map[c.Char] = c
	}

	poemChars1, err := n.client.PoemChar.Query().
		Where(poemchar.CharIn(mapKeys(char1Map)...)).
		WithPoem().
		All(ctx)
	if err != nil {
		return nil, err
	}

	type poemInfo struct {
		poem   *ent.Poem
		chars1 map[string]*PoemCharWithSource
		chars2 map[string]*PoemCharWithSource
	}

	poemMap := make(map[int]*poemInfo)
	for _, pc := range poemChars1 {
		p := pc.Edges.Poem
		if p == nil {
			continue
		}
		if _, ok := poemMap[pc.PoemID]; !ok {
			poemMap[pc.PoemID] = &poemInfo{poem: p, chars1: make(map[string]*PoemCharWithSource), chars2: make(map[string]*PoemCharWithSource)}
		}
		poemMap[pc.PoemID].chars1[pc.Char] = &PoemCharWithSource{PoemChar: pc, Poem: p}
	}

	poemChars2, err := n.client.PoemChar.Query().
		Where(poemchar.CharIn(mapKeys(char2Map)...)).
		WithPoem().
		All(ctx)
	if err != nil {
		return nil, err
	}

	for _, pc := range poemChars2 {
		if info, ok := poemMap[pc.PoemID]; ok {
			p := pc.Edges.Poem
			if p == nil && info.poem != nil {
				p = info.poem
			}
			info.chars2[pc.Char] = &PoemCharWithSource{PoemChar: pc, Poem: p}
		}
	}

	var pairs []PairResult
	for _, info := range poemMap {
		for c1Char, pc1 := range info.chars1 {
			for c2Char, pc2 := range info.chars2 {
				if c1Char == c2Char {
					continue
				}
				c1 := char1Map[c1Char]
				c2 := char2Map[c2Char]
				if c1 == nil || c2 == nil {
					continue
				}

				sameSentence := pc1.PoemChar.Sentence == pc2.PoemChar.Sentence && pc1.PoemChar.Sentence != ""

				score := 60.0
				if sameSentence {
					score += 25.0
				} else {
					score += 10.0
				}
				switch info.poem.Type {
				case "jing":
					score += 10.0
				case "shi":
					score += 5.0
				case "ci":
					score += 3.0
				}

				sources := []PoetrySource{{
					Title:    info.poem.Title,
					Author:   info.poem.Author,
					Dynasty:  info.poem.Dynasty,
					Sentence: pc1.PoemChar.Sentence,
					Type:     string(info.poem.Type),
				}}

				pairs = append(pairs, PairResult{
					Char1:   c1,
					Char2:   c2,
					Sources: sources,
					Score:   score,
				})
			}
		}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Score > pairs[j].Score
	})
	if limit > 0 && len(pairs) > limit {
		pairs = pairs[:limit]
	}

	return pairs, nil
}

func (n *Namer) FindCharPoetry(ctx context.Context, char string) ([]PoetrySource, error) {
	poemChars, err := n.client.PoemChar.Query().
		Where(poemchar.CharEQ(char)).
		WithPoem().
		Limit(5).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var sources []PoetrySource
	for _, pc := range poemChars {
		p := pc.Edges.Poem
		if p != nil {
			sources = append(sources, PoetrySource{
				Title:    p.Title,
				Author:   p.Author,
				Dynasty:  p.Dynasty,
				Sentence: pc.Sentence,
				Type:     string(p.Type),
			})
		}
	}
	return sources, nil
}

func mapKeys(m map[string]*ent.Character) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}
