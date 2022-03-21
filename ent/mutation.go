// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tereus-project/tereus-api/ent/predicate"
	"github.com/tereus-project/tereus-api/ent/submission"

	"entgo.io/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeSubmission = "Submission"
)

// SubmissionMutation represents an operation that mutates the Submission nodes in the graph.
type SubmissionMutation struct {
	config
	op              Op
	typ             string
	id              *uuid.UUID
	source_language *string
	target_language *string
	created_at      *time.Time
	clearedFields   map[string]struct{}
	done            bool
	oldValue        func(context.Context) (*Submission, error)
	predicates      []predicate.Submission
}

var _ ent.Mutation = (*SubmissionMutation)(nil)

// submissionOption allows management of the mutation configuration using functional options.
type submissionOption func(*SubmissionMutation)

// newSubmissionMutation creates new mutation for the Submission entity.
func newSubmissionMutation(c config, op Op, opts ...submissionOption) *SubmissionMutation {
	m := &SubmissionMutation{
		config:        c,
		op:            op,
		typ:           TypeSubmission,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withSubmissionID sets the ID field of the mutation.
func withSubmissionID(id uuid.UUID) submissionOption {
	return func(m *SubmissionMutation) {
		var (
			err   error
			once  sync.Once
			value *Submission
		)
		m.oldValue = func(ctx context.Context) (*Submission, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Submission.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withSubmission sets the old Submission of the mutation.
func withSubmission(node *Submission) submissionOption {
	return func(m *SubmissionMutation) {
		m.oldValue = func(context.Context) (*Submission, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SubmissionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SubmissionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Submission entities.
func (m *SubmissionMutation) SetID(id uuid.UUID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *SubmissionMutation) ID() (id uuid.UUID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *SubmissionMutation) IDs(ctx context.Context) ([]uuid.UUID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uuid.UUID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Submission.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetSourceLanguage sets the "source_language" field.
func (m *SubmissionMutation) SetSourceLanguage(s string) {
	m.source_language = &s
}

// SourceLanguage returns the value of the "source_language" field in the mutation.
func (m *SubmissionMutation) SourceLanguage() (r string, exists bool) {
	v := m.source_language
	if v == nil {
		return
	}
	return *v, true
}

// OldSourceLanguage returns the old "source_language" field's value of the Submission entity.
// If the Submission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubmissionMutation) OldSourceLanguage(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSourceLanguage is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSourceLanguage requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSourceLanguage: %w", err)
	}
	return oldValue.SourceLanguage, nil
}

// ResetSourceLanguage resets all changes to the "source_language" field.
func (m *SubmissionMutation) ResetSourceLanguage() {
	m.source_language = nil
}

// SetTargetLanguage sets the "target_language" field.
func (m *SubmissionMutation) SetTargetLanguage(s string) {
	m.target_language = &s
}

// TargetLanguage returns the value of the "target_language" field in the mutation.
func (m *SubmissionMutation) TargetLanguage() (r string, exists bool) {
	v := m.target_language
	if v == nil {
		return
	}
	return *v, true
}

// OldTargetLanguage returns the old "target_language" field's value of the Submission entity.
// If the Submission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubmissionMutation) OldTargetLanguage(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldTargetLanguage is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldTargetLanguage requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTargetLanguage: %w", err)
	}
	return oldValue.TargetLanguage, nil
}

// ResetTargetLanguage resets all changes to the "target_language" field.
func (m *SubmissionMutation) ResetTargetLanguage() {
	m.target_language = nil
}

// SetCreatedAt sets the "created_at" field.
func (m *SubmissionMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *SubmissionMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the Submission entity.
// If the Submission object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubmissionMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *SubmissionMutation) ResetCreatedAt() {
	m.created_at = nil
}

// Where appends a list predicates to the SubmissionMutation builder.
func (m *SubmissionMutation) Where(ps ...predicate.Submission) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *SubmissionMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Submission).
func (m *SubmissionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *SubmissionMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.source_language != nil {
		fields = append(fields, submission.FieldSourceLanguage)
	}
	if m.target_language != nil {
		fields = append(fields, submission.FieldTargetLanguage)
	}
	if m.created_at != nil {
		fields = append(fields, submission.FieldCreatedAt)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *SubmissionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case submission.FieldSourceLanguage:
		return m.SourceLanguage()
	case submission.FieldTargetLanguage:
		return m.TargetLanguage()
	case submission.FieldCreatedAt:
		return m.CreatedAt()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *SubmissionMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case submission.FieldSourceLanguage:
		return m.OldSourceLanguage(ctx)
	case submission.FieldTargetLanguage:
		return m.OldTargetLanguage(ctx)
	case submission.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	}
	return nil, fmt.Errorf("unknown Submission field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SubmissionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case submission.FieldSourceLanguage:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSourceLanguage(v)
		return nil
	case submission.FieldTargetLanguage:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTargetLanguage(v)
		return nil
	case submission.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	}
	return fmt.Errorf("unknown Submission field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *SubmissionMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *SubmissionMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SubmissionMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Submission numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *SubmissionMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *SubmissionMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *SubmissionMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Submission nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *SubmissionMutation) ResetField(name string) error {
	switch name {
	case submission.FieldSourceLanguage:
		m.ResetSourceLanguage()
		return nil
	case submission.FieldTargetLanguage:
		m.ResetTargetLanguage()
		return nil
	case submission.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	}
	return fmt.Errorf("unknown Submission field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *SubmissionMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *SubmissionMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *SubmissionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *SubmissionMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *SubmissionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *SubmissionMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *SubmissionMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Submission unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *SubmissionMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Submission edge %s", name)
}
