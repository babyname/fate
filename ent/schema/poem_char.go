package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type PoemChar struct {
	ent.Schema
}

func (PoemChar) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StorageKey("id"),
		field.Int("poem_id").
			Comment("关联诗词ID"),
		field.String("char").
			NotEmpty().
			Comment("汉字"),
		field.Int("position").
			NonNegative().
			Comment("在原文中的位置（从0开始）"),
		field.String("sentence").
			Optional().
			Comment("所在完整句子"),
		field.String("context").
			Optional().
			Comment("上下文（前后各5字）"),
	}
}

func (PoemChar) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("poem", Poem.Type).
			Ref("poem_chars").
			Unique().
			Required().
			Field("poem_id"),
	}
}

func (PoemChar) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("char"),
		index.Fields("poem_id"),
		index.Fields("char", "poem_id"),
	}
}

func (PoemChar) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "poem_char"},
	}
}
