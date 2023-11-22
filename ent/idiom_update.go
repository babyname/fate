// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"
	"github.com/babyname/fate/ent/idiom"
	"github.com/babyname/fate/ent/predicate"
)

// IdiomUpdate is the builder for updating Idiom entities.
type IdiomUpdate struct {
	config
	hooks    []Hook
	mutation *IdiomMutation
}

// Where appends a list predicates to the IdiomUpdate builder.
func (iu *IdiomUpdate) Where(ps ...predicate.Idiom) *IdiomUpdate {
	iu.mutation.Where(ps...)
	return iu
}

// SetPinYin sets the "pin_yin" field.
func (iu *IdiomUpdate) SetPinYin(s []string) *IdiomUpdate {
	iu.mutation.SetPinYin(s)
	return iu
}

// AppendPinYin appends s to the "pin_yin" field.
func (iu *IdiomUpdate) AppendPinYin(s []string) *IdiomUpdate {
	iu.mutation.AppendPinYin(s)
	return iu
}

// SetWord sets the "word" field.
func (iu *IdiomUpdate) SetWord(s string) *IdiomUpdate {
	iu.mutation.SetWord(s)
	return iu
}

// SetDerivation sets the "derivation" field.
func (iu *IdiomUpdate) SetDerivation(i int) *IdiomUpdate {
	iu.mutation.ResetDerivation()
	iu.mutation.SetDerivation(i)
	return iu
}

// AddDerivation adds i to the "derivation" field.
func (iu *IdiomUpdate) AddDerivation(i int) *IdiomUpdate {
	iu.mutation.AddDerivation(i)
	return iu
}

// SetExplanation sets the "explanation" field.
func (iu *IdiomUpdate) SetExplanation(s string) *IdiomUpdate {
	iu.mutation.SetExplanation(s)
	return iu
}

// SetAbbreviation sets the "abbreviation" field.
func (iu *IdiomUpdate) SetAbbreviation(i int) *IdiomUpdate {
	iu.mutation.ResetAbbreviation()
	iu.mutation.SetAbbreviation(i)
	return iu
}

// AddAbbreviation adds i to the "abbreviation" field.
func (iu *IdiomUpdate) AddAbbreviation(i int) *IdiomUpdate {
	iu.mutation.AddAbbreviation(i)
	return iu
}

// SetExample sets the "example" field.
func (iu *IdiomUpdate) SetExample(s string) *IdiomUpdate {
	iu.mutation.SetExample(s)
	return iu
}

// SetComment sets the "comment" field.
func (iu *IdiomUpdate) SetComment(s string) *IdiomUpdate {
	iu.mutation.SetComment(s)
	return iu
}

