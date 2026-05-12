package repository

import (
	"context"
	"fmt"

	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
)

type CharQuery func(query *ent.CharacterQuery) *ent.CharacterQuery
type StrokeQuery func(query *ent.CharacterQuery, s int) *ent.CharacterQuery

func (m *Repository) GetCharacter(ctx context.Context, filters ...CharQuery) (
	*ent.Character, error) {
	q := m.Character.Query()
	for _, filter := range filters {
		q = filter(q)
	}
	return q.First(ctx)
}

func (m *Repository) GetCharacters(ctx context.Context, filters ...CharQuery) ([]*ent.Character, error) {
	q := m.Character.Query()
	for _, filter := range filters {
		q = filter(q)
	}
	return q.All(ctx)
}

func (m *Repository) InsertOrUpdateCharacter(ctx context.Context, nch *ent.Character) (ch *ent.Character, e error) {
	tx, e := m.Tx(ctx)
	if e != nil {
		return nil, e
	}

	count, e := tx.Character.Query().Where(character.ID(nch.ID)).Count(ctx)
	if e != nil {
		return nil, fmt.Errorf("error updating character: %v,rollback: %v", e, tx.Rollback())
	}

	if count > 0 {
		ch, e = tx.Character.UpdateOne(nch).
			SetCharacter(nch).
			SetPinyin(nch.Pinyin).
			Save(ctx)
		if e != nil {
			return nil, fmt.Errorf("error updating character: %v,rollback: %v", e, tx.Rollback())

		}
		return ch, tx.Commit()
	}

	ch, e = tx.Character.Create().
		SetID(nch.ID).
		SetCharacter(nch).
		SetPinyin(nch.Pinyin).
		Save(ctx)
	if e != nil {
		return nil, fmt.Errorf("error updating character: %v,rollback: %v", e, tx.Rollback())

	}
	return ch, tx.Commit()
}

func Char(name string) CharQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.CharEQ(name))
	}
}

func Regular() func(query *ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.Regular(true))
	}
}

func StrokeKangxi(query *ent.CharacterQuery, s int) *ent.CharacterQuery {
	return query.Where(character.Or(character.KangxiStroke(s)))
}

func StrokeSimpleTotal(query *ent.CharacterQuery, s int) *ent.CharacterQuery {
	return query.Where(character.Or(character.SimplifiedStrokeEQ(s)))
}

func StrokeTraditionalTotal(query *ent.CharacterQuery, s int) *ent.CharacterQuery {
	return query.Where(character.Or(character.TraditionalStrokeEQ(s)))
}

func Stroke(s int, sqs ...StrokeQuery) CharQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		q := query.Where(character.ScienceStrokeEQ(s))
		for i := range sqs {
			q = sqs[i](q, s)
		}
		return q
	}

}
