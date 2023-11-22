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
		field.Strings("pin_yin"),
		field.String("word"),
		field.Int("derivation"),
		field.String("explanation"),
		field.Int("abbreviation"),
		field.String("example"),
		field.String("comment"),
	}
}

// Edges of the Idiom.
func (Idiom) Edges() []ent.Edge {
	return nil
}
