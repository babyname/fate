package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Poem struct {
	ent.Schema
}

func (Poem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StorageKey("id"),
		field.String("title").
			NotEmpty().
			Comment("诗名/词牌名"),
		field.String("author").
			Optional().
			Comment("作者"),
		field.String("dynasty").
			Optional().
			Comment("朝代，如唐/宋/先秦"),
		field.Text("content").
			NotEmpty().
			Comment("全文内容"),
		field.String("preface").
			Optional().
			Comment("序言"),
		field.Strings("keywords").
			Optional().
			Comment("可用于起名的关键字列表"),
		field.Strings("tags").
			Optional().
			Comment("风格标签，如豪放/婉约/送别/咏物"),
		field.Enum("type").
			Values("shi", "ci", "fu", "jing", "other").
			Default("shi").
			Comment("体裁：诗/词/赋/经/其他"),
		field.String("source").
			Optional().
			Comment("数据来源"),
	}
}

func (Poem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("poem_chars", PoemChar.Type).
			Comment("诗词中的字"),
	}
}

func (Poem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title", "author"),
		index.Fields("dynasty"),
		index.Fields("type"),
	}
}

func (Poem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "poem"},
	}
}
