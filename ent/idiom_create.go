// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/babyname/fate/ent/idiom"
)

// IdiomCreate is the builder for creating a Idiom entity.
type IdiomCreate struct {
	config
	mutation *IdiomMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetPinYin sets the "pin_yin" field.
func (ic *IdiomCreate) SetPinYin(s []string) *IdiomCreate {
	ic.mutation.SetPinYin(s)
	return ic
}

// SetWord sets the "word" field.
func (ic *IdiomCreate) SetWord(s string) *IdiomCreate {
	ic.mutation.SetWord(s)
	return ic
}

// SetDerivation sets the "derivation" field.
func (ic *IdiomCreate) SetDerivation(i int) *IdiomCreate {
	ic.mutation.SetDerivation(i)
	return ic
}

// SetExplanation sets the "explanation" field.
func (ic *IdiomCreate) SetExplanation(s string) *IdiomCreate {
	ic.mutation.SetExplanation(s)
	return ic
}

// SetAbbreviation sets the "abbreviation" field.
func (ic *IdiomCreate) SetAbbreviation(i int) *IdiomCreate {
	ic.mutation.SetAbbreviation(i)
	return ic
}

// SetExample sets the "example" field.
func (ic *IdiomCreate) SetExample(s string) *IdiomCreate {
	ic.mutation.SetExample(s)
	return ic
}

// SetComment sets the "comment" field.
func (ic *IdiomCreate) SetComment(s string) *IdiomCreate {
	ic.mutation.SetComment(s)
	return ic
}

// SetID sets the "id" field.
func (ic *IdiomCreate) SetID(i int32) *IdiomCreate {
	ic.mutation.SetID(i)
	return ic
}

// Mutation returns the IdiomMutation object of the builder.
func (ic *IdiomCreate) Mutation() *IdiomMutation {
	return ic.mutation
}

// Save creates the Idiom in the database.
func (ic *IdiomCreate) Save(ctx context.Context) (*Idiom, error) {
	return withHooks[*Idiom, IdiomMutation](ctx, ic.sqlSave, ic.mutation, ic.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ic *IdiomCreate) SaveX(ctx context.Context) *Idiom {
	v, err := ic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ic *IdiomCreate) Exec(ctx context.Context) error {
	_, err := ic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ic *IdiomCreate) ExecX(ctx context.Context) {
	if err := ic.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ic *IdiomCreate) check() error {
	if _, ok := ic.mutation.PinYin(); !ok {
		return &ValidationError{Name: "pin_yin", err: errors.New(`ent: missing required field "Idiom.pin_yin"`)}
	}
	if _, ok := ic.mutation.Word(); !ok {
		return &ValidationError{Name: "word", err: errors.New(`ent: missing required field "Idiom.word"`)}
	}
	if _, ok := ic.mutation.Derivation(); !ok {
		return &ValidationError{Name: "derivation", err: errors.New(`ent: missing required field "Idiom.derivation"`)}
	}
	if _, ok := ic.mutation.Explanation(); !ok {
		return &ValidationError{Name: "explanation", err: errors.New(`ent: missing required field "Idiom.explanation"`)}
	}
	if _, ok := ic.mutation.Abbreviation(); !ok {
		return &ValidationError{Name: "abbreviation", err: errors.New(`ent: missing required field "Idiom.abbreviation"`)}
	}
	if _, ok := ic.mutation.Example(); !ok {
		return &ValidationError{Name: "example", err: errors.New(`ent: missing required field "Idiom.example"`)}
	}
	if _, ok := ic.mutation.Comment(); !ok {
		return &ValidationError{Name: "comment", err: errors.New(`ent: missing required field "Idiom.comment"`)}
	}
	return nil
}

func (ic *IdiomCreate) sqlSave(ctx context.Context) (*Idiom, error) {
	if err := ic.check(); err != nil {
		return nil, err
	}
	_node, _spec := ic.createSpec()
	if err := sqlgraph.CreateNode(ctx, ic.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int32(id)
	}
	ic.mutation.id = &_node.ID
	ic.mutation.done = true
	return _node, nil
}

func (ic *IdiomCreate) createSpec() (*Idiom, *sqlgraph.CreateSpec) {
	var (
		_node = &Idiom{config: ic.config}
		_spec = sqlgraph.NewCreateSpec(idiom.Table, sqlgraph.NewFieldSpec(idiom.FieldID, field.TypeInt32))
	)
	_spec.OnConflict = ic.conflict
	if id, ok := ic.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := ic.mutation.PinYin(); ok {
		_spec.SetField(idiom.FieldPinYin, field.TypeJSON, value)
		_node.PinYin = value
	}
	if value, ok := ic.mutation.Word(); ok {
		_spec.SetField(idiom.FieldWord, field.TypeString, value)
		_node.Word = value
	}
	if value, ok := ic.mutation.Derivation(); ok {
		_spec.SetField(idiom.FieldDerivation, field.TypeInt, value)
		_node.Derivation = value
	}
	if value, ok := ic.mutation.Explanation(); ok {
		_spec.SetField(idiom.FieldExplanation, field.TypeString, value)
		_node.Explanation = value
	}
	if value, ok := ic.mutation.Abbreviation(); ok {
		_spec.SetField(idiom.FieldAbbreviation, field.TypeInt, value)
		_node.Abbreviation = value
	}
	if value, ok := ic.mutation.Example(); ok {
		_spec.SetField(idiom.FieldExample, field.TypeString, value)
		_node.Example = value
	}
	if value, ok := ic.mutation.Comment(); ok {
		_spec.SetField(idiom.FieldComment, field.TypeString, value)
		_node.Comment = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Idiom.Create().
//		SetPinYin(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.IdiomUpsert) {
//			SetPinYin(v+v).
//		}).
//		Exec(ctx)
func (ic *IdiomCreate) OnConflict(opts ...sql.ConflictOption) *IdiomUpsertOne {
	ic.conflict = opts
	return &IdiomUpsertOne{
		create: ic,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Idiom.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ic *IdiomCreate) OnConflictColumns(columns ...string) *IdiomUpsertOne {
	ic.conflict = append(ic.conflict, sql.ConflictColumns(columns...))
	return &IdiomUpsertOne{
		create: ic,
	}
}

type (
	// IdiomUpsertOne is the builder for "upsert"-ing
	//  one Idiom node.
	IdiomUpsertOne struct {
		create *IdiomCreate
	}

	// IdiomUpsert is the "OnConflict" setter.
	IdiomUpsert struct {
		*sql.UpdateSet
	}
)

// SetPinYin sets the "pin_yin" field.
func (u *IdiomUpsert) SetPinYin(v []string) *IdiomUpsert {
	u.Set(idiom.FieldPinYin, v)
	return u
}

// UpdatePinYin sets the "pin_yin" field to the value that was provided on create.
func (u *IdiomUpsert) UpdatePinYin() *IdiomUpsert {
	u.SetExcluded(idiom.FieldPinYin)
	return u
}

// SetWord sets the "word" field.
func (u *IdiomUpsert) SetWord(v string) *IdiomUpsert {
	u.Set(idiom.FieldWord, v)
	return u
}

// UpdateWord sets the "word" field to the value that was provided on create.
func (u *IdiomUpsert) UpdateWord() *IdiomUpsert {
	u.SetExcluded(idiom.FieldWord)
	return u
}

// SetDerivation sets the "derivation" field.
func (u *IdiomUpsert) SetDerivation(v int) *IdiomUpsert {
	u.Set(idiom.FieldDerivation, v)
	return u
}

// UpdateDerivation sets the "derivation" field to the value that was provided on create.
func (u *IdiomUpsert) UpdateDerivation() *IdiomUpsert {
	u.SetExcluded(idiom.FieldDerivation)
	return u
}

// AddDerivation adds v to the "derivation" field.
func (u *IdiomUpsert) AddDerivation(v int) *IdiomUpsert {
	u.Add(idiom.FieldDerivation, v)
	return u
}

// SetExplanation sets the "explanation" field.
func (u *IdiomUpsert) SetExplanation(v string) *IdiomUpsert {
	u.Set(idiom.FieldExplanation, v)
	return u
}

// UpdateExplanation sets the "explanation" field to the value that was provided on create.
func (u *IdiomUpsert) UpdateExplanation() *IdiomUpsert {
	u.SetExcluded(idiom.FieldExplanation)
	return u
}

// SetAbbreviation sets the "abbreviation" field.
func (u *IdiomUpsert) SetAbbreviation(v int) *IdiomUpsert {
	u.Set(idiom.FieldAbbreviation, v)
	return u
}

// UpdateAbbreviation sets the "abbreviation" field to the value that was provided on create.
func (u *IdiomUpsert) UpdateAbbreviation() *IdiomUpsert {
	u.SetExcluded(idiom.FieldAbbreviation)
	return u
}

// AddAbbreviation adds v to the "abbreviation" field.
func (u *IdiomUpsert) AddAbbreviation(v int) *IdiomUpsert {
	u.Add(idiom.FieldAbbreviation, v)
	return u
}

// SetExample sets the "example" field.
func (u *IdiomUpsert) SetExample(v string) *IdiomUpsert {
	u.Set(idiom.FieldExample, v)
	return u
}

// UpdateExample sets the "example" field to the value that was provided on create.
func (u *IdiomUpsert) UpdateExample() *IdiomUpsert {
	u.SetExcluded(idiom.FieldExample)
	return u
}

// SetComment sets the "comment" field.
func (u *IdiomUpsert) SetComment(v string) *IdiomUpsert {
	u.Set(idiom.FieldComment, v)
	return u
}

// UpdateComment sets the "comment" field to the value that was provided on create.
func (u *IdiomUpsert) UpdateComment() *IdiomUpsert {
	u.SetExcluded(idiom.FieldComment)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Idiom.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(idiom.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *IdiomUpsertOne) UpdateNewValues() *IdiomUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(idiom.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Idiom.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *IdiomUpsertOne) Ignore() *IdiomUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *IdiomUpsertOne) DoNothing() *IdiomUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the IdiomCreate.OnConflict
// documentation for more info.
func (u *IdiomUpsertOne) Update(set func(*IdiomUpsert)) *IdiomUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&IdiomUpsert{UpdateSet: update})
	}))
	return u
}

