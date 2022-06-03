// Code generated by entc, DO NOT EDIT.

package subscription

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/tereus-project/tereus-api/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// StripeCustomerID applies equality check predicate on the "stripe_customer_id" field. It's identical to StripeCustomerIDEQ.
func StripeCustomerID(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStripeCustomerID), v))
	})
}

// StripeSubscriptionID applies equality check predicate on the "stripe_subscription_id" field. It's identical to StripeSubscriptionIDEQ.
func StripeSubscriptionID(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStripeSubscriptionID), v))
	})
}

// ExpiresAt applies equality check predicate on the "expires_at" field. It's identical to ExpiresAtEQ.
func ExpiresAt(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExpiresAt), v))
	})
}

// Cancelled applies equality check predicate on the "cancelled" field. It's identical to CancelledEQ.
func Cancelled(v bool) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCancelled), v))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// StripeCustomerIDEQ applies the EQ predicate on the "stripe_customer_id" field.
func StripeCustomerIDEQ(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStripeCustomerID), v))
	})
}

// StripeCustomerIDNEQ applies the NEQ predicate on the "stripe_customer_id" field.
func StripeCustomerIDNEQ(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStripeCustomerID), v))
	})
}

// StripeCustomerIDIn applies the In predicate on the "stripe_customer_id" field.
func StripeCustomerIDIn(vs ...string) predicate.Subscription {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldStripeCustomerID), v...))
	})
}

// StripeCustomerIDNotIn applies the NotIn predicate on the "stripe_customer_id" field.
func StripeCustomerIDNotIn(vs ...string) predicate.Subscription {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldStripeCustomerID), v...))
	})
}

// StripeCustomerIDGT applies the GT predicate on the "stripe_customer_id" field.
func StripeCustomerIDGT(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStripeCustomerID), v))
	})
}

// StripeCustomerIDGTE applies the GTE predicate on the "stripe_customer_id" field.
func StripeCustomerIDGTE(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStripeCustomerID), v))
	})
}

// StripeCustomerIDLT applies the LT predicate on the "stripe_customer_id" field.
func StripeCustomerIDLT(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStripeCustomerID), v))
	})
}

// StripeCustomerIDLTE applies the LTE predicate on the "stripe_customer_id" field.
func StripeCustomerIDLTE(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStripeCustomerID), v))
	})
}

// StripeCustomerIDContains applies the Contains predicate on the "stripe_customer_id" field.
func StripeCustomerIDContains(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldStripeCustomerID), v))
	})
}

// StripeCustomerIDHasPrefix applies the HasPrefix predicate on the "stripe_customer_id" field.
func StripeCustomerIDHasPrefix(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldStripeCustomerID), v))
	})
}

// StripeCustomerIDHasSuffix applies the HasSuffix predicate on the "stripe_customer_id" field.
func StripeCustomerIDHasSuffix(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldStripeCustomerID), v))
	})
}

// StripeCustomerIDEqualFold applies the EqualFold predicate on the "stripe_customer_id" field.
func StripeCustomerIDEqualFold(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldStripeCustomerID), v))
	})
}

// StripeCustomerIDContainsFold applies the ContainsFold predicate on the "stripe_customer_id" field.
func StripeCustomerIDContainsFold(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldStripeCustomerID), v))
	})
}

// StripeSubscriptionIDEQ applies the EQ predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDEQ(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStripeSubscriptionID), v))
	})
}

// StripeSubscriptionIDNEQ applies the NEQ predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDNEQ(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStripeSubscriptionID), v))
	})
}

// StripeSubscriptionIDIn applies the In predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDIn(vs ...string) predicate.Subscription {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldStripeSubscriptionID), v...))
	})
}

// StripeSubscriptionIDNotIn applies the NotIn predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDNotIn(vs ...string) predicate.Subscription {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldStripeSubscriptionID), v...))
	})
}

// StripeSubscriptionIDGT applies the GT predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDGT(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStripeSubscriptionID), v))
	})
}

// StripeSubscriptionIDGTE applies the GTE predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDGTE(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStripeSubscriptionID), v))
	})
}

