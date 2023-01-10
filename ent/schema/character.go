package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Character struct {
	ent.Schema
}

func (Character) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("hash"),
		field.String("pin_yin"),
		field.String("ch"),
		field.String("radical"),
		field.Int("radical_stroke"),
		field.Int("stroke"),
		field.Bool("is_kang_xi"),
		field.String("kang_xi"),
		field.Int("kang_xi_stroke"),
		field.String("simple_radical"),
		field.Int("simple_radical_stroke"),
		field.Int("simple_total_stroke"),
		field.String("traditional_radical"),
		field.Int("traditional_radical_stroke"),
		field.Int("traditional_total_stroke"),
		field.Bool("name_science"),
		field.String("wu_xing"),
		field.String("lucky"),
		field.Bool("regular"),
		field.String("traditional_character"),
		field.String("variant_character"),
		field.String("comment"),
		field.Int("science_stroke")}
}
func (Character) Edges() []ent.Edge {
	return nil
}
func (Character) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "character"},
	}
}