// SetPinYin sets the "pin_yin" field.
func (u *IdiomUpsertOne) SetPinYin(v []string) *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.SetPinYin(v)
	})
}

// UpdatePinYin sets the "pin_yin" field to the value that was provided on create.
func (u *IdiomUpsertOne) UpdatePinYin() *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdatePinYin()
	})
}

// SetWord sets the "word" field.
func (u *IdiomUpsertOne) SetWord(v string) *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.SetWord(v)
	})
}

// UpdateWord sets the "word" field to the value that was provided on create.
func (u *IdiomUpsertOne) UpdateWord() *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateWord()
	})
}

// SetDerivation sets the "derivation" field.
func (u *IdiomUpsertOne) SetDerivation(v int) *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.SetDerivation(v)
	})
}

// AddDerivation adds v to the "derivation" field.
func (u *IdiomUpsertOne) AddDerivation(v int) *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.AddDerivation(v)
	})
}

// UpdateDerivation sets the "derivation" field to the value that was provided on create.
func (u *IdiomUpsertOne) UpdateDerivation() *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateDerivation()
	})
}

// SetExplanation sets the "explanation" field.
func (u *IdiomUpsertOne) SetExplanation(v string) *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.SetExplanation(v)
	})
}

// UpdateExplanation sets the "explanation" field to the value that was provided on create.
func (u *IdiomUpsertOne) UpdateExplanation() *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateExplanation()
	})
}

