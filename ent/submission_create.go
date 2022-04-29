// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/ent/user"
)

// SubmissionCreate is the builder for creating a Submission entity.
type SubmissionCreate struct {
	config
	mutation *SubmissionMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetSourceLanguage sets the "source_language" field.
func (sc *SubmissionCreate) SetSourceLanguage(s string) *SubmissionCreate {
	sc.mutation.SetSourceLanguage(s)
	return sc
}

// SetTargetLanguage sets the "target_language" field.
func (sc *SubmissionCreate) SetTargetLanguage(s string) *SubmissionCreate {
	sc.mutation.SetTargetLanguage(s)
	return sc
}

// SetStatus sets the "status" field.
func (sc *SubmissionCreate) SetStatus(s submission.Status) *SubmissionCreate {
	sc.mutation.SetStatus(s)
	return sc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (sc *SubmissionCreate) SetNillableStatus(s *submission.Status) *SubmissionCreate {
	if s != nil {
		sc.SetStatus(*s)
	}
	return sc
}

// SetReason sets the "reason" field.
func (sc *SubmissionCreate) SetReason(s string) *SubmissionCreate {
	sc.mutation.SetReason(s)
	return sc
}

// SetNillableReason sets the "reason" field if the given value is not nil.
func (sc *SubmissionCreate) SetNillableReason(s *string) *SubmissionCreate {
	if s != nil {
		sc.SetReason(*s)
	}
	return sc
}

// SetGitRepo sets the "git_repo" field.
func (sc *SubmissionCreate) SetGitRepo(s string) *SubmissionCreate {
	sc.mutation.SetGitRepo(s)
	return sc
}

// SetNillableGitRepo sets the "git_repo" field if the given value is not nil.
func (sc *SubmissionCreate) SetNillableGitRepo(s *string) *SubmissionCreate {
	if s != nil {
		sc.SetGitRepo(*s)
	}
	return sc
}

// SetCreatedAt sets the "created_at" field.
func (sc *SubmissionCreate) SetCreatedAt(t time.Time) *SubmissionCreate {
	sc.mutation.SetCreatedAt(t)
	return sc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (sc *SubmissionCreate) SetNillableCreatedAt(t *time.Time) *SubmissionCreate {
	if t != nil {
		sc.SetCreatedAt(*t)
	}
	return sc
}

// SetID sets the "id" field.
func (sc *SubmissionCreate) SetID(u uuid.UUID) *SubmissionCreate {
	sc.mutation.SetID(u)
	return sc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (sc *SubmissionCreate) SetNillableID(u *uuid.UUID) *SubmissionCreate {
	if u != nil {
		sc.SetID(*u)
	}
	return sc
}

// SetUserID sets the "user" edge to the User entity by ID.
func (sc *SubmissionCreate) SetUserID(id uuid.UUID) *SubmissionCreate {
	sc.mutation.SetUserID(id)
	return sc
}

// SetUser sets the "user" edge to the User entity.
func (sc *SubmissionCreate) SetUser(u *User) *SubmissionCreate {
	return sc.SetUserID(u.ID)
}

// Mutation returns the SubmissionMutation object of the builder.
func (sc *SubmissionCreate) Mutation() *SubmissionMutation {
	return sc.mutation
}

// Save creates the Submission in the database.
func (sc *SubmissionCreate) Save(ctx context.Context) (*Submission, error) {
	var (
		err  error
		node *Submission
	)
	sc.defaults()
	if len(sc.hooks) == 0 {
		if err = sc.check(); err != nil {
			return nil, err
		}
		node, err = sc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SubmissionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = sc.check(); err != nil {
				return nil, err
			}
			sc.mutation = mutation
			if node, err = sc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(sc.hooks) - 1; i >= 0; i-- {
			if sc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = sc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SubmissionCreate) SaveX(ctx context.Context) *Submission {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SubmissionCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SubmissionCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sc *SubmissionCreate) defaults() {
	if _, ok := sc.mutation.Status(); !ok {
		v := submission.DefaultStatus
		sc.mutation.SetStatus(v)
	}
	if _, ok := sc.mutation.CreatedAt(); !ok {
		v := submission.DefaultCreatedAt()
		sc.mutation.SetCreatedAt(v)
	}
	if _, ok := sc.mutation.ID(); !ok {
		v := submission.DefaultID()
		sc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SubmissionCreate) check() error {
	if _, ok := sc.mutation.SourceLanguage(); !ok {
		return &ValidationError{Name: "source_language", err: errors.New(`ent: missing required field "Submission.source_language"`)}
	}
	if _, ok := sc.mutation.TargetLanguage(); !ok {
		return &ValidationError{Name: "target_language", err: errors.New(`ent: missing required field "Submission.target_language"`)}
	}
	if _, ok := sc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "Submission.status"`)}
	}
	if v, ok := sc.mutation.Status(); ok {
		if err := submission.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Submission.status": %w`, err)}
		}
	}
	if _, ok := sc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Submission.created_at"`)}
	}
	if _, ok := sc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user", err: errors.New(`ent: missing required edge "Submission.user"`)}
	}
	return nil
}

func (sc *SubmissionCreate) sqlSave(ctx context.Context) (*Submission, error) {
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	return _node, nil
}

func (sc *SubmissionCreate) createSpec() (*Submission, *sqlgraph.CreateSpec) {
	var (
		_node = &Submission{config: sc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: submission.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: submission.FieldID,
			},
		}
	)
	_spec.OnConflict = sc.conflict
	if id, ok := sc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := sc.mutation.SourceLanguage(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: submission.FieldSourceLanguage,
		})
		_node.SourceLanguage = value
	}
	if value, ok := sc.mutation.TargetLanguage(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: submission.FieldTargetLanguage,
		})
		_node.TargetLanguage = value
	}
	if value, ok := sc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: submission.FieldStatus,
		})
		_node.Status = value
	}
	if value, ok := sc.mutation.Reason(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: submission.FieldReason,
		})
		_node.Reason = value
	}
	if value, ok := sc.mutation.GitRepo(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: submission.FieldGitRepo,
		})
		_node.GitRepo = value
	}
	if value, ok := sc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: submission.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if nodes := sc.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   submission.UserTable,
			Columns: []string{submission.UserColumn},
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
		_node.user_submissions = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Submission.Create().
