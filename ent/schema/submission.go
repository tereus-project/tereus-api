package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Submission holds the schema definition for the Submission entity.
type Submission struct {
	ent.Schema
}

// Fields of the Submission.
func (Submission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("source_language"),
		field.String("target_language"),
		field.Bool("is_inline").Default(false),
		field.Bool("is_public").Default(false),
		field.Enum("status").Values("pending", "processing", "done", "failed", "cleaned").Default("pending"),
		field.String("reason").Optional(),
		field.String("git_repo").Optional(),
		field.Time("created_at").Default(time.Now),
		field.String("share_id").Optional().Unique(),
		field.Int("submission_source_size_bytes").Default(0),
		field.Int("submission_target_size_bytes").Default(0),
	}
}

// Edges of the Submission.
func (Submission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("submissions").
			Unique().
			Required(),
	}
}
