// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/tereus-project/tereus-api/ent/migrate"

	"github.com/tereus-project/tereus-api/ent/submission"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Submission is the client for interacting with the Submission builders.
	Submission *SubmissionClient
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
	c.Submission = NewSubmissionClient(c.config)
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
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
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
		Submission: NewSubmissionClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
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
		Submission: NewSubmissionClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Submission.
//		Query().
//		Count(ctx)
//
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
	c.Submission.Use(hooks...)
}

// SubmissionClient is a client for the Submission schema.
type SubmissionClient struct {
	config
}

// NewSubmissionClient returns a client for the Submission from the given config.
func NewSubmissionClient(c config) *SubmissionClient {
	return &SubmissionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `submission.Hooks(f(g(h())))`.
func (c *SubmissionClient) Use(hooks ...Hook) {
	c.hooks.Submission = append(c.hooks.Submission, hooks...)
}

// Create returns a create builder for Submission.
func (c *SubmissionClient) Create() *SubmissionCreate {
	mutation := newSubmissionMutation(c.config, OpCreate)
	return &SubmissionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Submission entities.
func (c *SubmissionClient) CreateBulk(builders ...*SubmissionCreate) *SubmissionCreateBulk {
	return &SubmissionCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Submission.
func (c *SubmissionClient) Update() *SubmissionUpdate {
	mutation := newSubmissionMutation(c.config, OpUpdate)
	return &SubmissionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SubmissionClient) UpdateOne(s *Submission) *SubmissionUpdateOne {
	mutation := newSubmissionMutation(c.config, OpUpdateOne, withSubmission(s))
	return &SubmissionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SubmissionClient) UpdateOneID(id uuid.UUID) *SubmissionUpdateOne {
	mutation := newSubmissionMutation(c.config, OpUpdateOne, withSubmissionID(id))
	return &SubmissionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Submission.
func (c *SubmissionClient) Delete() *SubmissionDelete {
	mutation := newSubmissionMutation(c.config, OpDelete)
	return &SubmissionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SubmissionClient) DeleteOne(s *Submission) *SubmissionDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SubmissionClient) DeleteOneID(id uuid.UUID) *SubmissionDeleteOne {
	builder := c.Delete().Where(submission.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SubmissionDeleteOne{builder}
}

// Query returns a query builder for Submission.
func (c *SubmissionClient) Query() *SubmissionQuery {
	return &SubmissionQuery{
		config: c.config,
	}
}

// Get returns a Submission entity by its id.
func (c *SubmissionClient) Get(ctx context.Context, id uuid.UUID) (*Submission, error) {
	return c.Query().Where(submission.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SubmissionClient) GetX(ctx context.Context, id uuid.UUID) *Submission {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *SubmissionClient) Hooks() []Hook {
	return c.hooks.Submission
}
