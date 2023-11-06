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

func (NCharacter) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("hash"),
		field.String("pin_yin"),
		field.Int64("ch_id"),
		field.String("ch"),
		field.String("radical"),
		field.Int("radical_stroke"),
		field.Int("total_stroke"),
		field.Bool("is_kang_xi"),
		field.String("relate_kang_xi"),
		field.String("relate_simple"),
		field.String("relate_traditional"),
		field.String("relate_variant"),
		//field.String("kang_xi"),
		//field.Int("kang_xi_stroke"),
		//field.String("simple_radical"),
		//field.Int("simple_radical_stroke"),
		//field.Int("simple_total_stroke"),
		//field.String("traditional_radical"),
		//field.Int("traditional_radical_stroke"),
		//field.Int("traditional_total_stroke"),
		field.Bool("name_science"),
		field.Int("science_stroke"),
		field.String("wu_xing"),
		field.String("lucky"),
		field.Bool("regular"),
		//field.String("traditional_character"),
		//field.String("variant_character"),
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
