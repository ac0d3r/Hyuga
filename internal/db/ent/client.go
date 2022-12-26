// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"

	"hyuga/internal/db/ent/migrate"

	"hyuga/internal/db/ent/record"
	"hyuga/internal/db/ent/systemconfig"
	"hyuga/internal/db/ent/user"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Record is the client for interacting with the Record builders.
	Record *RecordClient
	// SystemConfig is the client for interacting with the SystemConfig builders.
	SystemConfig *SystemConfigClient
	// User is the client for interacting with the User builders.
	User *UserClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Record = NewRecordClient(c.config)
	c.SystemConfig = NewSystemConfigClient(c.config)
	c.User = NewUserClient(c.config)
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
		ctx:          ctx,
		config:       cfg,
		Record:       NewRecordClient(cfg),
		SystemConfig: NewSystemConfigClient(cfg),
		User:         NewUserClient(cfg),
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
		ctx:          ctx,
		config:       cfg,
		Record:       NewRecordClient(cfg),
		SystemConfig: NewSystemConfigClient(cfg),
		User:         NewUserClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Record.
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
	c.Record.Use(hooks...)
	c.SystemConfig.Use(hooks...)
	c.User.Use(hooks...)
}

// RecordClient is a client for the Record schema.
type RecordClient struct {
	config
}

// NewRecordClient returns a client for the Record from the given config.
func NewRecordClient(c config) *RecordClient {
	return &RecordClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `record.Hooks(f(g(h())))`.
func (c *RecordClient) Use(hooks ...Hook) {
	c.hooks.Record = append(c.hooks.Record, hooks...)
}

// Create returns a builder for creating a Record entity.
func (c *RecordClient) Create() *RecordCreate {
	mutation := newRecordMutation(c.config, OpCreate)
	return &RecordCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Record entities.
func (c *RecordClient) CreateBulk(builders ...*RecordCreate) *RecordCreateBulk {
	return &RecordCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Record.
func (c *RecordClient) Update() *RecordUpdate {
	mutation := newRecordMutation(c.config, OpUpdate)
	return &RecordUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *RecordClient) UpdateOne(r *Record) *RecordUpdateOne {
	mutation := newRecordMutation(c.config, OpUpdateOne, withRecord(r))
	return &RecordUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *RecordClient) UpdateOneID(id int) *RecordUpdateOne {
	mutation := newRecordMutation(c.config, OpUpdateOne, withRecordID(id))
	return &RecordUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Record.
func (c *RecordClient) Delete() *RecordDelete {
	mutation := newRecordMutation(c.config, OpDelete)
	return &RecordDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *RecordClient) DeleteOne(r *Record) *RecordDeleteOne {
	return c.DeleteOneID(r.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *RecordClient) DeleteOneID(id int) *RecordDeleteOne {
	builder := c.Delete().Where(record.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &RecordDeleteOne{builder}
}

// Query returns a query builder for Record.
func (c *RecordClient) Query() *RecordQuery {
	return &RecordQuery{
		config: c.config,
	}
}

// Get returns a Record entity by its id.
func (c *RecordClient) Get(ctx context.Context, id int) (*Record, error) {
	return c.Query().Where(record.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *RecordClient) GetX(ctx context.Context, id int) *Record {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *RecordClient) Hooks() []Hook {
	return c.hooks.Record
}

// SystemConfigClient is a client for the SystemConfig schema.
type SystemConfigClient struct {
	config
}

// NewSystemConfigClient returns a client for the SystemConfig from the given config.
func NewSystemConfigClient(c config) *SystemConfigClient {
	return &SystemConfigClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `systemconfig.Hooks(f(g(h())))`.
func (c *SystemConfigClient) Use(hooks ...Hook) {
	c.hooks.SystemConfig = append(c.hooks.SystemConfig, hooks...)
}

// Create returns a builder for creating a SystemConfig entity.
func (c *SystemConfigClient) Create() *SystemConfigCreate {
	mutation := newSystemConfigMutation(c.config, OpCreate)
	return &SystemConfigCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of SystemConfig entities.
func (c *SystemConfigClient) CreateBulk(builders ...*SystemConfigCreate) *SystemConfigCreateBulk {
	return &SystemConfigCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for SystemConfig.
func (c *SystemConfigClient) Update() *SystemConfigUpdate {
	mutation := newSystemConfigMutation(c.config, OpUpdate)
	return &SystemConfigUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SystemConfigClient) UpdateOne(sc *SystemConfig) *SystemConfigUpdateOne {
	mutation := newSystemConfigMutation(c.config, OpUpdateOne, withSystemConfig(sc))
	return &SystemConfigUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SystemConfigClient) UpdateOneID(id int) *SystemConfigUpdateOne {
	mutation := newSystemConfigMutation(c.config, OpUpdateOne, withSystemConfigID(id))
	return &SystemConfigUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for SystemConfig.
func (c *SystemConfigClient) Delete() *SystemConfigDelete {
	mutation := newSystemConfigMutation(c.config, OpDelete)
	return &SystemConfigDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *SystemConfigClient) DeleteOne(sc *SystemConfig) *SystemConfigDeleteOne {
	return c.DeleteOneID(sc.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *SystemConfigClient) DeleteOneID(id int) *SystemConfigDeleteOne {
	builder := c.Delete().Where(systemconfig.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SystemConfigDeleteOne{builder}
}

// Query returns a query builder for SystemConfig.
func (c *SystemConfigClient) Query() *SystemConfigQuery {
	return &SystemConfigQuery{
		config: c.config,
	}
}

// Get returns a SystemConfig entity by its id.
func (c *SystemConfigClient) Get(ctx context.Context, id int) (*SystemConfig, error) {
	return c.Query().Where(systemconfig.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SystemConfigClient) GetX(ctx context.Context, id int) *SystemConfig {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *SystemConfigClient) Hooks() []Hook {
	return c.hooks.SystemConfig
}

// UserClient is a client for the User schema.
type UserClient struct {
	config
}

// NewUserClient returns a client for the User from the given config.
func NewUserClient(c config) *UserClient {
	return &UserClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `user.Hooks(f(g(h())))`.
func (c *UserClient) Use(hooks ...Hook) {
	c.hooks.User = append(c.hooks.User, hooks...)
}

// Create returns a builder for creating a User entity.
func (c *UserClient) Create() *UserCreate {
	mutation := newUserMutation(c.config, OpCreate)
	return &UserCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of User entities.
func (c *UserClient) CreateBulk(builders ...*UserCreate) *UserCreateBulk {
	return &UserCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for User.
func (c *UserClient) Update() *UserUpdate {
	mutation := newUserMutation(c.config, OpUpdate)
	return &UserUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserClient) UpdateOne(u *User) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUser(u))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *UserClient) UpdateOneID(id int) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUserID(id))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for User.
func (c *UserClient) Delete() *UserDelete {
	mutation := newUserMutation(c.config, OpDelete)
	return &UserDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *UserClient) DeleteOne(u *User) *UserDeleteOne {
	return c.DeleteOneID(u.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *UserClient) DeleteOneID(id int) *UserDeleteOne {
	builder := c.Delete().Where(user.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UserDeleteOne{builder}
}

// Query returns a query builder for User.
func (c *UserClient) Query() *UserQuery {
	return &UserQuery{
		config: c.config,
	}
}

// Get returns a User entity by its id.
func (c *UserClient) Get(ctx context.Context, id int) (*User, error) {
	return c.Query().Where(user.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserClient) GetX(ctx context.Context, id int) *User {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *UserClient) Hooks() []Hook {
	return c.hooks.User
}