// SetAbbreviation sets the "abbreviation" field.
func (u *IdiomUpsertOne) SetAbbreviation(v int) *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.SetAbbreviation(v)
	})
}

// AddAbbreviation adds v to the "abbreviation" field.
func (u *IdiomUpsertOne) AddAbbreviation(v int) *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.AddAbbreviation(v)
	})
}

// UpdateAbbreviation sets the "abbreviation" field to the value that was provided on create.
func (u *IdiomUpsertOne) UpdateAbbreviation() *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateAbbreviation()
	})
}

// SetExample sets the "example" field.
func (u *IdiomUpsertOne) SetExample(v string) *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.SetExample(v)
	})
}

// UpdateExample sets the "example" field to the value that was provided on create.
func (u *IdiomUpsertOne) UpdateExample() *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateExample()
	})
}

// SetComment sets the "comment" field.
func (u *IdiomUpsertOne) SetComment(v string) *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.SetComment(v)
	})
}

// UpdateComment sets the "comment" field to the value that was provided on create.
func (u *IdiomUpsertOne) UpdateComment() *IdiomUpsertOne {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateComment()
	})
}

// Exec executes the query.
func (u *IdiomUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for IdiomCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *IdiomUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *IdiomUpsertOne) ID(ctx context.Context) (id int32, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *IdiomUpsertOne) IDX(ctx context.Context) int32 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// IdiomCreateBulk is the builder for creating many Idiom entities in bulk.
type IdiomCreateBulk struct {
	config
	builders []*IdiomCreate
	conflict []sql.ConflictOption
}

