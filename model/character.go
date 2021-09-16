package model

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	"github.com/godcong/fate/ent"
	"github.com/godcong/fate/ent/character"
)

func (m Model) GetCharacter(ctx context.Context, filters ...func(query *ent.CharacterQuery) *ent.CharacterQuery) (
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

// InsertOrUpdateCharacter ...
func (m Model) InsertOrUpdateCharacter(ctx context.Context, nch *ent.Character) (ch *ent.Character, e error) {

	tx, e := m.Tx(ctx)
	if e != nil {
		return nil, e
	}

	count, e := tx.Character.Query().Where(character.ID(nch.ID)).Count(ctx)
	if e != nil {
		return nil, fmt.Errorf("error updating character: %v,rollback: %v", e, tx.Rollback())
	}

	if count > 0 {
		ch, e = tx.Character.UpdateOne(nch).Save(ctx)
		if e != nil {
			return nil, fmt.Errorf("error updating character: %v,rollback: %v", e, tx.Rollback())

		}
		return ch, tx.Commit()
	}

	ch, e = tx.Character.Create().Save(ctx)
	if e != nil {
		return nil, fmt.Errorf("error updating character: %v,rollback: %v", e, tx.Rollback())

	}
	return ch, tx.Commit()
}

// Char ...
func Char(name string) func(query *ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.ChEQ(name),
			character.Or(character.KangxiEQ(name)),
			character.Or(func(selector *sql.Selector) {
				sqljson.ValueContains(character.FieldTraditionalCharacter, name)
			}))
	}
}
