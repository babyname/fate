package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Character struct {
	ent.Schema
}

func (Character) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StorageKey("id"),
		field.String("char").
			NotEmpty().
			Comment("汉字本身（Unicode标准形式）"),
		field.String("unicode").
			Optional().
			Comment("Unicode码点，如U+4E00"),
		field.Bool("is_simplified").
			Default(false).
			Comment("是否为简体字形"),
		field.Bool("is_traditional").
			Default(false).
			Comment("是否为繁体字形"),
		field.Bool("is_kangxi").
			Default(false).
			Comment("是否为康熙字典字形"),
		field.Bool("is_variant").
			Default(false).
			Comment("是否为异体字"),
		field.Bool("is_ancient").
			Default(false).
			Comment("是否为古字/旧字形"),
		field.Strings("pinyin").
			Optional().
			Comment("拼音列表（多音字支持），如[zhōng,zhòng]"),
		field.String("radical").
			Optional().
			Comment("部首"),
		field.Int("radical_stroke").
			Optional().
			NonNegative().
			Comment("部首笔画数"),
		field.Int("simplified_stroke").
			Optional().
			NonNegative().
			Comment("简体笔画数（GB13000/通用规范汉字表）"),
		field.Int("traditional_stroke").
			Optional().
			NonNegative().
			Comment("繁体笔画数"),
		field.Int("kangxi_stroke").
			Optional().
			NonNegative().
			Comment("康熙字典笔画数（起名用）"),
		field.Int("science_stroke").
			Optional().
			NonNegative().
			Comment("姓名学笔画数（基于康熙字典，含部首变形修正）"),
		field.String("wu_xing").
			Optional().
			Comment("五行属性：木/火/土/金/水"),
		field.Bool("regular").
			Default(false).
			Comment("是否常用字（通用规范汉字表内）"),
		field.Int("common_level").
			Optional().
			NonNegative().
			Comment("常用字等级 1-5（1最常用）"),
		field.String("gender_hint").
			Optional().
			Comment("性别倾向：male/female/neutral"),
		field.Bool("nameable").
			Default(true).
			Comment("是否可用于起名"),
		field.String("meaning").
			Optional().
			Comment("字义简释"),
		field.String("source").
			Optional().
			Comment("数据来源标识，如unihan/kangxi/custom"),
		field.Float("source_confidence").
			Optional().
			Comment("数据来源可信度 0-1"),
		field.String("comment").
			Optional().
			Comment("备注"),
	}
}

func (Character) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("simplified_of", Character.Type).
			Ref("traditional_to_simplified").
			Comment("简体字对应的繁体字"),
		edge.To("traditional_to_simplified", Character.Type).
			Comment("繁体字对应的简体字").
			Unique(),
		edge.From("variant_of", Character.Type).
			Ref("standard_to_variant").
			Comment("异体字对应的标准字").
			Unique(),
		edge.To("standard_to_variant", Character.Type).
			Comment("标准字对应的异体字"),
	}
}

func (Character) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("char").Unique(),
		index.Fields("unicode"),
		index.Fields("is_simplified"),
		index.Fields("is_traditional"),
		index.Fields("is_kangxi"),
		index.Fields("is_variant"),
		index.Fields("wu_xing"),
		index.Fields("simplified_stroke"),
		index.Fields("traditional_stroke"),
		index.Fields("kangxi_stroke"),
		index.Fields("science_stroke"),
		index.Fields("regular"),
		index.Fields("common_level"),
		index.Fields("nameable"),
		index.Fields("wu_xing", "science_stroke"),
		index.Fields("wu_xing", "kangxi_stroke"),
		index.Fields("regular", "nameable"),
		index.Fields("is_simplified", "is_traditional", "is_kangxi"),
	}
}

func (Character) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "character"},
	}
}
