// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/babyname/fate/ent/migrate"
	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/babyname/fate/ent/character"
	"github.com/babyname/fate/ent/idiom"
	"github.com/babyname/fate/ent/ncharacter"
	"github.com/babyname/fate/ent/version"
	"github.com/babyname/fate/ent/wugelucky"
	"github.com/babyname/fate/ent/wuxing"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Character is the client for interacting with the Character builders.
	Character *CharacterClient
	// Idiom is the client for interacting with the Idiom builders.
	Idiom *IdiomClient
	// NCharacter is the client for interacting with the NCharacter builders.
	NCharacter *NCharacterClient
	// Version is the client for interacting with the Version builders.
	Version *VersionClient
	// WuGeLucky is the client for interacting with the WuGeLucky builders.
	WuGeLucky *WuGeLuckyClient
	// WuXing is the client for interacting with the WuXing builders.
	WuXing *WuXingClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Character = NewCharacterClient(c.config)
	c.Idiom = NewIdiomClient(c.config)
	c.NCharacter = NewNCharacterClient(c.config)
	c.Version = NewVersionClient(c.config)
	c.WuGeLucky = NewWuGeLuckyClient(c.config)
	c.WuXing = NewWuXingClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:        ctx,
		config:     cfg,
		Character:  NewCharacterClient(cfg),
		Idiom:      NewIdiomClient(cfg),
		NCharacter: NewNCharacterClient(cfg),
		Version:    NewVersionClient(cfg),
		WuGeLucky:  NewWuGeLuckyClient(cfg),
		WuXing:     NewWuXingClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:        ctx,
		config:     cfg,
		Character:  NewCharacterClient(cfg),
		Idiom:      NewIdiomClient(cfg),
		NCharacter: NewNCharacterClient(cfg),
		Version:    NewVersionClient(cfg),
		WuGeLucky:  NewWuGeLuckyClient(cfg),
		WuXing:     NewWuXingClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Character.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	for _, n := range []interface{ Use(...Hook) }{
		c.Character, c.Idiom, c.NCharacter, c.Version, c.WuGeLucky, c.WuXing,
	} {
		n.Use(hooks...)
	}
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	for _, n := range []interface{ Intercept(...Interceptor) }{
		c.Character, c.Idiom, c.NCharacter, c.Version, c.WuGeLucky, c.WuXing,
	} {
		n.Intercept(interceptors...)
	}
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *CharacterMutation:
		return c.Character.mutate(ctx, m)
	case *IdiomMutation:
		return c.Idiom.mutate(ctx, m)
	case *NCharacterMutation:
		return c.NCharacter.mutate(ctx, m)
	case *VersionMutation:
		return c.Version.mutate(ctx, m)
	case *WuGeLuckyMutation:
		return c.WuGeLucky.mutate(ctx, m)
	case *WuXingMutation:
		return c.WuXing.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// CharacterClient is a client for the Character schema.
type CharacterClient struct {
	config
}

