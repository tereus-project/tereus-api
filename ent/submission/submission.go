// Code generated by entc, DO NOT EDIT.

package submission

import (
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
	// FieldCompleted holds the string denoting the completed field in the database.
	FieldCompleted = "completed"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// Table holds the table name of the submission in the database.
	Table = "submissions"
)

// Columns holds all SQL columns for submission fields.
var Columns = []string{
	FieldID,
	FieldSourceLanguage,
	FieldTargetLanguage,
	FieldCompleted,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCompleted holds the default value on creation for the "completed" field.
	DefaultCompleted bool
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)