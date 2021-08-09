package model

import (
	"context"

	"github.com/godcong/fate/ent"
)

func (m Model) GetCharacter(ctx context.Context, char string, filters ...func(query *ent.CharacterQuery) *ent.CharacterQuery) (
	*ent.Character, error) {
	q := m.Character.Query()
	for _, filter := range filters {
		q = filter(q)
	}
	return q.First(ctx)
}

func (m Model) GetCharacters(ctx context.Context, filters ...func(query *ent.CharacterQuery) *ent.CharacterQuery) (
	[]*ent.Character, error) {
	q := m.Character.Query()
	for _, filter := range filters {
		q = filter(q)
	}
	return q.All(ctx)
}
