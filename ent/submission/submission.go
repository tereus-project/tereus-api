// Code generated by entc, DO NOT EDIT.

package submission

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the submission type in the database.
	Label = "submission"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldSourceLanguage holds the string denoting the source_language field in the database.
	FieldSourceLanguage = "source_language"
	// FieldTargetLanguage holds the string denoting the target_language field in the database.
	FieldTargetLanguage = "target_language"
	// FieldIsInline holds the string denoting the is_inline field in the database.
	FieldIsInline = "is_inline"
	// FieldIsPublic holds the string denoting the is_public field in the database.
	FieldIsPublic = "is_public"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldReason holds the string denoting the reason field in the database.
	FieldReason = "reason"
	// FieldGitRepo holds the string denoting the git_repo field in the database.
	FieldGitRepo = "git_repo"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldShareID holds the string denoting the share_id field in the database.
	FieldShareID = "share_id"
	// FieldSubmissionSourceSizeBytes holds the string denoting the submission_source_size_bytes field in the database.
	FieldSubmissionSourceSizeBytes = "submission_source_size_bytes"
	// FieldSubmissionTargetSizeBytes holds the string denoting the submission_target_size_bytes field in the database.
	FieldSubmissionTargetSizeBytes = "submission_target_size_bytes"
	// FieldProcessingStartedAt holds the string denoting the processing_started_at field in the database.
	FieldProcessingStartedAt = "processing_started_at"
	// FieldProcessingFinishedAt holds the string denoting the processing_finished_at field in the database.
	FieldProcessingFinishedAt = "processing_finished_at"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// Table holds the table name of the submission in the database.
	Table = "submissions"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "submissions"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_submissions"
)

// Columns holds all SQL columns for submission fields.
var Columns = []string{
	FieldID,
	FieldSourceLanguage,
	FieldTargetLanguage,
	FieldIsInline,
	FieldIsPublic,
	FieldStatus,
	FieldReason,
	FieldGitRepo,
	FieldCreatedAt,
	FieldShareID,
	FieldSubmissionSourceSizeBytes,
	FieldSubmissionTargetSizeBytes,
	FieldProcessingStartedAt,
	FieldProcessingFinishedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "submissions"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_submissions",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultIsInline holds the default value on creation for the "is_inline" field.
	DefaultIsInline bool
	// DefaultIsPublic holds the default value on creation for the "is_public" field.
	DefaultIsPublic bool
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultSubmissionSourceSizeBytes holds the default value on creation for the "submission_source_size_bytes" field.
	DefaultSubmissionSourceSizeBytes int
	// DefaultSubmissionTargetSizeBytes holds the default value on creation for the "submission_target_size_bytes" field.
	DefaultSubmissionTargetSizeBytes int
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// Status defines the type for the "status" enum field.
type Status string

// StatusPending is the default value of the Status enum.
const DefaultStatus = StatusPending

// Status values.
const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusDone       Status = "done"
	StatusFailed     Status = "failed"
	StatusCleaned    Status = "cleaned"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusPending, StatusProcessing, StatusDone, StatusFailed, StatusCleaned:
		return nil
	default:
		return fmt.Errorf("submission: invalid enum value for status field: %q", s)
	}
}
