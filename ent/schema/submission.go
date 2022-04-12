package schema

import (
	"time"

	"entgo.io/ent"
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
		field.Enum("status").Values("pending", "processing", "done", "failed").Default("pending"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Submission.
func (Submission) Edges() []ent.Edge {
	return nil
}
