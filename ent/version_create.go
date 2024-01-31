// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/babyname/fate/ent/version"
)

// VersionCreate is the builder for creating a Version entity.
type VersionCreate struct {
	config
	mutation *VersionMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCurrentVersion sets the "current_version" field.
func (vc *VersionCreate) SetCurrentVersion(i int) *VersionCreate {
	vc.mutation.SetCurrentVersion(i)
	return vc
}

// SetUpdatedUnix sets the "updated_unix" field.
func (vc *VersionCreate) SetUpdatedUnix(i int) *VersionCreate {
	vc.mutation.SetUpdatedUnix(i)
	return vc
}

// Mutation returns the VersionMutation object of the builder.
func (vc *VersionCreate) Mutation() *VersionMutation {
	return vc.mutation
}

// Save creates the Version in the database.
func (vc *VersionCreate) Save(ctx context.Context) (*Version, error) {
	return withHooks[*Version, VersionMutation](ctx, vc.sqlSave, vc.mutation, vc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (vc *VersionCreate) SaveX(ctx context.Context) *Version {
	v, err := vc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (vc *VersionCreate) Exec(ctx context.Context) error {
	_, err := vc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vc *VersionCreate) ExecX(ctx context.Context) {
	if err := vc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (vc *VersionCreate) check() error {
	if _, ok := vc.mutation.CurrentVersion(); !ok {
		return &ValidationError{Name: "current_version", err: errors.New(`ent: missing required field "Version.current_version"`)}
	}
	if _, ok := vc.mutation.UpdatedUnix(); !ok {
		return &ValidationError{Name: "updated_unix", err: errors.New(`ent: missing required field "Version.updated_unix"`)}
	}
	return nil
}

func (vc *VersionCreate) sqlSave(ctx context.Context) (*Version, error) {
	if err := vc.check(); err != nil {
		return nil, err
	}
	_node, _spec := vc.createSpec()
	if err := sqlgraph.CreateNode(ctx, vc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	vc.mutation.id = &_node.ID
	vc.mutation.done = true
	return _node, nil
}

func (vc *VersionCreate) createSpec() (*Version, *sqlgraph.CreateSpec) {
	var (
		_node = &Version{config: vc.config}
		_spec = sqlgraph.NewCreateSpec(version.Table, sqlgraph.NewFieldSpec(version.FieldID, field.TypeInt))
	)
	_spec.OnConflict = vc.conflict
	if value, ok := vc.mutation.CurrentVersion(); ok {
		_spec.SetField(version.FieldCurrentVersion, field.TypeInt, value)
		_node.CurrentVersion = value
	}
	if value, ok := vc.mutation.UpdatedUnix(); ok {
		_spec.SetField(version.FieldUpdatedUnix, field.TypeInt, value)
		_node.UpdatedUnix = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Version.Create().
//		SetCurrentVersion(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.VersionUpsert) {
//			SetCurrentVersion(v+v).
//		}).
//		Exec(ctx)
func (vc *VersionCreate) OnConflict(opts ...sql.ConflictOption) *VersionUpsertOne {
	vc.conflict = opts
	return &VersionUpsertOne{
		create: vc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Version.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (vc *VersionCreate) OnConflictColumns(columns ...string) *VersionUpsertOne {
	vc.conflict = append(vc.conflict, sql.ConflictColumns(columns...))
	return &VersionUpsertOne{
		create: vc,
	}
}

type (
	// VersionUpsertOne is the builder for "upsert"-ing
	//  one Version node.
	VersionUpsertOne struct {
		create *VersionCreate
	}

	// VersionUpsert is the "OnConflict" setter.
	VersionUpsert struct {
		*sql.UpdateSet
	}
)

// SetCurrentVersion sets the "current_version" field.
func (u *VersionUpsert) SetCurrentVersion(v int) *VersionUpsert {
	u.Set(version.FieldCurrentVersion, v)
	return u
}

// UpdateCurrentVersion sets the "current_version" field to the value that was provided on create.
func (u *VersionUpsert) UpdateCurrentVersion() *VersionUpsert {
	u.SetExcluded(version.FieldCurrentVersion)
	return u
}

// AddCurrentVersion adds v to the "current_version" field.
func (u *VersionUpsert) AddCurrentVersion(v int) *VersionUpsert {
	u.Add(version.FieldCurrentVersion, v)
	return u
}

// SetUpdatedUnix sets the "updated_unix" field.
func (u *VersionUpsert) SetUpdatedUnix(v int) *VersionUpsert {
	u.Set(version.FieldUpdatedUnix, v)
	return u
}

// UpdateUpdatedUnix sets the "updated_unix" field to the value that was provided on create.
func (u *VersionUpsert) UpdateUpdatedUnix() *VersionUpsert {
	u.SetExcluded(version.FieldUpdatedUnix)
	return u
}

// AddUpdatedUnix adds v to the "updated_unix" field.
func (u *VersionUpsert) AddUpdatedUnix(v int) *VersionUpsert {
	u.Add(version.FieldUpdatedUnix, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Version.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *VersionUpsertOne) UpdateNewValues() *VersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Version.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *VersionUpsertOne) Ignore() *VersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *VersionUpsertOne) DoNothing() *VersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the VersionCreate.OnConflict
// documentation for more info.
func (u *VersionUpsertOne) Update(set func(*VersionUpsert)) *VersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&VersionUpsert{UpdateSet: update})
	}))
	return u
}

// SetCurrentVersion sets the "current_version" field.
func (u *VersionUpsertOne) SetCurrentVersion(v int) *VersionUpsertOne {
	return u.Update(func(s *VersionUpsert) {
		s.SetCurrentVersion(v)
	})
}

// AddCurrentVersion adds v to the "current_version" field.
func (u *VersionUpsertOne) AddCurrentVersion(v int) *VersionUpsertOne {
	return u.Update(func(s *VersionUpsert) {
		s.AddCurrentVersion(v)
	})
}

// UpdateCurrentVersion sets the "current_version" field to the value that was provided on create.
func (u *VersionUpsertOne) UpdateCurrentVersion() *VersionUpsertOne {
	return u.Update(func(s *VersionUpsert) {
		s.UpdateCurrentVersion()
	})
}

// SetUpdatedUnix sets the "updated_unix" field.
func (u *VersionUpsertOne) SetUpdatedUnix(v int) *VersionUpsertOne {
	return u.Update(func(s *VersionUpsert) {
		s.SetUpdatedUnix(v)
	})
}

// AddUpdatedUnix adds v to the "updated_unix" field.
func (u *VersionUpsertOne) AddUpdatedUnix(v int) *VersionUpsertOne {
	return u.Update(func(s *VersionUpsert) {
		s.AddUpdatedUnix(v)
	})
}

// UpdateUpdatedUnix sets the "updated_unix" field to the value that was provided on create.
func (u *VersionUpsertOne) UpdateUpdatedUnix() *VersionUpsertOne {
	return u.Update(func(s *VersionUpsert) {
		s.UpdateUpdatedUnix()
	})
}

// Exec executes the query.
func (u *VersionUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for VersionCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *VersionUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *VersionUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *VersionUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// VersionCreateBulk is the builder for creating many Version entities in bulk.
type VersionCreateBulk struct {
	config
	builders []*VersionCreate
	conflict []sql.ConflictOption
}

// Save creates the Version entities in the database.
func (vcb *VersionCreateBulk) Save(ctx context.Context) ([]*Version, error) {
	specs := make([]*sqlgraph.CreateSpec, len(vcb.builders))
	nodes := make([]*Version, len(vcb.builders))
	mutators := make([]Mutator, len(vcb.builders))
	for i := range vcb.builders {
		func(i int, root context.Context) {
			builder := vcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*VersionMutation)
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
					_, err = mutators[i+1].Mutate(root, vcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = vcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, vcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
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
		if _, err := mutators[0].Mutate(ctx, vcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (vcb *VersionCreateBulk) SaveX(ctx context.Context) []*Version {
	v, err := vcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (vcb *VersionCreateBulk) Exec(ctx context.Context) error {
	_, err := vcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vcb *VersionCreateBulk) ExecX(ctx context.Context) {
	if err := vcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Version.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.VersionUpsert) {
//			SetCurrentVersion(v+v).
//		}).
//		Exec(ctx)
func (vcb *VersionCreateBulk) OnConflict(opts ...sql.ConflictOption) *VersionUpsertBulk {
	vcb.conflict = opts
	return &VersionUpsertBulk{
		create: vcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Version.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (vcb *VersionCreateBulk) OnConflictColumns(columns ...string) *VersionUpsertBulk {
	vcb.conflict = append(vcb.conflict, sql.ConflictColumns(columns...))
	return &VersionUpsertBulk{
		create: vcb,
	}
}

// VersionUpsertBulk is the builder for "upsert"-ing
// a bulk of Version nodes.
type VersionUpsertBulk struct {
	create *VersionCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Version.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *VersionUpsertBulk) UpdateNewValues() *VersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Version.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *VersionUpsertBulk) Ignore() *VersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *VersionUpsertBulk) DoNothing() *VersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the VersionCreateBulk.OnConflict
// documentation for more info.
func (u *VersionUpsertBulk) Update(set func(*VersionUpsert)) *VersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&VersionUpsert{UpdateSet: update})
	}))
	return u
}

// SetCurrentVersion sets the "current_version" field.
func (u *VersionUpsertBulk) SetCurrentVersion(v int) *VersionUpsertBulk {
	return u.Update(func(s *VersionUpsert) {
		s.SetCurrentVersion(v)
	})
}

// AddCurrentVersion adds v to the "current_version" field.
func (u *VersionUpsertBulk) AddCurrentVersion(v int) *VersionUpsertBulk {
	return u.Update(func(s *VersionUpsert) {
		s.AddCurrentVersion(v)
	})
}

// UpdateCurrentVersion sets the "current_version" field to the value that was provided on create.
func (u *VersionUpsertBulk) UpdateCurrentVersion() *VersionUpsertBulk {
	return u.Update(func(s *VersionUpsert) {
		s.UpdateCurrentVersion()
	})
}

// SetUpdatedUnix sets the "updated_unix" field.
func (u *VersionUpsertBulk) SetUpdatedUnix(v int) *VersionUpsertBulk {
	return u.Update(func(s *VersionUpsert) {
		s.SetUpdatedUnix(v)
	})
}

// AddUpdatedUnix adds v to the "updated_unix" field.
func (u *VersionUpsertBulk) AddUpdatedUnix(v int) *VersionUpsertBulk {
	return u.Update(func(s *VersionUpsert) {
		s.AddUpdatedUnix(v)
	})
}

// UpdateUpdatedUnix sets the "updated_unix" field to the value that was provided on create.
func (u *VersionUpsertBulk) UpdateUpdatedUnix() *VersionUpsertBulk {
	return u.Update(func(s *VersionUpsert) {
		s.UpdateUpdatedUnix()
	})
}

// Exec executes the query.
func (u *VersionUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the VersionCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for VersionCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *VersionUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
