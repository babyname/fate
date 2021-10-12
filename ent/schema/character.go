package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Character holds the schema definition for the Character entity.
type Character struct {
	ent.Schema
}

const (
	WuXingGold  = "金"
	WuXingWood  = "木"
	WuXingWater = "水"
	WuXingFire  = "火"
	WuXingSoil  = "土"
	WuXingNone  = "无"
)

//const "大凶"
//const "凶"
//const "凶多于吉"
//const "吉凶参半"
//const "吉多于凶"
//const "吉"
//const "大吉"

// Fields of the Character.
func (Character) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		//拼音
		field.Strings("pin_yin").Optional(),
		//字符
		field.String("ch"),
		//科学笔画
		field.Int("science_stroke"),
		//部首
		field.String("radical"),
		//部首笔画
		field.Int("radical_stroke"),
		//总笔画数
		field.Int("stroke"),
		//是否康熙字典
		field.Bool("is_kangxi"),
		//康熙
		field.String("kangxi"),
		//康熙笔画
		field.Int("kangxi_stroke"),
		//简体部首
		field.String("simple_radical"),
		//简体部首笔画
		field.Int("simple_radical_stroke"),
		//简体笔画
		field.Int("simple_total_stroke"),
		//繁体部首
		field.String("traditional_radical"),
		//繁体部首笔画
		field.Int("traditional_radical_stroke"),
		//繁体部首笔画
		field.Int("traditional_total_stroke"),
		//姓名学
		field.Bool("is_name_science"),
		//五行
		field.Enum("wu_xing").Values(
			WuXingGold,
			WuXingWood,
			WuXingWater,
			WuXingFire,
			WuXingSoil,
			WuXingNone).
			Default(WuXingNone),
		//吉凶寓意
		field.String("lucky"),
		//常用
		field.Bool("is_regular"),
		//繁体字
		field.Strings("traditional_character").Optional(),
		//异体字
		field.Strings("variant_character").Optional(),
		//说明
		field.Strings("comment"),
	}
}

// Edges of the Character.
func (Character) Edges() []ent.Edge {
	return nil
}
