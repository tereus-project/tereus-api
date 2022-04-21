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
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/ent/user"
)

// SubmissionUpdate is the builder for updating Submission entities.
type SubmissionUpdate struct {
	config
	hooks    []Hook
	mutation *SubmissionMutation
}

// Where appends a list predicates to the SubmissionUpdate builder.
func (su *SubmissionUpdate) Where(ps ...predicate.Submission) *SubmissionUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetSourceLanguage sets the "source_language" field.
func (su *SubmissionUpdate) SetSourceLanguage(s string) *SubmissionUpdate {
	su.mutation.SetSourceLanguage(s)
	return su
}

// SetTargetLanguage sets the "target_language" field.
func (su *SubmissionUpdate) SetTargetLanguage(s string) *SubmissionUpdate {
	su.mutation.SetTargetLanguage(s)
	return su
}

// SetStatus sets the "status" field.
func (su *SubmissionUpdate) SetStatus(s submission.Status) *SubmissionUpdate {
	su.mutation.SetStatus(s)
	return su
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (su *SubmissionUpdate) SetNillableStatus(s *submission.Status) *SubmissionUpdate {
	if s != nil {
		su.SetStatus(*s)
	}
	return su
}

// SetGitRepo sets the "git_repo" field.
func (su *SubmissionUpdate) SetGitRepo(s string) *SubmissionUpdate {
	su.mutation.SetGitRepo(s)
	return su
}

// SetNillableGitRepo sets the "git_repo" field if the given value is not nil.
func (su *SubmissionUpdate) SetNillableGitRepo(s *string) *SubmissionUpdate {
	if s != nil {
		su.SetGitRepo(*s)
	}
	return su
}

// ClearGitRepo clears the value of the "git_repo" field.
func (su *SubmissionUpdate) ClearGitRepo() *SubmissionUpdate {
	su.mutation.ClearGitRepo()
	return su
}

// SetCreatedAt sets the "created_at" field.
func (su *SubmissionUpdate) SetCreatedAt(t time.Time) *SubmissionUpdate {
	su.mutation.SetCreatedAt(t)
	return su
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (su *SubmissionUpdate) SetNillableCreatedAt(t *time.Time) *SubmissionUpdate {
	if t != nil {
		su.SetCreatedAt(*t)
	}
	return su
}

// SetUserID sets the "user" edge to the User entity by ID.
func (su *SubmissionUpdate) SetUserID(id uuid.UUID) *SubmissionUpdate {
	su.mutation.SetUserID(id)
	return su
}

// SetUser sets the "user" edge to the User entity.
func (su *SubmissionUpdate) SetUser(u *User) *SubmissionUpdate {
	return su.SetUserID(u.ID)
}

// Mutation returns the SubmissionMutation object of the builder.
func (su *SubmissionUpdate) Mutation() *SubmissionMutation {
	return su.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (su *SubmissionUpdate) ClearUser() *SubmissionUpdate {
	su.mutation.ClearUser()
	return su
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SubmissionUpdate) Save(ctx context.Context) (int, error) {
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
			mutation, ok := m.(*SubmissionMutation)
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
func (su *SubmissionUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SubmissionUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SubmissionUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *SubmissionUpdate) check() error {
	if v, ok := su.mutation.Status(); ok {
		if err := submission.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Submission.status": %w`, err)}
		}
	}
	if _, ok := su.mutation.UserID(); su.mutation.UserCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Submission.user"`)
	}
	return nil
}

func (su *SubmissionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   submission.Table,
			Columns: submission.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: submission.FieldID,
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
	if value, ok := su.mutation.SourceLanguage(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: submission.FieldSourceLanguage,
		})
	}
	if value, ok := su.mutation.TargetLanguage(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: submission.FieldTargetLanguage,
		})
	}
	if value, ok := su.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: submission.FieldStatus,
		})
	}
	if value, ok := su.mutation.GitRepo(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: submission.FieldGitRepo,
		})
	}
	if su.mutation.GitRepoCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: submission.FieldGitRepo,
		})
	}
	if value, ok := su.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: submission.FieldCreatedAt,
		})
	}
	if su.mutation.UserCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.UserIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{submission.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// SubmissionUpdateOne is the builder for updating a single Submission entity.
type SubmissionUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SubmissionMutation
}

