// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/babyname/fate/ent/idiom"
)

// Idiom is the model entity for the Idiom schema.
type Idiom struct {
	config `json:"-"`
	// ID of the ent.
	ID int32 `json:"id,omitempty"`
	// PinYin holds the value of the "pin_yin" field.
	PinYin []string `json:"pin_yin,omitempty"`
	// Word holds the value of the "word" field.
	Word string `json:"word,omitempty"`
	// Derivation holds the value of the "derivation" field.
	Derivation int `json:"derivation,omitempty"`
	// Explanation holds the value of the "explanation" field.
	Explanation string `json:"explanation,omitempty"`
	// Abbreviation holds the value of the "abbreviation" field.
	Abbreviation int `json:"abbreviation,omitempty"`
	// Example holds the value of the "example" field.
	Example string `json:"example,omitempty"`
	// Comment holds the value of the "comment" field.
	Comment string `json:"comment,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Idiom) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case idiom.FieldPinYin:
			values[i] = new([]byte)
		case idiom.FieldID, idiom.FieldDerivation, idiom.FieldAbbreviation:
			values[i] = new(sql.NullInt64)
		case idiom.FieldWord, idiom.FieldExplanation, idiom.FieldExample, idiom.FieldComment:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Idiom", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Idiom fields.
func (i *Idiom) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for j := range columns {
		switch columns[j] {
		case idiom.FieldID:
			value, ok := values[j].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			i.ID = int32(value.Int64)
		case idiom.FieldPinYin:
			if value, ok := values[j].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field pin_yin", values[j])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &i.PinYin); err != nil {
					return fmt.Errorf("unmarshal field pin_yin: %w", err)
				}
			}
		case idiom.FieldWord:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field word", values[j])
			} else if value.Valid {
				i.Word = value.String
			}
		case idiom.FieldDerivation:
			if value, ok := values[j].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field derivation", values[j])
			} else if value.Valid {
				i.Derivation = int(value.Int64)
			}
		case idiom.FieldExplanation:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field explanation", values[j])
			} else if value.Valid {
				i.Explanation = value.String
			}
		case idiom.FieldAbbreviation:
			if value, ok := values[j].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field abbreviation", values[j])
			} else if value.Valid {
				i.Abbreviation = int(value.Int64)
			}
		case idiom.FieldExample:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field example", values[j])
			} else if value.Valid {
				i.Example = value.String
			}
		case idiom.FieldComment:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field comment", values[j])
			} else if value.Valid {
				i.Comment = value.String
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Idiom.
// Note that you need to call Idiom.Unwrap() before calling this method if this Idiom
// was returned from a transaction, and the transaction was committed or rolled back.
func (i *Idiom) Update() *IdiomUpdateOne {
	return NewIdiomClient(i.config).UpdateOne(i)
}

// Unwrap unwraps the Idiom entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (i *Idiom) Unwrap() *Idiom {
	_tx, ok := i.config.driver.(*txDriver)
	if !ok {
		panic("ent: Idiom is not a transactional entity")
	}
	i.config.driver = _tx.drv
	return i
}

// String implements the fmt.Stringer.
func (i *Idiom) String() string {
	var builder strings.Builder
	builder.WriteString("Idiom(")
	builder.WriteString(fmt.Sprintf("id=%v, ", i.ID))
	builder.WriteString("pin_yin=")
	builder.WriteString(fmt.Sprintf("%v", i.PinYin))
	builder.WriteString(", ")
	builder.WriteString("word=")
	builder.WriteString(i.Word)
	builder.WriteString(", ")
	builder.WriteString("derivation=")
	builder.WriteString(fmt.Sprintf("%v", i.Derivation))
	builder.WriteString(", ")
	builder.WriteString("explanation=")
	builder.WriteString(i.Explanation)
	builder.WriteString(", ")
	builder.WriteString("abbreviation=")
	builder.WriteString(fmt.Sprintf("%v", i.Abbreviation))
	builder.WriteString(", ")
	builder.WriteString("example=")
	builder.WriteString(i.Example)
	builder.WriteString(", ")
	builder.WriteString("comment=")
	builder.WriteString(i.Comment)
	builder.WriteByte(')')
	return builder.String()
}

// Idioms is a parsable slice of Idiom.
type Idioms []*Idiom
