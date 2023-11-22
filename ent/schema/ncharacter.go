package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type NCharacter struct {
	ent.Schema
}

const (
	CharTypeUnknown     = 0x00
	CharTypeSimple      = 0x01
	CharTypeTraditional = 0x02
	CharTypeKangXi      = 0x04
	CharTypeVariant     = 0x08
)

func (NCharacter) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("id"), //rune code
		field.String("pin_yin"),
		field.String("ch"),
		field.Int("ch_stroke"),
		field.Int("ch_type").Default(CharTypeUnknown),
		field.String("radical"),
		field.Int("radical_stroke"),
		field.Int32("relate"), // relate simple chinese
		field.Int32("relate_kang_xi"),
		field.Int32("relate_traditional"),
		field.Strings("relate_variant"), //relate other variant characters
		field.Bool("is_name_science"),
		field.Int("name_science_ch_stroke"),
		field.Bool("is_regular"),
		field.String("wu_xing"),
		field.String("lucky"),
		field.String("explanation"),
		field.String("comment"),
	}
}
func (NCharacter) Edges() []ent.Edge {
	return nil
}
func (NCharacter) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "n_character"},
	}
}
