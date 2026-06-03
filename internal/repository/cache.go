// Package repository 提供字符数据的存储、缓存与查询功能。
package repository

import (
	"context"
	"sync"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/ent/character"
)

// CharCache 按笔画数缓存字符数据的并发安全缓存。
type CharCache struct {
	mu       sync.RWMutex
	byStroke map[int][]*ent.Character
	all      []*ent.Character
}

// ModelCache 管理字符数据的缓存模型。
type ModelCache struct {
	chars *CharCache
}

func newCharCache() *CharCache {
	return &CharCache{
		byStroke: make(map[int][]*ent.Character, 128),
	}
}

// Get 根据笔画数获取缓存的字符列表。
func (c *CharCache) Get(stroke int) ([]*ent.Character, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	chars, ok := c.byStroke[stroke]
	return chars, ok
}

// Set 设置指定笔画数的字符缓存。
func (c *CharCache) Set(stroke int, chars []*ent.Character) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.byStroke[stroke] = chars
}

// All 返回所有已缓存的字符列表。
func (c *CharCache) All() []*ent.Character {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.all
}

// SetAll 设置全部字符的缓存。
func (c *CharCache) SetAll(chars []*ent.Character) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.all = chars
}

// PreloadCharacters 预加载所有字符数据到缓存。
func (m *Repository) PreloadCharacters(ctx context.Context) error {
	chars, err := m.Character.Query().All(ctx)
	if err != nil {
		return err
	}
	m.cache.chars.SetAll(chars)
	return nil
}

// PreloadCharactersByStroke 按笔画类型预加载常用字符数据到缓存。
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

// GetCharactersCached 根据笔画数和过滤条件获取字符，优先从缓存读取。
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

// NewModelCache 创建新的缓存模型实例。
func NewModelCache() *ModelCache {
	return &ModelCache{
		chars: newCharCache(),
	}
}
