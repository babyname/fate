package repository

import (
	"sync"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/ent/wugelucky"
	"golang.org/x/net/context"
)

type CharCache struct {
	mu       sync.RWMutex
	byStroke map[int][]*ent.Character
	all      []*ent.Character
}

type LuckyCache struct {
	mu        sync.RWMutex
	byStrokes map[[2]int][]*ent.WuGeLucky
}

type ModelCache struct {
	chars *CharCache
	lucky *LuckyCache
}

func newCharCache() *CharCache {
	return &CharCache{
		byStroke: make(map[int][]*ent.Character, 128),
	}
}

func newLuckyCache() *LuckyCache {
	return &LuckyCache{
		byStrokes: make(map[[2]int][]*ent.WuGeLucky, 64),
	}
}

func (c *CharCache) Get(stroke int) ([]*ent.Character, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	chars, ok := c.byStroke[stroke]
	return chars, ok
}

func (c *CharCache) Set(stroke int, chars []*ent.Character) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.byStroke[stroke] = chars
}

func (c *CharCache) All() []*ent.Character {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.all
}

func (c *CharCache) SetAll(chars []*ent.Character) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.all = chars
}

func (c *LuckyCache) Get(strokes [2]int) ([]*ent.WuGeLucky, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	luckies, ok := c.byStrokes[strokes]
	return luckies, ok
}

func (c *LuckyCache) Set(strokes [2]int, luckies []*ent.WuGeLucky) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.byStrokes[strokes] = luckies
}

func (m *Repository) PreloadCharacters(ctx context.Context) error {
	chars, err := m.Character.Query().All(ctx)
	if err != nil {
		return err
	}
	m.cache.chars.SetAll(chars)
	return nil
}

func (m *Repository) PreloadCharactersByStroke(ctx context.Context, strokeType string) error {
	chars, err := m.Character.Query().
		Where(character.RegularEQ(true)).
		All(ctx)
	if err != nil {
		return err
	}
	grouped := make(map[int][]*ent.Character, 128)
	for _, c := range chars {
		var stroke int
		switch strokeType {
		case "chs":
			stroke = c.SimplifiedStroke
		case "cht":
			stroke = c.TraditionalStroke
		case "kangxi":
			stroke = c.KangxiStroke
		default:
			stroke = c.ScienceStroke
		}
		grouped[stroke] = append(grouped[stroke], c)
	}
	for s, cs := range grouped {
		m.cache.chars.Set(s, cs)
	}
	return nil
}

func (m *Repository) GetCharactersCached(ctx context.Context, stroke int, strokeFilter func(*ent.CharacterQuery) *ent.CharacterQuery, regularFilter func(*ent.CharacterQuery) *ent.CharacterQuery) ([]*ent.Character, error) {
	if cached, ok := m.cache.chars.Get(stroke); ok {
		return cached, nil
	}
	q := m.Character.Query()
	if strokeFilter != nil {
		q = strokeFilter(q)
	}
	if regularFilter != nil {
		q = regularFilter(q)
	}
	chars, err := q.All(ctx)
	if err != nil {
		return nil, err
	}
	m.cache.chars.Set(stroke, chars)
	return chars, nil
}

func (m *Repository) GetWuGeLuckyCached(ctx context.Context, strokes [2]int) ([]*ent.WuGeLucky, error) {
	if cached, ok := m.cache.lucky.Get(strokes); ok {
		return cached, nil
	}
	query := m.WuGeLucky.Query().
		Where(wugelucky.LastStroke1EQ(strokes[0])).
		Where(wugelucky.And(wugelucky.LastStroke2EQ(strokes[1]))).
		Where(wugelucky.And(wugelucky.ZongLuckyEQ(true)))
	luckies, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	m.cache.lucky.Set(strokes, luckies)
	return luckies, nil
}

func NewModelCache() *ModelCache {
	return &ModelCache{
		chars: newCharCache(),
		lucky: newLuckyCache(),
	}
}
