package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// NCharacter holds the schema definition for the NCharacter entity.
type NCharacter struct {
	ent.Schema
}

const (
	CharTypeUnknown     = 0x00
	CharTypeSimple      = 0x01
	CharTypeTraditional = 0x02
	CharTypeKangXi      = 0x04
	CharTypeVariant     = 0x08
)

// Fields of the NCharacter.
func (NCharacter) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Default(0),                    //char rune code
		field.Strings("pin_yin").Default([]string{}),  //pinyin
		field.String("char").Default(""),              //char
		field.Int("char_stroke").Default(0),           //char stroke
		field.String("radical").Default(""),           //radical
		field.Int("radical_stroke").Default(0),        //radical stroke
		field.Bool("is_regular").Default(false),       //is regular
		field.Bool("is_simplified").Default(false),    //relate simple chinese
		field.Ints("simplified_id").Default([]int{}),  //relate simple chinese
		field.Bool("is_traditional").Default(false),   //relate traditional chinese
		field.Ints("traditional_id").Default([]int{}), //relate traditional chinese
		field.Bool("is_kang_xi").Default(false),       //relate kang xi chinese
		field.Ints("kang_xi_id").Default([]int{}),     //relate kang xi chinese
		field.Bool("is_variant").Default(false),       //relate other variant characters
		field.Ints("variant_id").Default([]int{}),     //relate other variant characters
		field.Bool("is_science").Default(false),       //is name science
		field.Int("science_stroke").Default(0),        //name science stroke
		field.String("wu_xing").Default(""),           //wu xing
		field.String("lucky").Default(""),             //lucky
		field.String("explanation").Default(""),
		field.String("comment").Default(""),
		field.Bool("need_fix").Default(false),
	}
}

// Edges of the NCharacter.
func (NCharacter) Edges() []ent.Edge {
	return nil
}

// Annotations of the NCharacter.
func (NCharacter) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "n_character"},
	}
}