// SetSourceLanguage sets the "source_language" field.
func (suo *SubmissionUpdateOne) SetSourceLanguage(s string) *SubmissionUpdateOne {
	suo.mutation.SetSourceLanguage(s)
	return suo
}

// SetTargetLanguage sets the "target_language" field.
func (suo *SubmissionUpdateOne) SetTargetLanguage(s string) *SubmissionUpdateOne {
	suo.mutation.SetTargetLanguage(s)
	return suo
}

// SetStatus sets the "status" field.
func (suo *SubmissionUpdateOne) SetStatus(s submission.Status) *SubmissionUpdateOne {
	suo.mutation.SetStatus(s)
	return suo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (suo *SubmissionUpdateOne) SetNillableStatus(s *submission.Status) *SubmissionUpdateOne {
	if s != nil {
		suo.SetStatus(*s)
	}
	return suo
}

// SetGitRepo sets the "git_repo" field.
func (suo *SubmissionUpdateOne) SetGitRepo(s string) *SubmissionUpdateOne {
	suo.mutation.SetGitRepo(s)
	return suo
}

// SetNillableGitRepo sets the "git_repo" field if the given value is not nil.
func (suo *SubmissionUpdateOne) SetNillableGitRepo(s *string) *SubmissionUpdateOne {
	if s != nil {
		suo.SetGitRepo(*s)
	}
	return suo
}

// ClearGitRepo clears the value of the "git_repo" field.
func (suo *SubmissionUpdateOne) ClearGitRepo() *SubmissionUpdateOne {
	suo.mutation.ClearGitRepo()
	return suo
}

// SetCreatedAt sets the "created_at" field.
func (suo *SubmissionUpdateOne) SetCreatedAt(t time.Time) *SubmissionUpdateOne {
	suo.mutation.SetCreatedAt(t)
	return suo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (suo *SubmissionUpdateOne) SetNillableCreatedAt(t *time.Time) *SubmissionUpdateOne {
	if t != nil {
		suo.SetCreatedAt(*t)
	}
	return suo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (suo *SubmissionUpdateOne) SetUserID(id uuid.UUID) *SubmissionUpdateOne {
	suo.mutation.SetUserID(id)
	return suo
}

// SetUser sets the "user" edge to the User entity.
func (suo *SubmissionUpdateOne) SetUser(u *User) *SubmissionUpdateOne {
	return suo.SetUserID(u.ID)
}

// Mutation returns the SubmissionMutation object of the builder.
func (suo *SubmissionUpdateOne) Mutation() *SubmissionMutation {
	return suo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (suo *SubmissionUpdateOne) ClearUser() *SubmissionUpdateOne {
	suo.mutation.ClearUser()
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SubmissionUpdateOne) Select(field string, fields ...string) *SubmissionUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Submission entity.
func (suo *SubmissionUpdateOne) Save(ctx context.Context) (*Submission, error) {
	var (
		err  error
		node *Submission
	)
	if len(suo.hooks) == 0 {
		if err = suo.check(); err != nil {
			return nil, err
		}
		node, err = suo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SubmissionMutation)
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
func (suo *SubmissionUpdateOne) SaveX(ctx context.Context) *Submission {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SubmissionUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SubmissionUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *SubmissionUpdateOne) check() error {
	if v, ok := suo.mutation.Status(); ok {
		if err := submission.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Submission.status": %w`, err)}
		}
	}
	if _, ok := suo.mutation.UserID(); suo.mutation.UserCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Submission.user"`)
	}
	return nil
}

func (suo *SubmissionUpdateOne) sqlSave(ctx context.Context) (_node *Submission, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   submission.Table,
			Columns: submission.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: submission.FieldID,
			},
		},
	}
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Submission.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, submission.FieldID)
		for _, f := range fields {
			if !submission.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != submission.FieldID {
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
	if value, ok := suo.mutation.SourceLanguage(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: submission.FieldSourceLanguage,
		})
	}
	if value, ok := suo.mutation.TargetLanguage(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: submission.FieldTargetLanguage,
		})
	}
	if value, ok := suo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: submission.FieldStatus,
		})
	}
	if value, ok := suo.mutation.GitRepo(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: submission.FieldGitRepo,
		})
	}
	if suo.mutation.GitRepoCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: submission.FieldGitRepo,
		})
	}
	if value, ok := suo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: submission.FieldCreatedAt,
		})
	}
	if suo.mutation.UserCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.UserIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Submission{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{submission.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