//		SetSourceLanguage(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SubmissionUpsert) {
//			SetSourceLanguage(v+v).
//		}).
//		Exec(ctx)
//
func (sc *SubmissionCreate) OnConflict(opts ...sql.ConflictOption) *SubmissionUpsertOne {
	sc.conflict = opts
	return &SubmissionUpsertOne{
		create: sc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Submission.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (sc *SubmissionCreate) OnConflictColumns(columns ...string) *SubmissionUpsertOne {
	sc.conflict = append(sc.conflict, sql.ConflictColumns(columns...))
	return &SubmissionUpsertOne{
		create: sc,
	}
}

type (
	// SubmissionUpsertOne is the builder for "upsert"-ing
	//  one Submission node.
	SubmissionUpsertOne struct {
		create *SubmissionCreate
	}

	// SubmissionUpsert is the "OnConflict" setter.
	SubmissionUpsert struct {
		*sql.UpdateSet
	}
)

// SetSourceLanguage sets the "source_language" field.
func (u *SubmissionUpsert) SetSourceLanguage(v string) *SubmissionUpsert {
	u.Set(submission.FieldSourceLanguage, v)
	return u
}

// UpdateSourceLanguage sets the "source_language" field to the value that was provided on create.
func (u *SubmissionUpsert) UpdateSourceLanguage() *SubmissionUpsert {
	u.SetExcluded(submission.FieldSourceLanguage)
	return u
}

// SetTargetLanguage sets the "target_language" field.
func (u *SubmissionUpsert) SetTargetLanguage(v string) *SubmissionUpsert {
	u.Set(submission.FieldTargetLanguage, v)
	return u
}

// UpdateTargetLanguage sets the "target_language" field to the value that was provided on create.
func (u *SubmissionUpsert) UpdateTargetLanguage() *SubmissionUpsert {
	u.SetExcluded(submission.FieldTargetLanguage)
	return u
}

// SetStatus sets the "status" field.
func (u *SubmissionUpsert) SetStatus(v submission.Status) *SubmissionUpsert {
	u.Set(submission.FieldStatus, v)
	return u
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *SubmissionUpsert) UpdateStatus() *SubmissionUpsert {
	u.SetExcluded(submission.FieldStatus)
	return u
}

// SetReason sets the "reason" field.
func (u *SubmissionUpsert) SetReason(v string) *SubmissionUpsert {
	u.Set(submission.FieldReason, v)
	return u
}

// UpdateReason sets the "reason" field to the value that was provided on create.
func (u *SubmissionUpsert) UpdateReason() *SubmissionUpsert {
	u.SetExcluded(submission.FieldReason)
	return u
}

// ClearReason clears the value of the "reason" field.
func (u *SubmissionUpsert) ClearReason() *SubmissionUpsert {
	u.SetNull(submission.FieldReason)
	return u
}

// SetGitRepo sets the "git_repo" field.
func (u *SubmissionUpsert) SetGitRepo(v string) *SubmissionUpsert {
	u.Set(submission.FieldGitRepo, v)
	return u
}

// UpdateGitRepo sets the "git_repo" field to the value that was provided on create.
func (u *SubmissionUpsert) UpdateGitRepo() *SubmissionUpsert {
	u.SetExcluded(submission.FieldGitRepo)
	return u
}

// ClearGitRepo clears the value of the "git_repo" field.
func (u *SubmissionUpsert) ClearGitRepo() *SubmissionUpsert {
	u.SetNull(submission.FieldGitRepo)
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *SubmissionUpsert) SetCreatedAt(v time.Time) *SubmissionUpsert {
	u.Set(submission.FieldCreatedAt, v)
	return u
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *SubmissionUpsert) UpdateCreatedAt() *SubmissionUpsert {
	u.SetExcluded(submission.FieldCreatedAt)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Submission.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(submission.FieldID)
//			}),
//		).
//		Exec(ctx)
//
func (u *SubmissionUpsertOne) UpdateNewValues() *SubmissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(submission.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.Submission.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *SubmissionUpsertOne) Ignore() *SubmissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SubmissionUpsertOne) DoNothing() *SubmissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SubmissionCreate.OnConflict
// documentation for more info.
func (u *SubmissionUpsertOne) Update(set func(*SubmissionUpsert)) *SubmissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SubmissionUpsert{UpdateSet: update})
	}))
	return u
}

