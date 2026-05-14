package repository

import (
	"sync"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
	"golang.org/x/net/context"
)

type CharCache struct {
	mu       sync.RWMutex
	byStroke map[int][]*ent.Character
	all      []*ent.Character
}

type ModelCache struct {
	chars *CharCache
}

func newCharCache() *CharCache {
	return &CharCache{
		byStroke: make(map[int][]*ent.Character, 128),
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

func NewModelCache() *ModelCache {
	return &ModelCache{
		chars: newCharCache(),
	}
}
