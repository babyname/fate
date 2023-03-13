package fate

import (
	"github.com/babyname/fate/ent"
	"github.com/babyname/fate/ent/character"
)

type Filter struct {
	MinCharacter int
	MaxCharacter int
}

func DefaultProperty() *Filter {
	return &Filter{
		MinCharacter: 3,
		MaxCharacter: 18,
	}
}

func (p Filter) CharacterFilter(query *ent.CharacterQuery) *ent.CharacterQuery {
	return query.Where(character.StrokeGTE(p.MinCharacter)).Where(character.And(character.StrokeLTE(p.MaxCharacter)))
}

//func (p Filter) CharacterFilter(query *ent.CharacterQuery) *ent.CharacterQuery {
//	return query.Where(character.StrokeGTE(p.MinCharacter)).Where(character.And(character.StrokeLTE(p.MaxCharacter)))
//}
