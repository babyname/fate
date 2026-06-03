package repository

import (
	"context"
	"fmt"

	"github.com/babyname/fate/v4/ent"
	"github.com/babyname/fate/v4/ent/character"
)

// CharQuery 定义字符查询的过滤函数类型。
type CharQuery func(query *ent.CharacterQuery) *ent.CharacterQuery

// StrokeQuery 定义按笔画查询的过滤函数类型。
type StrokeQuery func(query *ent.CharacterQuery, s int) *ent.CharacterQuery

// GetCharacter 根据过滤条件查询单个字符。
func (m *Repository) GetCharacter(ctx context.Context, filters ...CharQuery) (
	*ent.Character, error) {
	q := m.Character.Query()
	for _, filter := range filters {
		q = filter(q)
	}
	return q.First(ctx)
}

// GetCharacters 根据过滤条件查询多个字符。
func (m *Repository) GetCharacters(ctx context.Context, filters ...CharQuery) ([]*ent.Character, error) {
	q := m.Character.Query()
	for _, filter := range filters {
		q = filter(q)
	}
	return q.All(ctx)
}

// InsertOrUpdateCharacter 插入或更新字符记录。
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

// Char 返回按字符名称过滤的查询函数。
func Char(name string) CharQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.CharEQ(name))
	}
}

// Regular 返回筛选常用字的查询函数。
func Regular() func(query *ent.CharacterQuery) *ent.CharacterQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		return query.Where(character.Regular(true))
	}
}

// StrokeKangxi 返回按康熙笔画筛选的查询条件。
func StrokeKangxi(query *ent.CharacterQuery, s int) *ent.CharacterQuery {
	return query.Where(character.Or(character.KangxiStroke(s)))
}

// StrokeSimpleTotal 返回按简体笔画筛选的查询条件。
func StrokeSimpleTotal(query *ent.CharacterQuery, s int) *ent.CharacterQuery {
	return query.Where(character.Or(character.SimplifiedStrokeEQ(s)))
}

// StrokeTraditionalTotal 返回按繁体笔画筛选的查询条件。
func StrokeTraditionalTotal(query *ent.CharacterQuery, s int) *ent.CharacterQuery {
	return query.Where(character.Or(character.TraditionalStrokeEQ(s)))
}

// Stroke 返回按科学笔画筛选的查询函数，支持附加笔画查询条件。
func Stroke(s int, sqs ...StrokeQuery) CharQuery {
	return func(query *ent.CharacterQuery) *ent.CharacterQuery {
		q := query.Where(character.ScienceStrokeEQ(s))
		for i := range sqs {
			q = sqs[i](q, s)
		}
		return q
	}

}
