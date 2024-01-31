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
	"github.com/babyname/fate/ent/predicate"
	"github.com/babyname/fate/ent/wuxing"
)

// WuXingQuery is the builder for querying WuXing entities.
type WuXingQuery struct {
	config
	ctx        *QueryContext
	order      []OrderFunc
	inters     []Interceptor
	predicates []predicate.WuXing
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the WuXingQuery builder.
func (wxq *WuXingQuery) Where(ps ...predicate.WuXing) *WuXingQuery {
	wxq.predicates = append(wxq.predicates, ps...)
	return wxq
}

// Limit the number of records to be returned by this query.
func (wxq *WuXingQuery) Limit(limit int) *WuXingQuery {
	wxq.ctx.Limit = &limit
	return wxq
}

// Offset to start from.
func (wxq *WuXingQuery) Offset(offset int) *WuXingQuery {
	wxq.ctx.Offset = &offset
	return wxq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (wxq *WuXingQuery) Unique(unique bool) *WuXingQuery {
	wxq.ctx.Unique = &unique
	return wxq
}

// Order specifies how the records should be ordered.
func (wxq *WuXingQuery) Order(o ...OrderFunc) *WuXingQuery {
	wxq.order = append(wxq.order, o...)
	return wxq
}

// First returns the first WuXing entity from the query.
// Returns a *NotFoundError when no WuXing was found.
func (wxq *WuXingQuery) First(ctx context.Context) (*WuXing, error) {
	nodes, err := wxq.Limit(1).All(setContextOp(ctx, wxq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{wuxing.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (wxq *WuXingQuery) FirstX(ctx context.Context) *WuXing {
	node, err := wxq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first WuXing ID from the query.
// Returns a *NotFoundError when no WuXing ID was found.
func (wxq *WuXingQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = wxq.Limit(1).IDs(setContextOp(ctx, wxq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{wuxing.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (wxq *WuXingQuery) FirstIDX(ctx context.Context) string {
	id, err := wxq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single WuXing entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one WuXing entity is found.
// Returns a *NotFoundError when no WuXing entities are found.
func (wxq *WuXingQuery) Only(ctx context.Context) (*WuXing, error) {
	nodes, err := wxq.Limit(2).All(setContextOp(ctx, wxq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{wuxing.Label}
	default:
		return nil, &NotSingularError{wuxing.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (wxq *WuXingQuery) OnlyX(ctx context.Context) *WuXing {
	node, err := wxq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only WuXing ID in the query.
// Returns a *NotSingularError when more than one WuXing ID is found.
// Returns a *NotFoundError when no entities are found.
func (wxq *WuXingQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = wxq.Limit(2).IDs(setContextOp(ctx, wxq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{wuxing.Label}
	default:
		err = &NotSingularError{wuxing.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (wxq *WuXingQuery) OnlyIDX(ctx context.Context) string {
	id, err := wxq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of WuXings.
func (wxq *WuXingQuery) All(ctx context.Context) ([]*WuXing, error) {
	ctx = setContextOp(ctx, wxq.ctx, "All")
	if err := wxq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*WuXing, *WuXingQuery]()
	return withInterceptors[[]*WuXing](ctx, wxq, qr, wxq.inters)
}

// AllX is like All, but panics if an error occurs.
func (wxq *WuXingQuery) AllX(ctx context.Context) []*WuXing {
	nodes, err := wxq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of WuXing IDs.
func (wxq *WuXingQuery) IDs(ctx context.Context) (ids []string, err error) {
	if wxq.ctx.Unique == nil && wxq.path != nil {
		wxq.Unique(true)
	}
	ctx = setContextOp(ctx, wxq.ctx, "IDs")
	if err = wxq.Select(wuxing.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (wxq *WuXingQuery) IDsX(ctx context.Context) []string {
	ids, err := wxq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (wxq *WuXingQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, wxq.ctx, "Count")
	if err := wxq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, wxq, querierCount[*WuXingQuery](), wxq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (wxq *WuXingQuery) CountX(ctx context.Context) int {
	count, err := wxq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (wxq *WuXingQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, wxq.ctx, "Exist")
	switch _, err := wxq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (wxq *WuXingQuery) ExistX(ctx context.Context) bool {
	exist, err := wxq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the WuXingQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (wxq *WuXingQuery) Clone() *WuXingQuery {
	if wxq == nil {
		return nil
	}
	return &WuXingQuery{
		config:     wxq.config,
		ctx:        wxq.ctx.Clone(),
		order:      append([]OrderFunc{}, wxq.order...),
		inters:     append([]Interceptor{}, wxq.inters...),
		predicates: append([]predicate.WuXing{}, wxq.predicates...),
		// clone intermediate query.
		sql:  wxq.sql.Clone(),
		path: wxq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Created time.Time `json:"created,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.WuXing.Query().
//		GroupBy(wuxing.FieldCreated).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (wxq *WuXingQuery) GroupBy(field string, fields ...string) *WuXingGroupBy {
	wxq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &WuXingGroupBy{build: wxq}
	grbuild.flds = &wxq.ctx.Fields
	grbuild.label = wuxing.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Created time.Time `json:"created,omitempty"`
//	}
//
//	client.WuXing.Query().
//		Select(wuxing.FieldCreated).
//		Scan(ctx, &v)
func (wxq *WuXingQuery) Select(fields ...string) *WuXingSelect {
	wxq.ctx.Fields = append(wxq.ctx.Fields, fields...)
	sbuild := &WuXingSelect{WuXingQuery: wxq}
	sbuild.label = wuxing.Label
	sbuild.flds, sbuild.scan = &wxq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a WuXingSelect configured with the given aggregations.
func (wxq *WuXingQuery) Aggregate(fns ...AggregateFunc) *WuXingSelect {
	return wxq.Select().Aggregate(fns...)
}

func (wxq *WuXingQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range wxq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, wxq); err != nil {
				return err
			}
		}
	}
	for _, f := range wxq.ctx.Fields {
		if !wuxing.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if wxq.path != nil {
		prev, err := wxq.path(ctx)
		if err != nil {
			return err
		}
		wxq.sql = prev
	}
	return nil
}

func (wxq *WuXingQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*WuXing, error) {
	var (
		nodes = []*WuXing{}
		_spec = wxq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*WuXing).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &WuXing{config: wxq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(wxq.modifiers) > 0 {
		_spec.Modifiers = wxq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, wxq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (wxq *WuXingQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := wxq.querySpec()
	if len(wxq.modifiers) > 0 {
		_spec.Modifiers = wxq.modifiers
	}
	_spec.Node.Columns = wxq.ctx.Fields
	if len(wxq.ctx.Fields) > 0 {
		_spec.Unique = wxq.ctx.Unique != nil && *wxq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, wxq.driver, _spec)
}

func (wxq *WuXingQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(wuxing.Table, wuxing.Columns, sqlgraph.NewFieldSpec(wuxing.FieldID, field.TypeString))
	_spec.From = wxq.sql
	if unique := wxq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if wxq.path != nil {
		_spec.Unique = true
	}
	if fields := wxq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, wuxing.FieldID)
		for i := range fields {
			if fields[i] != wuxing.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := wxq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := wxq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := wxq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := wxq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (wxq *WuXingQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(wxq.driver.Dialect())
	t1 := builder.Table(wuxing.Table)
	columns := wxq.ctx.Fields
	if len(columns) == 0 {
		columns = wuxing.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if wxq.sql != nil {
		selector = wxq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if wxq.ctx.Unique != nil && *wxq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range wxq.modifiers {
		m(selector)
	}
	for _, p := range wxq.predicates {
		p(selector)
	}
	for _, p := range wxq.order {
		p(selector)
	}
	if offset := wxq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := wxq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (wxq *WuXingQuery) ForUpdate(opts ...sql.LockOption) *WuXingQuery {
	if wxq.driver.Dialect() == dialect.Postgres {
		wxq.Unique(false)
	}
	wxq.modifiers = append(wxq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return wxq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (wxq *WuXingQuery) ForShare(opts ...sql.LockOption) *WuXingQuery {
	if wxq.driver.Dialect() == dialect.Postgres {
		wxq.Unique(false)
	}
	wxq.modifiers = append(wxq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return wxq
}

// WuXingGroupBy is the group-by builder for WuXing entities.
type WuXingGroupBy struct {
	selector
	build *WuXingQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (wxgb *WuXingGroupBy) Aggregate(fns ...AggregateFunc) *WuXingGroupBy {
	wxgb.fns = append(wxgb.fns, fns...)
	return wxgb
}

// Scan applies the selector query and scans the result into the given value.
func (wxgb *WuXingGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, wxgb.build.ctx, "GroupBy")
	if err := wxgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*WuXingQuery, *WuXingGroupBy](ctx, wxgb.build, wxgb, wxgb.build.inters, v)
}

func (wxgb *WuXingGroupBy) sqlScan(ctx context.Context, root *WuXingQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(wxgb.fns))
	for _, fn := range wxgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*wxgb.flds)+len(wxgb.fns))
		for _, f := range *wxgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*wxgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := wxgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// WuXingSelect is the builder for selecting fields of WuXing entities.
type WuXingSelect struct {
	*WuXingQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (wxs *WuXingSelect) Aggregate(fns ...AggregateFunc) *WuXingSelect {
	wxs.fns = append(wxs.fns, fns...)
	return wxs
}

// Scan applies the selector query and scans the result into the given value.
func (wxs *WuXingSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, wxs.ctx, "Select")
	if err := wxs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*WuXingQuery, *WuXingSelect](ctx, wxs.WuXingQuery, wxs, wxs.inters, v)
}

func (wxs *WuXingSelect) sqlScan(ctx context.Context, root *WuXingQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(wxs.fns))
	for _, fn := range wxs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*wxs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := wxs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
