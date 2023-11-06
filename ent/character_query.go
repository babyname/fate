// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/ent/predicate"
)

// CharacterQuery is the builder for querying Character entities.
type CharacterQuery struct {
	config
	ctx        *QueryContext
	order      []character.OrderOption
	inters     []Interceptor
	predicates []predicate.Character
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CharacterQuery builder.
func (cq *CharacterQuery) Where(ps ...predicate.Character) *CharacterQuery {
	cq.predicates = append(cq.predicates, ps...)
	return cq
}

// Limit the number of records to be returned by this query.
func (cq *CharacterQuery) Limit(limit int) *CharacterQuery {
	cq.ctx.Limit = &limit
	return cq
}

// Offset to start from.
func (cq *CharacterQuery) Offset(offset int) *CharacterQuery {
	cq.ctx.Offset = &offset
	return cq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (cq *CharacterQuery) Unique(unique bool) *CharacterQuery {
	cq.ctx.Unique = &unique
	return cq
}

// Order specifies how the records should be ordered.
func (cq *CharacterQuery) Order(o ...character.OrderOption) *CharacterQuery {
	cq.order = append(cq.order, o...)
	return cq
}

// First returns the first Character entity from the query.
// Returns a *NotFoundError when no Character was found.
func (cq *CharacterQuery) First(ctx context.Context) (*Character, error) {
	nodes, err := cq.Limit(1).All(setContextOp(ctx, cq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{character.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (cq *CharacterQuery) FirstX(ctx context.Context) *Character {
	node, err := cq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Character ID from the query.
// Returns a *NotFoundError when no Character ID was found.
func (cq *CharacterQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = cq.Limit(1).IDs(setContextOp(ctx, cq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{character.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (cq *CharacterQuery) FirstIDX(ctx context.Context) string {
	id, err := cq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Character entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Character entity is found.
// Returns a *NotFoundError when no Character entities are found.
func (cq *CharacterQuery) Only(ctx context.Context) (*Character, error) {
	nodes, err := cq.Limit(2).All(setContextOp(ctx, cq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{character.Label}
	default:
		return nil, &NotSingularError{character.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (cq *CharacterQuery) OnlyX(ctx context.Context) *Character {
	node, err := cq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Character ID in the query.
// Returns a *NotSingularError when more than one Character ID is found.
// Returns a *NotFoundError when no entities are found.
func (cq *CharacterQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = cq.Limit(2).IDs(setContextOp(ctx, cq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{character.Label}
	default:
		err = &NotSingularError{character.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (cq *CharacterQuery) OnlyIDX(ctx context.Context) string {
	id, err := cq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Characters.
func (cq *CharacterQuery) All(ctx context.Context) ([]*Character, error) {
	ctx = setContextOp(ctx, cq.ctx, "All")
	if err := cq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Character, *CharacterQuery]()
	return withInterceptors[[]*Character](ctx, cq, qr, cq.inters)
}

// AllX is like All, but panics if an error occurs.
func (cq *CharacterQuery) AllX(ctx context.Context) []*Character {
	nodes, err := cq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Character IDs.
func (cq *CharacterQuery) IDs(ctx context.Context) (ids []string, err error) {
	if cq.ctx.Unique == nil && cq.path != nil {
		cq.Unique(true)
	}
	ctx = setContextOp(ctx, cq.ctx, "IDs")
	if err = cq.Select(character.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (cq *CharacterQuery) IDsX(ctx context.Context) []string {
	ids, err := cq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (cq *CharacterQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, cq.ctx, "Count")
	if err := cq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, cq, querierCount[*CharacterQuery](), cq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (cq *CharacterQuery) CountX(ctx context.Context) int {
	count, err := cq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (cq *CharacterQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, cq.ctx, "Exist")
	switch _, err := cq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (cq *CharacterQuery) ExistX(ctx context.Context) bool {
	exist, err := cq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CharacterQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (cq *CharacterQuery) Clone() *CharacterQuery {
	if cq == nil {
		return nil
	}
	return &CharacterQuery{
		config:     cq.config,
		ctx:        cq.ctx.Clone(),
		order:      append([]character.OrderOption{}, cq.order...),
		inters:     append([]Interceptor{}, cq.inters...),
		predicates: append([]predicate.Character{}, cq.predicates...),
		// clone intermediate query.
		sql:  cq.sql.Clone(),
		path: cq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		PinYin string `json:"pin_yin,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Character.Query().
//		GroupBy(character.FieldPinYin).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (cq *CharacterQuery) GroupBy(field string, fields ...string) *CharacterGroupBy {
	cq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &CharacterGroupBy{build: cq}
	grbuild.flds = &cq.ctx.Fields
	grbuild.label = character.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		PinYin string `json:"pin_yin,omitempty"`
//	}
//
//	client.Character.Query().
//		Select(character.FieldPinYin).
//		Scan(ctx, &v)
func (cq *CharacterQuery) Select(fields ...string) *CharacterSelect {
	cq.ctx.Fields = append(cq.ctx.Fields, fields...)
	sbuild := &CharacterSelect{CharacterQuery: cq}
	sbuild.label = character.Label
	sbuild.flds, sbuild.scan = &cq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a CharacterSelect configured with the given aggregations.
func (cq *CharacterQuery) Aggregate(fns ...AggregateFunc) *CharacterSelect {
	return cq.Select().Aggregate(fns...)
}

func (cq *CharacterQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range cq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, cq); err != nil {
				return err
			}
		}
	}
	for _, f := range cq.ctx.Fields {
		if !character.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if cq.path != nil {
		prev, err := cq.path(ctx)
		if err != nil {
			return err
		}
		cq.sql = prev
	}
	return nil
}

func (cq *CharacterQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Character, error) {
	var (
		nodes = []*Character{}
		_spec = cq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Character).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Character{config: cq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(cq.modifiers) > 0 {
		_spec.Modifiers = cq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, cq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (cq *CharacterQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := cq.querySpec()
	if len(cq.modifiers) > 0 {
		_spec.Modifiers = cq.modifiers
	}
	_spec.Node.Columns = cq.ctx.Fields
	if len(cq.ctx.Fields) > 0 {
		_spec.Unique = cq.ctx.Unique != nil && *cq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, cq.driver, _spec)
}

func (cq *CharacterQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(character.Table, character.Columns, sqlgraph.NewFieldSpec(character.FieldID, field.TypeString))
	_spec.From = cq.sql
	if unique := cq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if cq.path != nil {
		_spec.Unique = true
	}
	if fields := cq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, character.FieldID)
		for i := range fields {
			if fields[i] != character.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := cq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := cq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := cq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := cq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (cq *CharacterQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(cq.driver.Dialect())
	t1 := builder.Table(character.Table)
	columns := cq.ctx.Fields
	if len(columns) == 0 {
		columns = character.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if cq.sql != nil {
		selector = cq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if cq.ctx.Unique != nil && *cq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range cq.modifiers {
		m(selector)
	}
	for _, p := range cq.predicates {
		p(selector)
	}
	for _, p := range cq.order {
		p(selector)
	}
	if offset := cq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := cq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (cq *CharacterQuery) ForUpdate(opts ...sql.LockOption) *CharacterQuery {
	if cq.driver.Dialect() == dialect.Postgres {
		cq.Unique(false)
	}
	cq.modifiers = append(cq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return cq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (cq *CharacterQuery) ForShare(opts ...sql.LockOption) *CharacterQuery {
	if cq.driver.Dialect() == dialect.Postgres {
		cq.Unique(false)
	}
	cq.modifiers = append(cq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return cq
}

// CharacterGroupBy is the group-by builder for Character entities.
type CharacterGroupBy struct {
	selector
	build *CharacterQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (cgb *CharacterGroupBy) Aggregate(fns ...AggregateFunc) *CharacterGroupBy {
	cgb.fns = append(cgb.fns, fns...)
	return cgb
}

// Scan applies the selector query and scans the result into the given value.
func (cgb *CharacterGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, cgb.build.ctx, "GroupBy")
	if err := cgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CharacterQuery, *CharacterGroupBy](ctx, cgb.build, cgb, cgb.build.inters, v)
}

func (cgb *CharacterGroupBy) sqlScan(ctx context.Context, root *CharacterQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(cgb.fns))
	for _, fn := range cgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*cgb.flds)+len(cgb.fns))
		for _, f := range *cgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*cgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// CharacterSelect is the builder for selecting fields of Character entities.
type CharacterSelect struct {
	*CharacterQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (cs *CharacterSelect) Aggregate(fns ...AggregateFunc) *CharacterSelect {
	cs.fns = append(cs.fns, fns...)
	return cs
}

// Scan applies the selector query and scans the result into the given value.
func (cs *CharacterSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, cs.ctx, "Select")
	if err := cs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CharacterQuery, *CharacterSelect](ctx, cs.CharacterQuery, cs, cs.inters, v)
}

func (cs *CharacterSelect) sqlScan(ctx context.Context, root *CharacterQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(cs.fns))
	for _, fn := range cs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*cs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