// NewCharacterClient returns a client for the Character from the given config.
func NewCharacterClient(c config) *CharacterClient {
	return &CharacterClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `character.Hooks(f(g(h())))`.
func (c *CharacterClient) Use(hooks ...Hook) {
	c.hooks.Character = append(c.hooks.Character, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `character.Intercept(f(g(h())))`.
func (c *CharacterClient) Intercept(interceptors ...Interceptor) {
	c.inters.Character = append(c.inters.Character, interceptors...)
}

// Create returns a builder for creating a Character entity.
func (c *CharacterClient) Create() *CharacterCreate {
	mutation := newCharacterMutation(c.config, OpCreate)
	return &CharacterCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Character entities.
func (c *CharacterClient) CreateBulk(builders ...*CharacterCreate) *CharacterCreateBulk {
	return &CharacterCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Character.
func (c *CharacterClient) Update() *CharacterUpdate {
	mutation := newCharacterMutation(c.config, OpUpdate)
	return &CharacterUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *CharacterClient) UpdateOne(ch *Character) *CharacterUpdateOne {
	mutation := newCharacterMutation(c.config, OpUpdateOne, withCharacter(ch))
	return &CharacterUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *CharacterClient) UpdateOneID(id string) *CharacterUpdateOne {
	mutation := newCharacterMutation(c.config, OpUpdateOne, withCharacterID(id))
	return &CharacterUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Character.
func (c *CharacterClient) Delete() *CharacterDelete {
	mutation := newCharacterMutation(c.config, OpDelete)
	return &CharacterDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *CharacterClient) DeleteOne(ch *Character) *CharacterDeleteOne {
	return c.DeleteOneID(ch.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *CharacterClient) DeleteOneID(id string) *CharacterDeleteOne {
	builder := c.Delete().Where(character.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &CharacterDeleteOne{builder}
}

// Query returns a query builder for Character.
func (c *CharacterClient) Query() *CharacterQuery {
	return &CharacterQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeCharacter},
		inters: c.Interceptors(),
	}
}

// Get returns a Character entity by its id.
func (c *CharacterClient) Get(ctx context.Context, id string) (*Character, error) {
	return c.Query().Where(character.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CharacterClient) GetX(ctx context.Context, id string) *Character {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *CharacterClient) Hooks() []Hook {
	return c.hooks.Character
}

// Interceptors returns the client interceptors.
func (c *CharacterClient) Interceptors() []Interceptor {
	return c.inters.Character
}

func (c *CharacterClient) mutate(ctx context.Context, m *CharacterMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&CharacterCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&CharacterUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&CharacterUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&CharacterDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Character mutation op: %q", m.Op())
	}
}

// IdiomClient is a client for the Idiom schema.
type IdiomClient struct {
	config
}

// NewIdiomClient returns a client for the Idiom from the given config.
func NewIdiomClient(c config) *IdiomClient {
	return &IdiomClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `idiom.Hooks(f(g(h())))`.
func (c *IdiomClient) Use(hooks ...Hook) {
	c.hooks.Idiom = append(c.hooks.Idiom, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `idiom.Intercept(f(g(h())))`.
func (c *IdiomClient) Intercept(interceptors ...Interceptor) {
	c.inters.Idiom = append(c.inters.Idiom, interceptors...)
}

// Create returns a builder for creating a Idiom entity.
func (c *IdiomClient) Create() *IdiomCreate {
	mutation := newIdiomMutation(c.config, OpCreate)
	return &IdiomCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Idiom entities.
func (c *IdiomClient) CreateBulk(builders ...*IdiomCreate) *IdiomCreateBulk {
	return &IdiomCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Idiom.
func (c *IdiomClient) Update() *IdiomUpdate {
	mutation := newIdiomMutation(c.config, OpUpdate)
	return &IdiomUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *IdiomClient) UpdateOne(i *Idiom) *IdiomUpdateOne {
	mutation := newIdiomMutation(c.config, OpUpdateOne, withIdiom(i))
	return &IdiomUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *IdiomClient) UpdateOneID(id int32) *IdiomUpdateOne {
	mutation := newIdiomMutation(c.config, OpUpdateOne, withIdiomID(id))
	return &IdiomUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Idiom.
func (c *IdiomClient) Delete() *IdiomDelete {
	mutation := newIdiomMutation(c.config, OpDelete)
	return &IdiomDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *IdiomClient) DeleteOne(i *Idiom) *IdiomDeleteOne {
	return c.DeleteOneID(i.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *IdiomClient) DeleteOneID(id int32) *IdiomDeleteOne {
	builder := c.Delete().Where(idiom.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &IdiomDeleteOne{builder}
}

// Query returns a query builder for Idiom.
func (c *IdiomClient) Query() *IdiomQuery {
	return &IdiomQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeIdiom},
		inters: c.Interceptors(),
	}
}

// Get returns a Idiom entity by its id.
func (c *IdiomClient) Get(ctx context.Context, id int32) (*Idiom, error) {
	return c.Query().Where(idiom.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *IdiomClient) GetX(ctx context.Context, id int32) *Idiom {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *IdiomClient) Hooks() []Hook {
	return c.hooks.Idiom
}

// Interceptors returns the client interceptors.
func (c *IdiomClient) Interceptors() []Interceptor {
	return c.inters.Idiom
}

func (c *IdiomClient) mutate(ctx context.Context, m *IdiomMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&IdiomCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&IdiomUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&IdiomUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&IdiomDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Idiom mutation op: %q", m.Op())
	}
}

// NCharacterClient is a client for the NCharacter schema.
type NCharacterClient struct {
	config
}

// NewNCharacterClient returns a client for the NCharacter from the given config.
func NewNCharacterClient(c config) *NCharacterClient {
	return &NCharacterClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `ncharacter.Hooks(f(g(h())))`.
func (c *NCharacterClient) Use(hooks ...Hook) {
	c.hooks.NCharacter = append(c.hooks.NCharacter, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `ncharacter.Intercept(f(g(h())))`.
func (c *NCharacterClient) Intercept(interceptors ...Interceptor) {
	c.inters.NCharacter = append(c.inters.NCharacter, interceptors...)
}

// Create returns a builder for creating a NCharacter entity.
func (c *NCharacterClient) Create() *NCharacterCreate {
	mutation := newNCharacterMutation(c.config, OpCreate)
	return &NCharacterCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of NCharacter entities.
func (c *NCharacterClient) CreateBulk(builders ...*NCharacterCreate) *NCharacterCreateBulk {
	return &NCharacterCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for NCharacter.
func (c *NCharacterClient) Update() *NCharacterUpdate {
	mutation := newNCharacterMutation(c.config, OpUpdate)
	return &NCharacterUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *NCharacterClient) UpdateOne(n *NCharacter) *NCharacterUpdateOne {
	mutation := newNCharacterMutation(c.config, OpUpdateOne, withNCharacter(n))
	return &NCharacterUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *NCharacterClient) UpdateOneID(id int) *NCharacterUpdateOne {
	mutation := newNCharacterMutation(c.config, OpUpdateOne, withNCharacterID(id))
	return &NCharacterUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for NCharacter.
func (c *NCharacterClient) Delete() *NCharacterDelete {
	mutation := newNCharacterMutation(c.config, OpDelete)
	return &NCharacterDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *NCharacterClient) DeleteOne(n *NCharacter) *NCharacterDeleteOne {
	return c.DeleteOneID(n.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *NCharacterClient) DeleteOneID(id int) *NCharacterDeleteOne {
	builder := c.Delete().Where(ncharacter.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &NCharacterDeleteOne{builder}
}

// Query returns a query builder for NCharacter.
func (c *NCharacterClient) Query() *NCharacterQuery {
	return &NCharacterQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeNCharacter},
		inters: c.Interceptors(),
	}
}

// Get returns a NCharacter entity by its id.
func (c *NCharacterClient) Get(ctx context.Context, id int) (*NCharacter, error) {
	return c.Query().Where(ncharacter.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *NCharacterClient) GetX(ctx context.Context, id int) *NCharacter {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *NCharacterClient) Hooks() []Hook {
	return c.hooks.NCharacter
}

// Interceptors returns the client interceptors.
func (c *NCharacterClient) Interceptors() []Interceptor {
	return c.inters.NCharacter
}

func (c *NCharacterClient) mutate(ctx context.Context, m *NCharacterMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&NCharacterCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&NCharacterUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&NCharacterUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&NCharacterDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown NCharacter mutation op: %q", m.Op())
	}
}

// VersionClient is a client for the Version schema.
type VersionClient struct {
	config
}

// NewVersionClient returns a client for the Version from the given config.
func NewVersionClient(c config) *VersionClient {
	return &VersionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `version.Hooks(f(g(h())))`.
func (c *VersionClient) Use(hooks ...Hook) {
	c.hooks.Version = append(c.hooks.Version, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `version.Intercept(f(g(h())))`.
func (c *VersionClient) Intercept(interceptors ...Interceptor) {
	c.inters.Version = append(c.inters.Version, interceptors...)
}

// Create returns a builder for creating a Version entity.
func (c *VersionClient) Create() *VersionCreate {
	mutation := newVersionMutation(c.config, OpCreate)
	return &VersionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Version entities.
func (c *VersionClient) CreateBulk(builders ...*VersionCreate) *VersionCreateBulk {
	return &VersionCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Version.
func (c *VersionClient) Update() *VersionUpdate {
	mutation := newVersionMutation(c.config, OpUpdate)
	return &VersionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *VersionClient) UpdateOne(v *Version) *VersionUpdateOne {
	mutation := newVersionMutation(c.config, OpUpdateOne, withVersion(v))
	return &VersionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *VersionClient) UpdateOneID(id int) *VersionUpdateOne {
	mutation := newVersionMutation(c.config, OpUpdateOne, withVersionID(id))
	return &VersionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Version.
func (c *VersionClient) Delete() *VersionDelete {
	mutation := newVersionMutation(c.config, OpDelete)
	return &VersionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *VersionClient) DeleteOne(v *Version) *VersionDeleteOne {
	return c.DeleteOneID(v.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *VersionClient) DeleteOneID(id int) *VersionDeleteOne {
	builder := c.Delete().Where(version.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &VersionDeleteOne{builder}
}

// Query returns a query builder for Version.
func (c *VersionClient) Query() *VersionQuery {
	return &VersionQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeVersion},
		inters: c.Interceptors(),
	}
}

// Get returns a Version entity by its id.
func (c *VersionClient) Get(ctx context.Context, id int) (*Version, error) {
	return c.Query().Where(version.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *VersionClient) GetX(ctx context.Context, id int) *Version {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *VersionClient) Hooks() []Hook {
	return c.hooks.Version
}

// Interceptors returns the client interceptors.
func (c *VersionClient) Interceptors() []Interceptor {
	return c.inters.Version
}

func (c *VersionClient) mutate(ctx context.Context, m *VersionMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&VersionCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&VersionUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&VersionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&VersionDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Version mutation op: %q", m.Op())
	}
}

// WuGeLuckyClient is a client for the WuGeLucky schema.
type WuGeLuckyClient struct {
	config
}

// NewWuGeLuckyClient returns a client for the WuGeLucky from the given config.
func NewWuGeLuckyClient(c config) *WuGeLuckyClient {
	return &WuGeLuckyClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `wugelucky.Hooks(f(g(h())))`.
func (c *WuGeLuckyClient) Use(hooks ...Hook) {
	c.hooks.WuGeLucky = append(c.hooks.WuGeLucky, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `wugelucky.Intercept(f(g(h())))`.
func (c *WuGeLuckyClient) Intercept(interceptors ...Interceptor) {
	c.inters.WuGeLucky = append(c.inters.WuGeLucky, interceptors...)
}

// Create returns a builder for creating a WuGeLucky entity.
func (c *WuGeLuckyClient) Create() *WuGeLuckyCreate {
	mutation := newWuGeLuckyMutation(c.config, OpCreate)
	return &WuGeLuckyCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of WuGeLucky entities.
func (c *WuGeLuckyClient) CreateBulk(builders ...*WuGeLuckyCreate) *WuGeLuckyCreateBulk {
	return &WuGeLuckyCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for WuGeLucky.
func (c *WuGeLuckyClient) Update() *WuGeLuckyUpdate {
	mutation := newWuGeLuckyMutation(c.config, OpUpdate)
	return &WuGeLuckyUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *WuGeLuckyClient) UpdateOne(wgl *WuGeLucky) *WuGeLuckyUpdateOne {
	mutation := newWuGeLuckyMutation(c.config, OpUpdateOne, withWuGeLucky(wgl))
	return &WuGeLuckyUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *WuGeLuckyClient) UpdateOneID(id uuid.UUID) *WuGeLuckyUpdateOne {
	mutation := newWuGeLuckyMutation(c.config, OpUpdateOne, withWuGeLuckyID(id))
	return &WuGeLuckyUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for WuGeLucky.
func (c *WuGeLuckyClient) Delete() *WuGeLuckyDelete {
	mutation := newWuGeLuckyMutation(c.config, OpDelete)
	return &WuGeLuckyDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *WuGeLuckyClient) DeleteOne(wgl *WuGeLucky) *WuGeLuckyDeleteOne {
	return c.DeleteOneID(wgl.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *WuGeLuckyClient) DeleteOneID(id uuid.UUID) *WuGeLuckyDeleteOne {
	builder := c.Delete().Where(wugelucky.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &WuGeLuckyDeleteOne{builder}
}

// Query returns a query builder for WuGeLucky.
func (c *WuGeLuckyClient) Query() *WuGeLuckyQuery {
	return &WuGeLuckyQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeWuGeLucky},
		inters: c.Interceptors(),
	}
}

// Get returns a WuGeLucky entity by its id.
func (c *WuGeLuckyClient) Get(ctx context.Context, id uuid.UUID) (*WuGeLucky, error) {
	return c.Query().Where(wugelucky.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WuGeLuckyClient) GetX(ctx context.Context, id uuid.UUID) *WuGeLucky {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *WuGeLuckyClient) Hooks() []Hook {
	return c.hooks.WuGeLucky
}

// Interceptors returns the client interceptors.
func (c *WuGeLuckyClient) Interceptors() []Interceptor {
	return c.inters.WuGeLucky
}

func (c *WuGeLuckyClient) mutate(ctx context.Context, m *WuGeLuckyMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&WuGeLuckyCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&WuGeLuckyUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&WuGeLuckyUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&WuGeLuckyDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown WuGeLucky mutation op: %q", m.Op())
	}
}

// WuXingClient is a client for the WuXing schema.
type WuXingClient struct {
	config
}

// NewWuXingClient returns a client for the WuXing from the given config.
func NewWuXingClient(c config) *WuXingClient {
	return &WuXingClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `wuxing.Hooks(f(g(h())))`.
func (c *WuXingClient) Use(hooks ...Hook) {
	c.hooks.WuXing = append(c.hooks.WuXing, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `wuxing.Intercept(f(g(h())))`.
func (c *WuXingClient) Intercept(interceptors ...Interceptor) {
	c.inters.WuXing = append(c.inters.WuXing, interceptors...)
}

// Create returns a builder for creating a WuXing entity.
func (c *WuXingClient) Create() *WuXingCreate {
	mutation := newWuXingMutation(c.config, OpCreate)
	return &WuXingCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of WuXing entities.
func (c *WuXingClient) CreateBulk(builders ...*WuXingCreate) *WuXingCreateBulk {
	return &WuXingCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for WuXing.
func (c *WuXingClient) Update() *WuXingUpdate {
	mutation := newWuXingMutation(c.config, OpUpdate)
	return &WuXingUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *WuXingClient) UpdateOne(wx *WuXing) *WuXingUpdateOne {
	mutation := newWuXingMutation(c.config, OpUpdateOne, withWuXing(wx))
	return &WuXingUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *WuXingClient) UpdateOneID(id string) *WuXingUpdateOne {
	mutation := newWuXingMutation(c.config, OpUpdateOne, withWuXingID(id))
	return &WuXingUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for WuXing.
func (c *WuXingClient) Delete() *WuXingDelete {
	mutation := newWuXingMutation(c.config, OpDelete)
	return &WuXingDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *WuXingClient) DeleteOne(wx *WuXing) *WuXingDeleteOne {
	return c.DeleteOneID(wx.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *WuXingClient) DeleteOneID(id string) *WuXingDeleteOne {
	builder := c.Delete().Where(wuxing.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &WuXingDeleteOne{builder}
}

// Query returns a query builder for WuXing.
func (c *WuXingClient) Query() *WuXingQuery {
	return &WuXingQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeWuXing},
		inters: c.Interceptors(),
	}
}

// Get returns a WuXing entity by its id.
func (c *WuXingClient) Get(ctx context.Context, id string) (*WuXing, error) {
	return c.Query().Where(wuxing.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WuXingClient) GetX(ctx context.Context, id string) *WuXing {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *WuXingClient) Hooks() []Hook {
	return c.hooks.WuXing
}

// Interceptors returns the client interceptors.
func (c *WuXingClient) Interceptors() []Interceptor {
	return c.inters.WuXing
}

func (c *WuXingClient) mutate(ctx context.Context, m *WuXingMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&WuXingCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&WuXingUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&WuXingUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&WuXingDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown WuXing mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Character, Idiom, NCharacter, Version, WuGeLucky, WuXing []ent.Hook
	}
	inters struct {
		Character, Idiom, NCharacter, Version, WuGeLucky, WuXing []ent.Interceptor
	}
)