// StripeSubscriptionIDLT applies the LT predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDLT(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStripeSubscriptionID), v))
	})
}

// StripeSubscriptionIDLTE applies the LTE predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDLTE(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStripeSubscriptionID), v))
	})
}

// StripeSubscriptionIDContains applies the Contains predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDContains(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldStripeSubscriptionID), v))
	})
}

// StripeSubscriptionIDHasPrefix applies the HasPrefix predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDHasPrefix(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldStripeSubscriptionID), v))
	})
}

// StripeSubscriptionIDHasSuffix applies the HasSuffix predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDHasSuffix(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldStripeSubscriptionID), v))
	})
}

// StripeSubscriptionIDIsNil applies the IsNil predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDIsNil() predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldStripeSubscriptionID)))
	})
}

// StripeSubscriptionIDNotNil applies the NotNil predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDNotNil() predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldStripeSubscriptionID)))
	})
}

// StripeSubscriptionIDEqualFold applies the EqualFold predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDEqualFold(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldStripeSubscriptionID), v))
	})
}

// StripeSubscriptionIDContainsFold applies the ContainsFold predicate on the "stripe_subscription_id" field.
func StripeSubscriptionIDContainsFold(v string) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldStripeSubscriptionID), v))
	})
}

// TierEQ applies the EQ predicate on the "tier" field.
func TierEQ(v Tier) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTier), v))
	})
}

// TierNEQ applies the NEQ predicate on the "tier" field.
func TierNEQ(v Tier) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTier), v))
	})
}

// TierIn applies the In predicate on the "tier" field.
func TierIn(vs ...Tier) predicate.Subscription {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTier), v...))
	})
}

// TierNotIn applies the NotIn predicate on the "tier" field.
func TierNotIn(vs ...Tier) predicate.Subscription {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTier), v...))
	})
}

// ExpiresAtEQ applies the EQ predicate on the "expires_at" field.
func ExpiresAtEQ(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExpiresAt), v))
	})
}

// ExpiresAtNEQ applies the NEQ predicate on the "expires_at" field.
func ExpiresAtNEQ(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldExpiresAt), v))
	})
}

// ExpiresAtIn applies the In predicate on the "expires_at" field.
func ExpiresAtIn(vs ...time.Time) predicate.Subscription {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldExpiresAt), v...))
	})
}

// ExpiresAtNotIn applies the NotIn predicate on the "expires_at" field.
func ExpiresAtNotIn(vs ...time.Time) predicate.Subscription {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldExpiresAt), v...))
	})
}

// ExpiresAtGT applies the GT predicate on the "expires_at" field.
func ExpiresAtGT(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldExpiresAt), v))
	})
}

// ExpiresAtGTE applies the GTE predicate on the "expires_at" field.
func ExpiresAtGTE(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldExpiresAt), v))
	})
}

// ExpiresAtLT applies the LT predicate on the "expires_at" field.
func ExpiresAtLT(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldExpiresAt), v))
	})
}

// ExpiresAtLTE applies the LTE predicate on the "expires_at" field.
func ExpiresAtLTE(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldExpiresAt), v))
	})
}

// ExpiresAtIsNil applies the IsNil predicate on the "expires_at" field.
func ExpiresAtIsNil() predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldExpiresAt)))
	})
}

// ExpiresAtNotNil applies the NotNil predicate on the "expires_at" field.
func ExpiresAtNotNil() predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldExpiresAt)))
	})
}

// CancelledEQ applies the EQ predicate on the "cancelled" field.
func CancelledEQ(v bool) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCancelled), v))
	})
}

// CancelledNEQ applies the NEQ predicate on the "cancelled" field.
func CancelledNEQ(v bool) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCancelled), v))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Subscription {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Subscription {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Subscription(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(UserTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(UserInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Subscription) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Subscription) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Subscription) predicate.Subscription {
	return predicate.Subscription(func(s *sql.Selector) {
		p(s.Not())
	})
}
