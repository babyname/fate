package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

const CurrentDataVersion = 1

// Version holds the schema definition for the Version entity.
type Version struct {
	ent.Schema
}

// Fields of the Version.
func (Version) Fields() []ent.Field {
	return []ent.Field{
		field.Int("version"),
		field.Int64("UpdatedUnix"),
	}
}

// Edges of the Version.
func (Version) Edges() []ent.Edge {
	return nil
}
