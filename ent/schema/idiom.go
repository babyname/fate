package schema

import "entgo.io/ent"

// Idiom holds the schema definition for the Idiom entity.
type Idiom struct {
	ent.Schema
}

// Fields of the Idiom.
func (Idiom) Fields() []ent.Field {
	return nil
}

// Edges of the Idiom.
func (Idiom) Edges() []ent.Edge {
	return nil
}
