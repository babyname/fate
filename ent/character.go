// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/godcong/fate/ent/character"
)

// Character is the model entity for the Character schema.
type Character struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// PinYin holds the value of the "pin_yin" field.
	PinYin []string `json:"pin_yin,omitempty"`
	// Ch holds the value of the "ch" field.
	Ch string `json:"ch,omitempty"`
	// ScienceStroke holds the value of the "science_stroke" field.
	ScienceStroke int8 `json:"science_stroke,omitempty"`
	// Radical holds the value of the "radical" field.
	Radical string `json:"radical,omitempty"`
	// RadicalStroke holds the value of the "radical_stroke" field.
	RadicalStroke int8 `json:"radical_stroke,omitempty"`
	// Stroke holds the value of the "stroke" field.
	Stroke int8 `json:"stroke,omitempty"`
	// IsKangxi holds the value of the "is_kangxi" field.
	IsKangxi bool `json:"is_kangxi,omitempty"`
	// Kangxi holds the value of the "kangxi" field.
	Kangxi string `json:"kangxi,omitempty"`
	// KangxiStroke holds the value of the "kangxi_stroke" field.
	KangxiStroke string `json:"kangxi_stroke,omitempty"`
	// SimpleRadical holds the value of the "simple_radical" field.
	SimpleRadical string `json:"simple_radical,omitempty"`
	// SimpleRadicalStroke holds the value of the "simple_radical_stroke" field.
	SimpleRadicalStroke string `json:"simple_radical_stroke,omitempty"`
	// SimpleTotalStroke holds the value of the "simple_total_stroke" field.
	SimpleTotalStroke int8 `json:"simple_total_stroke,omitempty"`
	// TraditionalRadical holds the value of the "traditional_radical" field.
	TraditionalRadical string `json:"traditional_radical,omitempty"`
	// TraditionalRadicalStroke holds the value of the "traditional_radical_stroke" field.
	TraditionalRadicalStroke int8 `json:"traditional_radical_stroke,omitempty"`
	// TraditionalTotalStroke holds the value of the "traditional_total_stroke" field.
	TraditionalTotalStroke int8 `json:"traditional_total_stroke,omitempty"`
	// IsNameScience holds the value of the "is_name_science" field.
	IsNameScience bool `json:"is_name_science,omitempty"`
	// WuXing holds the value of the "wu_xing" field.
	WuXing string `json:"wu_xing,omitempty"`
	// Lucky holds the value of the "lucky" field.
	Lucky string `json:"lucky,omitempty"`
	// IsRegular holds the value of the "is_regular" field.
	IsRegular bool `json:"is_regular,omitempty"`
	// TraditionalCharacter holds the value of the "traditional_character" field.
	TraditionalCharacter []string `json:"traditional_character,omitempty"`
	// VariantCharacter holds the value of the "variant_character" field.
	VariantCharacter []string `json:"variant_character,omitempty"`
	// Comment holds the value of the "comment" field.
	Comment string `json:"comment,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Character) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case character.FieldPinYin, character.FieldTraditionalCharacter, character.FieldVariantCharacter:
			values[i] = new([]byte)
		case character.FieldIsKangxi, character.FieldIsNameScience, character.FieldIsRegular:
			values[i] = new(sql.NullBool)
		case character.FieldScienceStroke, character.FieldRadicalStroke, character.FieldStroke, character.FieldSimpleTotalStroke, character.FieldTraditionalRadicalStroke, character.FieldTraditionalTotalStroke:
			values[i] = new(sql.NullInt64)
		case character.FieldID, character.FieldCh, character.FieldRadical, character.FieldKangxi, character.FieldKangxiStroke, character.FieldSimpleRadical, character.FieldSimpleRadicalStroke, character.FieldTraditionalRadical, character.FieldWuXing, character.FieldLucky, character.FieldComment:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Character", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Character fields.
func (c *Character) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case character.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				c.ID = value.String
			}
		case character.FieldPinYin:

			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field pin_yin", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &c.PinYin); err != nil {
					return fmt.Errorf("unmarshal field pin_yin: %w", err)
				}
			}
		case character.FieldCh:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field ch", values[i])
			} else if value.Valid {
				c.Ch = value.String
			}
		case character.FieldScienceStroke:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field science_stroke", values[i])
			} else if value.Valid {
				c.ScienceStroke = int8(value.Int64)
			}
		case character.FieldRadical:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field radical", values[i])
			} else if value.Valid {
				c.Radical = value.String
			}
		case character.FieldRadicalStroke:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field radical_stroke", values[i])
			} else if value.Valid {
				c.RadicalStroke = int8(value.Int64)
			}
		case character.FieldStroke:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field stroke", values[i])
			} else if value.Valid {
				c.Stroke = int8(value.Int64)
			}
		case character.FieldIsKangxi:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_kangxi", values[i])
			} else if value.Valid {
				c.IsKangxi = value.Bool
			}
		case character.FieldKangxi:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kangxi", values[i])
			} else if value.Valid {
				c.Kangxi = value.String
			}
		case character.FieldKangxiStroke:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kangxi_stroke", values[i])
			} else if value.Valid {
				c.KangxiStroke = value.String
			}
		case character.FieldSimpleRadical:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field simple_radical", values[i])
			} else if value.Valid {
				c.SimpleRadical = value.String
			}
		case character.FieldSimpleRadicalStroke:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field simple_radical_stroke", values[i])
			} else if value.Valid {
				c.SimpleRadicalStroke = value.String
			}
		case character.FieldSimpleTotalStroke:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field simple_total_stroke", values[i])
			} else if value.Valid {
				c.SimpleTotalStroke = int8(value.Int64)
			}
		case character.FieldTraditionalRadical:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field traditional_radical", values[i])
			} else if value.Valid {
				c.TraditionalRadical = value.String
			}
		case character.FieldTraditionalRadicalStroke:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field traditional_radical_stroke", values[i])
			} else if value.Valid {
				c.TraditionalRadicalStroke = int8(value.Int64)
			}
		case character.FieldTraditionalTotalStroke:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field traditional_total_stroke", values[i])
			} else if value.Valid {
				c.TraditionalTotalStroke = int8(value.Int64)
			}
		case character.FieldIsNameScience:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_name_science", values[i])
			} else if value.Valid {
				c.IsNameScience = value.Bool
			}
		case character.FieldWuXing:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field wu_xing", values[i])
			} else if value.Valid {
				c.WuXing = value.String
			}
		case character.FieldLucky:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field lucky", values[i])
			} else if value.Valid {
				c.Lucky = value.String
			}
		case character.FieldIsRegular:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_regular", values[i])
			} else if value.Valid {
				c.IsRegular = value.Bool
			}
		case character.FieldTraditionalCharacter:

			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field traditional_character", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &c.TraditionalCharacter); err != nil {
					return fmt.Errorf("unmarshal field traditional_character: %w", err)
				}
			}
		case character.FieldVariantCharacter:

			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field variant_character", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &c.VariantCharacter); err != nil {
					return fmt.Errorf("unmarshal field variant_character: %w", err)
				}
			}
		case character.FieldComment:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field comment", values[i])
			} else if value.Valid {
				c.Comment = value.String
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Character.
// Note that you need to call Character.Unwrap() before calling this method if this Character
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Character) Update() *CharacterUpdateOne {
	return (&CharacterClient{config: c.config}).UpdateOne(c)
}

// Unwrap unwraps the Character entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Character) Unwrap() *Character {
	tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Character is not a transactional entity")
	}
	c.config.driver = tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Character) String() string {
	var builder strings.Builder
	builder.WriteString("Character(")
	builder.WriteString(fmt.Sprintf("id=%v", c.ID))
	builder.WriteString(", pin_yin=")
	builder.WriteString(fmt.Sprintf("%v", c.PinYin))
	builder.WriteString(", ch=")
	builder.WriteString(c.Ch)
	builder.WriteString(", science_stroke=")
	builder.WriteString(fmt.Sprintf("%v", c.ScienceStroke))
	builder.WriteString(", radical=")
	builder.WriteString(c.Radical)
	builder.WriteString(", radical_stroke=")
	builder.WriteString(fmt.Sprintf("%v", c.RadicalStroke))
	builder.WriteString(", stroke=")
	builder.WriteString(fmt.Sprintf("%v", c.Stroke))
	builder.WriteString(", is_kangxi=")
	builder.WriteString(fmt.Sprintf("%v", c.IsKangxi))
	builder.WriteString(", kangxi=")
	builder.WriteString(c.Kangxi)
	builder.WriteString(", kangxi_stroke=")
	builder.WriteString(c.KangxiStroke)
	builder.WriteString(", simple_radical=")
	builder.WriteString(c.SimpleRadical)
	builder.WriteString(", simple_radical_stroke=")
	builder.WriteString(c.SimpleRadicalStroke)
	builder.WriteString(", simple_total_stroke=")
	builder.WriteString(fmt.Sprintf("%v", c.SimpleTotalStroke))
	builder.WriteString(", traditional_radical=")
	builder.WriteString(c.TraditionalRadical)
	builder.WriteString(", traditional_radical_stroke=")
	builder.WriteString(fmt.Sprintf("%v", c.TraditionalRadicalStroke))
	builder.WriteString(", traditional_total_stroke=")
	builder.WriteString(fmt.Sprintf("%v", c.TraditionalTotalStroke))
	builder.WriteString(", is_name_science=")
	builder.WriteString(fmt.Sprintf("%v", c.IsNameScience))
	builder.WriteString(", wu_xing=")
	builder.WriteString(c.WuXing)
	builder.WriteString(", lucky=")
	builder.WriteString(c.Lucky)
	builder.WriteString(", is_regular=")
	builder.WriteString(fmt.Sprintf("%v", c.IsRegular))
	builder.WriteString(", traditional_character=")
	builder.WriteString(fmt.Sprintf("%v", c.TraditionalCharacter))
	builder.WriteString(", variant_character=")
	builder.WriteString(fmt.Sprintf("%v", c.VariantCharacter))
	builder.WriteString(", comment=")
	builder.WriteString(c.Comment)
	builder.WriteByte(')')
	return builder.String()
}

// Characters is a parsable slice of Character.
type Characters []*Character

func (c Characters) config(cfg config) {
	for _i := range c {
		c[_i].config = cfg
	}
}
