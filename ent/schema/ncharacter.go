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
		field.Int("id"),                               //char rune code
		field.Strings("pin_yin"),                      //pinyin
		field.String("ch"),                            //char
		field.Int("ch_stroke"),                        //char stroke
		field.Int("ch_type").Default(CharTypeUnknown), //char type
		field.String("radical"),                       //radical
		field.Int("radical_stroke"),                   //radical stroke
		field.Ints("relate_simplified"),               // relate simple chinese
		field.Ints("relate_kang_xi"),                  //relate kang xi chinese
		field.Ints("relate_traditional"),              //relate traditional chinese
		field.Ints("relate_variant"),                  //relate other variant characters
		field.Bool("is_name_science"),                 //is name science
		field.Int("name_science_ch_stroke"),           //name science stroke
		field.Bool("is_regular"),                      //is regular
		field.String("wu_xing"),                       //wu xing
		field.String("lucky"),                         //lucky
		field.String("explanation"),
		field.String("comment"),
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
