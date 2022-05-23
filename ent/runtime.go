// Code generated by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/tereus-project/tereus-api/ent/schema"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/ent/subscription"
	"github.com/tereus-project/tereus-api/ent/token"
	"github.com/tereus-project/tereus-api/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	submissionFields := schema.Submission{}.Fields()
	_ = submissionFields
	// submissionDescCreatedAt is the schema descriptor for created_at field.
	submissionDescCreatedAt := submissionFields[6].Descriptor()
	// submission.DefaultCreatedAt holds the default value on creation for the created_at field.
	submission.DefaultCreatedAt = submissionDescCreatedAt.Default.(func() time.Time)
	// submissionDescID is the schema descriptor for id field.
	submissionDescID := submissionFields[0].Descriptor()
	// submission.DefaultID holds the default value on creation for the id field.
	submission.DefaultID = submissionDescID.Default.(func() uuid.UUID)
	subscriptionFields := schema.Subscription{}.Fields()
	_ = subscriptionFields
	// subscriptionDescCreatedAt is the schema descriptor for created_at field.
	subscriptionDescCreatedAt := subscriptionFields[5].Descriptor()
	// subscription.DefaultCreatedAt holds the default value on creation for the created_at field.
	subscription.DefaultCreatedAt = subscriptionDescCreatedAt.Default.(func() time.Time)
	// subscriptionDescID is the schema descriptor for id field.
	subscriptionDescID := subscriptionFields[0].Descriptor()
	// subscription.DefaultID holds the default value on creation for the id field.
	subscription.DefaultID = subscriptionDescID.Default.(func() uuid.UUID)
	tokenFields := schema.Token{}.Fields()
	_ = tokenFields
	// tokenDescIsActive is the schema descriptor for is_active field.
	tokenDescIsActive := tokenFields[1].Descriptor()
	// token.DefaultIsActive holds the default value on creation for the is_active field.
	token.DefaultIsActive = tokenDescIsActive.Default.(bool)
	// tokenDescCreatedAt is the schema descriptor for created_at field.
	tokenDescCreatedAt := tokenFields[2].Descriptor()
	// token.DefaultCreatedAt holds the default value on creation for the created_at field.
	token.DefaultCreatedAt = tokenDescCreatedAt.Default.(func() time.Time)
	// tokenDescID is the schema descriptor for id field.
	tokenDescID := tokenFields[0].Descriptor()
	// token.DefaultID holds the default value on creation for the id field.
	token.DefaultID = tokenDescID.Default.(func() uuid.UUID)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userFields[4].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
}
