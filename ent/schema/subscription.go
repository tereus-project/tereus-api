package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Subscription holds the schema definition for the Subscription entity.
type Subscription struct {
	ent.Schema
}

// Fields of the Subscription.
func (Subscription) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("stripe_customer_id").Optional(),
		field.String("stripe_subscription_id").Optional(),
		field.Enum("tier").Values("free", "pro", "enterprise").Default("free"),
		field.Time("expires_at").Optional(),
		field.Bool("cancelled").Default(false),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Subscription.
func (Subscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("subscription").
			Unique().
			Required(),
	}
}
