package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type WuGeLucky struct {
	ent.Schema
}

func (WuGeLucky) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}),
		field.Int("last_stroke_1"),
		field.Int("last_stroke_2"),
		field.Int("first_stroke_1"),
		field.Int("first_stroke_2"),
		field.Int("tian_ge"),
		field.String("tian_da_yan"),
		field.Int("ren_ge"),
		field.String("ren_da_yan"),
		field.Int("di_ge"),
		field.String("di_da_yan"),
		field.Int("wai_ge"),
		field.String("wai_da_yan"),
		field.Int("zong_ge"),
		field.String("zong_da_yan"),
		field.Bool("zong_lucky"),
		field.Bool("zong_sex"),
		field.Bool("zong_max"),
	}
}
func (WuGeLucky) Edges() []ent.Edge {
	return nil
}

func (WuGeLucky) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "wu_ge_lucky"},
	}
}

func (WuGeLucky) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("last_stroke_1"),
		index.Fields("last_stroke_2"),
	}
}