// Save creates the Idiom entities in the database.
func (icb *IdiomCreateBulk) Save(ctx context.Context) ([]*Idiom, error) {
	specs := make([]*sqlgraph.CreateSpec, len(icb.builders))
	nodes := make([]*Idiom, len(icb.builders))
	mutators := make([]Mutator, len(icb.builders))
	for i := range icb.builders {
		func(i int, root context.Context) {
			builder := icb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*IdiomMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, icb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = icb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, icb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int32(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, icb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (icb *IdiomCreateBulk) SaveX(ctx context.Context) []*Idiom {
	v, err := icb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (icb *IdiomCreateBulk) Exec(ctx context.Context) error {
	_, err := icb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (icb *IdiomCreateBulk) ExecX(ctx context.Context) {
	if err := icb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Idiom.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.IdiomUpsert) {
//			SetPinYin(v+v).
//		}).
//		Exec(ctx)
func (icb *IdiomCreateBulk) OnConflict(opts ...sql.ConflictOption) *IdiomUpsertBulk {
	icb.conflict = opts
	return &IdiomUpsertBulk{
		create: icb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Idiom.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (icb *IdiomCreateBulk) OnConflictColumns(columns ...string) *IdiomUpsertBulk {
	icb.conflict = append(icb.conflict, sql.ConflictColumns(columns...))
	return &IdiomUpsertBulk{
		create: icb,
	}
}

// IdiomUpsertBulk is the builder for "upsert"-ing
// a bulk of Idiom nodes.
type IdiomUpsertBulk struct {
	create *IdiomCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Idiom.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(idiom.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *IdiomUpsertBulk) UpdateNewValues() *IdiomUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(idiom.FieldID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Idiom.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *IdiomUpsertBulk) Ignore() *IdiomUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *IdiomUpsertBulk) DoNothing() *IdiomUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the IdiomCreateBulk.OnConflict
// documentation for more info.
func (u *IdiomUpsertBulk) Update(set func(*IdiomUpsert)) *IdiomUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&IdiomUpsert{UpdateSet: update})
	}))
	return u
}

// SetPinYin sets the "pin_yin" field.
func (u *IdiomUpsertBulk) SetPinYin(v []string) *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.SetPinYin(v)
	})
}

// UpdatePinYin sets the "pin_yin" field to the value that was provided on create.
func (u *IdiomUpsertBulk) UpdatePinYin() *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdatePinYin()
	})
}

// SetWord sets the "word" field.
func (u *IdiomUpsertBulk) SetWord(v string) *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.SetWord(v)
	})
}

// UpdateWord sets the "word" field to the value that was provided on create.
func (u *IdiomUpsertBulk) UpdateWord() *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateWord()
	})
}

// SetDerivation sets the "derivation" field.
func (u *IdiomUpsertBulk) SetDerivation(v int) *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.SetDerivation(v)
	})
}

// AddDerivation adds v to the "derivation" field.
func (u *IdiomUpsertBulk) AddDerivation(v int) *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.AddDerivation(v)
	})
}

// UpdateDerivation sets the "derivation" field to the value that was provided on create.
func (u *IdiomUpsertBulk) UpdateDerivation() *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateDerivation()
	})
}

// SetExplanation sets the "explanation" field.
func (u *IdiomUpsertBulk) SetExplanation(v string) *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.SetExplanation(v)
	})
}

// UpdateExplanation sets the "explanation" field to the value that was provided on create.
func (u *IdiomUpsertBulk) UpdateExplanation() *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateExplanation()
	})
}

// SetAbbreviation sets the "abbreviation" field.
func (u *IdiomUpsertBulk) SetAbbreviation(v int) *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.SetAbbreviation(v)
	})
}

// AddAbbreviation adds v to the "abbreviation" field.
func (u *IdiomUpsertBulk) AddAbbreviation(v int) *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.AddAbbreviation(v)
	})
}

// UpdateAbbreviation sets the "abbreviation" field to the value that was provided on create.
func (u *IdiomUpsertBulk) UpdateAbbreviation() *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateAbbreviation()
	})
}

// SetExample sets the "example" field.
func (u *IdiomUpsertBulk) SetExample(v string) *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.SetExample(v)
	})
}

// UpdateExample sets the "example" field to the value that was provided on create.
func (u *IdiomUpsertBulk) UpdateExample() *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateExample()
	})
}

// SetComment sets the "comment" field.
func (u *IdiomUpsertBulk) SetComment(v string) *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.SetComment(v)
	})
}

// UpdateComment sets the "comment" field to the value that was provided on create.
func (u *IdiomUpsertBulk) UpdateComment() *IdiomUpsertBulk {
	return u.Update(func(s *IdiomUpsert) {
		s.UpdateComment()
	})
}

// Exec executes the query.
func (u *IdiomUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the IdiomCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for IdiomCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *IdiomUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