// SetSourceLanguage sets the "source_language" field.
func (u *SubmissionUpsertOne) SetSourceLanguage(v string) *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetSourceLanguage(v)
	})
}

// UpdateSourceLanguage sets the "source_language" field to the value that was provided on create.
func (u *SubmissionUpsertOne) UpdateSourceLanguage() *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateSourceLanguage()
	})
}

// SetTargetLanguage sets the "target_language" field.
func (u *SubmissionUpsertOne) SetTargetLanguage(v string) *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetTargetLanguage(v)
	})
}

// UpdateTargetLanguage sets the "target_language" field to the value that was provided on create.
func (u *SubmissionUpsertOne) UpdateTargetLanguage() *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateTargetLanguage()
	})
}

// SetStatus sets the "status" field.
func (u *SubmissionUpsertOne) SetStatus(v submission.Status) *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *SubmissionUpsertOne) UpdateStatus() *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateStatus()
	})
}

// SetReason sets the "reason" field.
func (u *SubmissionUpsertOne) SetReason(v string) *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetReason(v)
	})
}

// UpdateReason sets the "reason" field to the value that was provided on create.
func (u *SubmissionUpsertOne) UpdateReason() *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateReason()
	})
}

// ClearReason clears the value of the "reason" field.
func (u *SubmissionUpsertOne) ClearReason() *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.ClearReason()
	})
}

// SetGitRepo sets the "git_repo" field.
func (u *SubmissionUpsertOne) SetGitRepo(v string) *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetGitRepo(v)
	})
}

// UpdateGitRepo sets the "git_repo" field to the value that was provided on create.
func (u *SubmissionUpsertOne) UpdateGitRepo() *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateGitRepo()
	})
}

// ClearGitRepo clears the value of the "git_repo" field.
func (u *SubmissionUpsertOne) ClearGitRepo() *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.ClearGitRepo()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *SubmissionUpsertOne) SetCreatedAt(v time.Time) *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *SubmissionUpsertOne) UpdateCreatedAt() *SubmissionUpsertOne {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateCreatedAt()
	})
}

