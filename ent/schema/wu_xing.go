package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type WuXing struct {
	ent.Schema
}

func (WuXing) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.Time("created").Optional(),
		field.Time("updated").Optional(),
		field.Time("deleted").Optional(),
		field.Int("version").Optional(),
		field.String("first").Optional(),
		field.String("second").Optional(),
		field.String("third").Optional(),
		field.String("fortune").Optional()}
}
func (WuXing) Edges() []ent.Edge {
	return nil
}
func (WuXing) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "wu_xing"},
	}
}
