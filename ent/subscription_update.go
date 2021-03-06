// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/tereus-project/tereus-api/ent/predicate"
	"github.com/tereus-project/tereus-api/ent/subscription"
	"github.com/tereus-project/tereus-api/ent/user"
)

// SubscriptionUpdate is the builder for updating Subscription entities.
type SubscriptionUpdate struct {
	config
	hooks    []Hook
	mutation *SubscriptionMutation
}

// Where appends a list predicates to the SubscriptionUpdate builder.
func (su *SubscriptionUpdate) Where(ps ...predicate.Subscription) *SubscriptionUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetStripeCustomerID sets the "stripe_customer_id" field.
func (su *SubscriptionUpdate) SetStripeCustomerID(s string) *SubscriptionUpdate {
	su.mutation.SetStripeCustomerID(s)
	return su
}

// SetNillableStripeCustomerID sets the "stripe_customer_id" field if the given value is not nil.
func (su *SubscriptionUpdate) SetNillableStripeCustomerID(s *string) *SubscriptionUpdate {
	if s != nil {
		su.SetStripeCustomerID(*s)
	}
	return su
}

// ClearStripeCustomerID clears the value of the "stripe_customer_id" field.
func (su *SubscriptionUpdate) ClearStripeCustomerID() *SubscriptionUpdate {
	su.mutation.ClearStripeCustomerID()
	return su
}

// SetStripeSubscriptionID sets the "stripe_subscription_id" field.
func (su *SubscriptionUpdate) SetStripeSubscriptionID(s string) *SubscriptionUpdate {
	su.mutation.SetStripeSubscriptionID(s)
	return su
}

// SetNillableStripeSubscriptionID sets the "stripe_subscription_id" field if the given value is not nil.
func (su *SubscriptionUpdate) SetNillableStripeSubscriptionID(s *string) *SubscriptionUpdate {
	if s != nil {
		su.SetStripeSubscriptionID(*s)
	}
	return su
}

// ClearStripeSubscriptionID clears the value of the "stripe_subscription_id" field.
func (su *SubscriptionUpdate) ClearStripeSubscriptionID() *SubscriptionUpdate {
	su.mutation.ClearStripeSubscriptionID()
	return su
}

// SetTier sets the "tier" field.
func (su *SubscriptionUpdate) SetTier(s subscription.Tier) *SubscriptionUpdate {
	su.mutation.SetTier(s)
	return su
}

// SetNillableTier sets the "tier" field if the given value is not nil.
func (su *SubscriptionUpdate) SetNillableTier(s *subscription.Tier) *SubscriptionUpdate {
	if s != nil {
		su.SetTier(*s)
	}
	return su
}

// SetExpiresAt sets the "expires_at" field.
func (su *SubscriptionUpdate) SetExpiresAt(t time.Time) *SubscriptionUpdate {
	su.mutation.SetExpiresAt(t)
	return su
}

// SetNillableExpiresAt sets the "expires_at" field if the given value is not nil.
func (su *SubscriptionUpdate) SetNillableExpiresAt(t *time.Time) *SubscriptionUpdate {
	if t != nil {
		su.SetExpiresAt(*t)
	}
	return su
}

// ClearExpiresAt clears the value of the "expires_at" field.
func (su *SubscriptionUpdate) ClearExpiresAt() *SubscriptionUpdate {
	su.mutation.ClearExpiresAt()
	return su
}

// SetCancelled sets the "cancelled" field.
func (su *SubscriptionUpdate) SetCancelled(b bool) *SubscriptionUpdate {
	su.mutation.SetCancelled(b)
	return su
}

// SetNillableCancelled sets the "cancelled" field if the given value is not nil.
func (su *SubscriptionUpdate) SetNillableCancelled(b *bool) *SubscriptionUpdate {
	if b != nil {
		su.SetCancelled(*b)
	}
	return su
}

