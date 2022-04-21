// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// SubmissionsColumns holds the columns for the "submissions" table.
	SubmissionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "source_language", Type: field.TypeString},
		{Name: "target_language", Type: field.TypeString},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"pending", "processing", "done", "failed"}, Default: "pending"},
		{Name: "git_repo", Type: field.TypeString, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "user_submissions", Type: field.TypeUUID},
	}
	// SubmissionsTable holds the schema information for the "submissions" table.
	SubmissionsTable = &schema.Table{
		Name:       "submissions",
		Columns:    SubmissionsColumns,
		PrimaryKey: []*schema.Column{SubmissionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "submissions_users_submissions",
				Columns:    []*schema.Column{SubmissionsColumns[6]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// TokensColumns holds the columns for the "tokens" table.
	TokensColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "is_active", Type: field.TypeBool, Default: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "user_tokens", Type: field.TypeUUID},
	}
	// TokensTable holds the schema information for the "tokens" table.
	TokensTable = &schema.Table{
		Name:       "tokens",
		Columns:    TokensColumns,
		PrimaryKey: []*schema.Column{TokensColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "tokens_users_tokens",
				Columns:    []*schema.Column{TokensColumns[3]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "email", Type: field.TypeString, Unique: true},
		{Name: "password", Type: field.TypeString, Nullable: true},
		{Name: "github_access_token", Type: field.TypeString, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		SubmissionsTable,
		TokensTable,
		UsersTable,
	}
)

func init() {
	SubmissionsTable.ForeignKeys[0].RefTable = UsersTable
	TokensTable.ForeignKeys[0].RefTable = UsersTable
}
