package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Idiom holds the schema definition for the Idiom entity.
type Idiom struct {
	ent.Schema
}

// Fields of the Idiom.
func (Idiom) Fields() []ent.Field {
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

// Edges of the Idiom.
func (Idiom) Edges() []ent.Edge {
	return nil
}