// SetCreatedAt sets the "created_at" field.
func (su *SubscriptionUpdate) SetCreatedAt(t time.Time) *SubscriptionUpdate {
	su.mutation.SetCreatedAt(t)
	return su
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (su *SubscriptionUpdate) SetNillableCreatedAt(t *time.Time) *SubscriptionUpdate {
	if t != nil {
		su.SetCreatedAt(*t)
	}
	return su
}

// SetUserID sets the "user" edge to the User entity by ID.
func (su *SubscriptionUpdate) SetUserID(id uuid.UUID) *SubscriptionUpdate {
	su.mutation.SetUserID(id)
	return su
}

// SetUser sets the "user" edge to the User entity.
func (su *SubscriptionUpdate) SetUser(u *User) *SubscriptionUpdate {
	return su.SetUserID(u.ID)
}

// Mutation returns the SubscriptionMutation object of the builder.
func (su *SubscriptionUpdate) Mutation() *SubscriptionMutation {
	return su.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (su *SubscriptionUpdate) ClearUser() *SubscriptionUpdate {
	su.mutation.ClearUser()
	return su
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SubscriptionUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(su.hooks) == 0 {
		if err = su.check(); err != nil {
			return 0, err
		}
		affected, err = su.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SubscriptionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = su.check(); err != nil {
				return 0, err
			}
			su.mutation = mutation
			affected, err = su.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(su.hooks) - 1; i >= 0; i-- {
			if su.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = su.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, su.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (su *SubscriptionUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SubscriptionUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SubscriptionUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *SubscriptionUpdate) check() error {
	if v, ok := su.mutation.Tier(); ok {
		if err := subscription.TierValidator(v); err != nil {
			return &ValidationError{Name: "tier", err: fmt.Errorf(`ent: validator failed for field "Subscription.tier": %w`, err)}
		}
	}
	if _, ok := su.mutation.UserID(); su.mutation.UserCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Subscription.user"`)
	}
	return nil
}

func (su *SubscriptionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   subscription.Table,
			Columns: subscription.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: subscription.FieldID,
			},
		},
	}
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.StripeCustomerID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: subscription.FieldStripeCustomerID,
		})
	}
	if su.mutation.StripeCustomerIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: subscription.FieldStripeCustomerID,
		})
	}
	if value, ok := su.mutation.StripeSubscriptionID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: subscription.FieldStripeSubscriptionID,
		})
	}
	if su.mutation.StripeSubscriptionIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: subscription.FieldStripeSubscriptionID,
		})
	}
	if value, ok := su.mutation.Tier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: subscription.FieldTier,
		})
	}
	if value, ok := su.mutation.ExpiresAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscription.FieldExpiresAt,
		})
	}
	if su.mutation.ExpiresAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: subscription.FieldExpiresAt,
		})
	}
	if value, ok := su.mutation.Cancelled(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: subscription.FieldCancelled,
		})
	}
	if value, ok := su.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscription.FieldCreatedAt,
		})
	}
	if su.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   subscription.UserTable,
			Columns: []string{subscription.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   subscription.UserTable,
			Columns: []string{subscription.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{subscription.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// SubscriptionUpdateOne is the builder for updating a single Subscription entity.
type SubscriptionUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SubscriptionMutation
}

// SetStripeCustomerID sets the "stripe_customer_id" field.
func (suo *SubscriptionUpdateOne) SetStripeCustomerID(s string) *SubscriptionUpdateOne {
	suo.mutation.SetStripeCustomerID(s)
	return suo
}

// SetNillableStripeCustomerID sets the "stripe_customer_id" field if the given value is not nil.
func (suo *SubscriptionUpdateOne) SetNillableStripeCustomerID(s *string) *SubscriptionUpdateOne {
	if s != nil {
		suo.SetStripeCustomerID(*s)
	}
	return suo
}

// ClearStripeCustomerID clears the value of the "stripe_customer_id" field.
func (suo *SubscriptionUpdateOne) ClearStripeCustomerID() *SubscriptionUpdateOne {
	suo.mutation.ClearStripeCustomerID()
	return suo
}

// SetStripeSubscriptionID sets the "stripe_subscription_id" field.
func (suo *SubscriptionUpdateOne) SetStripeSubscriptionID(s string) *SubscriptionUpdateOne {
	suo.mutation.SetStripeSubscriptionID(s)
	return suo
}

// SetNillableStripeSubscriptionID sets the "stripe_subscription_id" field if the given value is not nil.
func (suo *SubscriptionUpdateOne) SetNillableStripeSubscriptionID(s *string) *SubscriptionUpdateOne {
	if s != nil {
		suo.SetStripeSubscriptionID(*s)
	}
	return suo
}

// ClearStripeSubscriptionID clears the value of the "stripe_subscription_id" field.
func (suo *SubscriptionUpdateOne) ClearStripeSubscriptionID() *SubscriptionUpdateOne {
	suo.mutation.ClearStripeSubscriptionID()
	return suo
}

// SetTier sets the "tier" field.
func (suo *SubscriptionUpdateOne) SetTier(s subscription.Tier) *SubscriptionUpdateOne {
	suo.mutation.SetTier(s)
	return suo
}

// SetNillableTier sets the "tier" field if the given value is not nil.
func (suo *SubscriptionUpdateOne) SetNillableTier(s *subscription.Tier) *SubscriptionUpdateOne {
	if s != nil {
		suo.SetTier(*s)
	}
	return suo
}

// SetExpiresAt sets the "expires_at" field.
func (suo *SubscriptionUpdateOne) SetExpiresAt(t time.Time) *SubscriptionUpdateOne {
	suo.mutation.SetExpiresAt(t)
	return suo
}

// SetNillableExpiresAt sets the "expires_at" field if the given value is not nil.
func (suo *SubscriptionUpdateOne) SetNillableExpiresAt(t *time.Time) *SubscriptionUpdateOne {
	if t != nil {
		suo.SetExpiresAt(*t)
	}
	return suo
}

// ClearExpiresAt clears the value of the "expires_at" field.
func (suo *SubscriptionUpdateOne) ClearExpiresAt() *SubscriptionUpdateOne {
	suo.mutation.ClearExpiresAt()
	return suo
}

// SetCancelled sets the "cancelled" field.
func (suo *SubscriptionUpdateOne) SetCancelled(b bool) *SubscriptionUpdateOne {
	suo.mutation.SetCancelled(b)
	return suo
}

// SetNillableCancelled sets the "cancelled" field if the given value is not nil.
func (suo *SubscriptionUpdateOne) SetNillableCancelled(b *bool) *SubscriptionUpdateOne {
	if b != nil {
		suo.SetCancelled(*b)
	}
	return suo
}

// SetCreatedAt sets the "created_at" field.
func (suo *SubscriptionUpdateOne) SetCreatedAt(t time.Time) *SubscriptionUpdateOne {
	suo.mutation.SetCreatedAt(t)
	return suo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (suo *SubscriptionUpdateOne) SetNillableCreatedAt(t *time.Time) *SubscriptionUpdateOne {
	if t != nil {
		suo.SetCreatedAt(*t)
	}
	return suo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (suo *SubscriptionUpdateOne) SetUserID(id uuid.UUID) *SubscriptionUpdateOne {
	suo.mutation.SetUserID(id)
	return suo
}

// SetUser sets the "user" edge to the User entity.
func (suo *SubscriptionUpdateOne) SetUser(u *User) *SubscriptionUpdateOne {
	return suo.SetUserID(u.ID)
}

// Mutation returns the SubscriptionMutation object of the builder.
func (suo *SubscriptionUpdateOne) Mutation() *SubscriptionMutation {
	return suo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (suo *SubscriptionUpdateOne) ClearUser() *SubscriptionUpdateOne {
	suo.mutation.ClearUser()
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SubscriptionUpdateOne) Select(field string, fields ...string) *SubscriptionUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Subscription entity.
func (suo *SubscriptionUpdateOne) Save(ctx context.Context) (*Subscription, error) {
	var (
		err  error
		node *Subscription
	)
	if len(suo.hooks) == 0 {
		if err = suo.check(); err != nil {
			return nil, err
		}
		node, err = suo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SubscriptionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = suo.check(); err != nil {
				return nil, err
			}
			suo.mutation = mutation
			node, err = suo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(suo.hooks) - 1; i >= 0; i-- {
			if suo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = suo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, suo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (suo *SubscriptionUpdateOne) SaveX(ctx context.Context) *Subscription {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SubscriptionUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SubscriptionUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *SubscriptionUpdateOne) check() error {
	if v, ok := suo.mutation.Tier(); ok {
		if err := subscription.TierValidator(v); err != nil {
			return &ValidationError{Name: "tier", err: fmt.Errorf(`ent: validator failed for field "Subscription.tier": %w`, err)}
		}
	}
	if _, ok := suo.mutation.UserID(); suo.mutation.UserCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Subscription.user"`)
	}
	return nil
}

func (suo *SubscriptionUpdateOne) sqlSave(ctx context.Context) (_node *Subscription, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   subscription.Table,
			Columns: subscription.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: subscription.FieldID,
			},
		},
	}
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Subscription.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, subscription.FieldID)
		for _, f := range fields {
			if !subscription.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != subscription.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.StripeCustomerID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: subscription.FieldStripeCustomerID,
		})
	}
	if suo.mutation.StripeCustomerIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: subscription.FieldStripeCustomerID,
		})
	}
	if value, ok := suo.mutation.StripeSubscriptionID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: subscription.FieldStripeSubscriptionID,
		})
	}
	if suo.mutation.StripeSubscriptionIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: subscription.FieldStripeSubscriptionID,
		})
	}
	if value, ok := suo.mutation.Tier(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: subscription.FieldTier,
		})
	}
	if value, ok := suo.mutation.ExpiresAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscription.FieldExpiresAt,
		})
	}
	if suo.mutation.ExpiresAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: subscription.FieldExpiresAt,
		})
	}
	if value, ok := suo.mutation.Cancelled(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: subscription.FieldCancelled,
		})
	}
	if value, ok := suo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: subscription.FieldCreatedAt,
		})
	}
	if suo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   subscription.UserTable,
			Columns: []string{subscription.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   subscription.UserTable,
			Columns: []string{subscription.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Subscription{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{subscription.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