// Exec executes the query.
func (u *SubmissionUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SubmissionCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SubmissionUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *SubmissionUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: SubmissionUpsertOne.ID is not supported by MySQL driver. Use SubmissionUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *SubmissionUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// SubmissionCreateBulk is the builder for creating many Submission entities in bulk.
type SubmissionCreateBulk struct {
	config
	builders []*SubmissionCreate
	conflict []sql.ConflictOption
}

// Save creates the Submission entities in the database.
func (scb *SubmissionCreateBulk) Save(ctx context.Context) ([]*Submission, error) {
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Submission, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SubmissionMutation)
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
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = scb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
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
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SubmissionCreateBulk) SaveX(ctx context.Context) []*Submission {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SubmissionCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SubmissionCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Submission.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SubmissionUpsert) {
//			SetSourceLanguage(v+v).
//		}).
//		Exec(ctx)
//
func (scb *SubmissionCreateBulk) OnConflict(opts ...sql.ConflictOption) *SubmissionUpsertBulk {
	scb.conflict = opts
	return &SubmissionUpsertBulk{
		create: scb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Submission.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (scb *SubmissionCreateBulk) OnConflictColumns(columns ...string) *SubmissionUpsertBulk {
	scb.conflict = append(scb.conflict, sql.ConflictColumns(columns...))
	return &SubmissionUpsertBulk{
		create: scb,
	}
}

// SubmissionUpsertBulk is the builder for "upsert"-ing
// a bulk of Submission nodes.
type SubmissionUpsertBulk struct {
	create *SubmissionCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Submission.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(submission.FieldID)
//			}),
//		).
//		Exec(ctx)
//
func (u *SubmissionUpsertBulk) UpdateNewValues() *SubmissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(submission.FieldID)
				return
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Submission.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
//
func (u *SubmissionUpsertBulk) Ignore() *SubmissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SubmissionUpsertBulk) DoNothing() *SubmissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SubmissionCreateBulk.OnConflict
// documentation for more info.
func (u *SubmissionUpsertBulk) Update(set func(*SubmissionUpsert)) *SubmissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SubmissionUpsert{UpdateSet: update})
	}))
	return u
}

// SetSourceLanguage sets the "source_language" field.
func (u *SubmissionUpsertBulk) SetSourceLanguage(v string) *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetSourceLanguage(v)
	})
}

// UpdateSourceLanguage sets the "source_language" field to the value that was provided on create.
func (u *SubmissionUpsertBulk) UpdateSourceLanguage() *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateSourceLanguage()
	})
}

// SetTargetLanguage sets the "target_language" field.
func (u *SubmissionUpsertBulk) SetTargetLanguage(v string) *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetTargetLanguage(v)
	})
}

// UpdateTargetLanguage sets the "target_language" field to the value that was provided on create.
func (u *SubmissionUpsertBulk) UpdateTargetLanguage() *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateTargetLanguage()
	})
}

// SetStatus sets the "status" field.
func (u *SubmissionUpsertBulk) SetStatus(v submission.Status) *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *SubmissionUpsertBulk) UpdateStatus() *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateStatus()
	})
}

// SetReason sets the "reason" field.
func (u *SubmissionUpsertBulk) SetReason(v string) *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetReason(v)
	})
}

// UpdateReason sets the "reason" field to the value that was provided on create.
func (u *SubmissionUpsertBulk) UpdateReason() *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateReason()
	})
}

// ClearReason clears the value of the "reason" field.
func (u *SubmissionUpsertBulk) ClearReason() *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.ClearReason()
	})
}

// SetGitRepo sets the "git_repo" field.
func (u *SubmissionUpsertBulk) SetGitRepo(v string) *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetGitRepo(v)
	})
}

// UpdateGitRepo sets the "git_repo" field to the value that was provided on create.
func (u *SubmissionUpsertBulk) UpdateGitRepo() *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateGitRepo()
	})
}

// ClearGitRepo clears the value of the "git_repo" field.
func (u *SubmissionUpsertBulk) ClearGitRepo() *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.ClearGitRepo()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *SubmissionUpsertBulk) SetCreatedAt(v time.Time) *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *SubmissionUpsertBulk) UpdateCreatedAt() *SubmissionUpsertBulk {
	return u.Update(func(s *SubmissionUpsert) {
		s.UpdateCreatedAt()
	})
}

// Exec executes the query.
func (u *SubmissionUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the SubmissionCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SubmissionCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SubmissionUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