// Mutation returns the IdiomMutation object of the builder.
func (iu *IdiomUpdate) Mutation() *IdiomMutation {
	return iu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (iu *IdiomUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, iu.sqlSave, iu.mutation, iu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (iu *IdiomUpdate) SaveX(ctx context.Context) int {
	affected, err := iu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (iu *IdiomUpdate) Exec(ctx context.Context) error {
	_, err := iu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iu *IdiomUpdate) ExecX(ctx context.Context) {
	if err := iu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (iu *IdiomUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(idiom.Table, idiom.Columns, sqlgraph.NewFieldSpec(idiom.FieldID, field.TypeInt32))
	if ps := iu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iu.mutation.PinYin(); ok {
		_spec.SetField(idiom.FieldPinYin, field.TypeJSON, value)
	}
	if value, ok := iu.mutation.AppendedPinYin(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, idiom.FieldPinYin, value)
		})
	}
	if value, ok := iu.mutation.Word(); ok {
		_spec.SetField(idiom.FieldWord, field.TypeString, value)
	}
	if value, ok := iu.mutation.Derivation(); ok {
		_spec.SetField(idiom.FieldDerivation, field.TypeInt, value)
	}
	if value, ok := iu.mutation.AddedDerivation(); ok {
		_spec.AddField(idiom.FieldDerivation, field.TypeInt, value)
	}
	if value, ok := iu.mutation.Explanation(); ok {
		_spec.SetField(idiom.FieldExplanation, field.TypeString, value)
	}
	if value, ok := iu.mutation.Abbreviation(); ok {
		_spec.SetField(idiom.FieldAbbreviation, field.TypeInt, value)
	}
	if value, ok := iu.mutation.AddedAbbreviation(); ok {
		_spec.AddField(idiom.FieldAbbreviation, field.TypeInt, value)
	}
	if value, ok := iu.mutation.Example(); ok {
		_spec.SetField(idiom.FieldExample, field.TypeString, value)
	}
	if value, ok := iu.mutation.Comment(); ok {
		_spec.SetField(idiom.FieldComment, field.TypeString, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, iu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{idiom.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	iu.mutation.done = true
	return n, nil
}

// IdiomUpdateOne is the builder for updating a single Idiom entity.
type IdiomUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *IdiomMutation
}

// SetPinYin sets the "pin_yin" field.
func (iuo *IdiomUpdateOne) SetPinYin(s []string) *IdiomUpdateOne {
	iuo.mutation.SetPinYin(s)
	return iuo
}

// AppendPinYin appends s to the "pin_yin" field.
func (iuo *IdiomUpdateOne) AppendPinYin(s []string) *IdiomUpdateOne {
	iuo.mutation.AppendPinYin(s)
	return iuo
}

// SetWord sets the "word" field.
func (iuo *IdiomUpdateOne) SetWord(s string) *IdiomUpdateOne {
	iuo.mutation.SetWord(s)
	return iuo
}

// SetDerivation sets the "derivation" field.
func (iuo *IdiomUpdateOne) SetDerivation(i int) *IdiomUpdateOne {
	iuo.mutation.ResetDerivation()
	iuo.mutation.SetDerivation(i)
	return iuo
}

// AddDerivation adds i to the "derivation" field.
func (iuo *IdiomUpdateOne) AddDerivation(i int) *IdiomUpdateOne {
	iuo.mutation.AddDerivation(i)
	return iuo
}

// SetExplanation sets the "explanation" field.
func (iuo *IdiomUpdateOne) SetExplanation(s string) *IdiomUpdateOne {
	iuo.mutation.SetExplanation(s)
	return iuo
}

// SetAbbreviation sets the "abbreviation" field.
func (iuo *IdiomUpdateOne) SetAbbreviation(i int) *IdiomUpdateOne {
	iuo.mutation.ResetAbbreviation()
	iuo.mutation.SetAbbreviation(i)
	return iuo
}

// AddAbbreviation adds i to the "abbreviation" field.
func (iuo *IdiomUpdateOne) AddAbbreviation(i int) *IdiomUpdateOne {
	iuo.mutation.AddAbbreviation(i)
	return iuo
}

// SetExample sets the "example" field.
func (iuo *IdiomUpdateOne) SetExample(s string) *IdiomUpdateOne {
	iuo.mutation.SetExample(s)
	return iuo
}

// SetComment sets the "comment" field.
func (iuo *IdiomUpdateOne) SetComment(s string) *IdiomUpdateOne {
	iuo.mutation.SetComment(s)
	return iuo
}

// Mutation returns the IdiomMutation object of the builder.
func (iuo *IdiomUpdateOne) Mutation() *IdiomMutation {
	return iuo.mutation
}

// Where appends a list predicates to the IdiomUpdate builder.
func (iuo *IdiomUpdateOne) Where(ps ...predicate.Idiom) *IdiomUpdateOne {
	iuo.mutation.Where(ps...)
	return iuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (iuo *IdiomUpdateOne) Select(field string, fields ...string) *IdiomUpdateOne {
	iuo.fields = append([]string{field}, fields...)
	return iuo
}

// Save executes the query and returns the updated Idiom entity.
func (iuo *IdiomUpdateOne) Save(ctx context.Context) (*Idiom, error) {
	return withHooks(ctx, iuo.sqlSave, iuo.mutation, iuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (iuo *IdiomUpdateOne) SaveX(ctx context.Context) *Idiom {
	node, err := iuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (iuo *IdiomUpdateOne) Exec(ctx context.Context) error {
	_, err := iuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iuo *IdiomUpdateOne) ExecX(ctx context.Context) {
	if err := iuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (iuo *IdiomUpdateOne) sqlSave(ctx context.Context) (_node *Idiom, err error) {
	_spec := sqlgraph.NewUpdateSpec(idiom.Table, idiom.Columns, sqlgraph.NewFieldSpec(idiom.FieldID, field.TypeInt32))
	id, ok := iuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Idiom.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := iuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, idiom.FieldID)
		for _, f := range fields {
			if !idiom.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != idiom.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := iuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iuo.mutation.PinYin(); ok {
		_spec.SetField(idiom.FieldPinYin, field.TypeJSON, value)
	}
	if value, ok := iuo.mutation.AppendedPinYin(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, idiom.FieldPinYin, value)
		})
	}
	if value, ok := iuo.mutation.Word(); ok {
		_spec.SetField(idiom.FieldWord, field.TypeString, value)
	}
	if value, ok := iuo.mutation.Derivation(); ok {
		_spec.SetField(idiom.FieldDerivation, field.TypeInt, value)
	}
	if value, ok := iuo.mutation.AddedDerivation(); ok {
		_spec.AddField(idiom.FieldDerivation, field.TypeInt, value)
	}
	if value, ok := iuo.mutation.Explanation(); ok {
		_spec.SetField(idiom.FieldExplanation, field.TypeString, value)
	}
	if value, ok := iuo.mutation.Abbreviation(); ok {
		_spec.SetField(idiom.FieldAbbreviation, field.TypeInt, value)
	}
	if value, ok := iuo.mutation.AddedAbbreviation(); ok {
		_spec.AddField(idiom.FieldAbbreviation, field.TypeInt, value)
	}
	if value, ok := iuo.mutation.Example(); ok {
		_spec.SetField(idiom.FieldExample, field.TypeString, value)
	}
	if value, ok := iuo.mutation.Comment(); ok {
		_spec.SetField(idiom.FieldComment, field.TypeString, value)
	}
	_node = &Idiom{config: iuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, iuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{idiom.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	iuo.mutation.done = true
	return _node, nil
}
