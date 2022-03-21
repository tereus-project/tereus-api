// Code generated by entc, DO NOT EDIT.

package submission

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/tereus-project/tereus-api/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
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
func IDNotIn(ids ...uuid.UUID) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
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
func IDGT(id uuid.UUID) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// SourceLanguage applies equality check predicate on the "source_language" field. It's identical to SourceLanguageEQ.
func SourceLanguage(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSourceLanguage), v))
	})
}

// TargetLanguage applies equality check predicate on the "target_language" field. It's identical to TargetLanguageEQ.
func TargetLanguage(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTargetLanguage), v))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// SourceLanguageEQ applies the EQ predicate on the "source_language" field.
func SourceLanguageEQ(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSourceLanguage), v))
	})
}

// SourceLanguageNEQ applies the NEQ predicate on the "source_language" field.
func SourceLanguageNEQ(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSourceLanguage), v))
	})
}

// SourceLanguageIn applies the In predicate on the "source_language" field.
func SourceLanguageIn(vs ...string) predicate.Submission {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Submission(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSourceLanguage), v...))
	})
}

// SourceLanguageNotIn applies the NotIn predicate on the "source_language" field.
func SourceLanguageNotIn(vs ...string) predicate.Submission {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Submission(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSourceLanguage), v...))
	})
}

// SourceLanguageGT applies the GT predicate on the "source_language" field.
func SourceLanguageGT(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSourceLanguage), v))
	})
}

// SourceLanguageGTE applies the GTE predicate on the "source_language" field.
func SourceLanguageGTE(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSourceLanguage), v))
	})
}

// SourceLanguageLT applies the LT predicate on the "source_language" field.
func SourceLanguageLT(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSourceLanguage), v))
	})
}

// SourceLanguageLTE applies the LTE predicate on the "source_language" field.
func SourceLanguageLTE(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSourceLanguage), v))
	})
}

// SourceLanguageContains applies the Contains predicate on the "source_language" field.
func SourceLanguageContains(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldSourceLanguage), v))
	})
}

// SourceLanguageHasPrefix applies the HasPrefix predicate on the "source_language" field.
func SourceLanguageHasPrefix(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldSourceLanguage), v))
	})
}

// SourceLanguageHasSuffix applies the HasSuffix predicate on the "source_language" field.
func SourceLanguageHasSuffix(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldSourceLanguage), v))
	})
}

// SourceLanguageEqualFold applies the EqualFold predicate on the "source_language" field.
func SourceLanguageEqualFold(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldSourceLanguage), v))
	})
}

// SourceLanguageContainsFold applies the ContainsFold predicate on the "source_language" field.
func SourceLanguageContainsFold(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldSourceLanguage), v))
	})
}

// TargetLanguageEQ applies the EQ predicate on the "target_language" field.
func TargetLanguageEQ(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTargetLanguage), v))
	})
}

// TargetLanguageNEQ applies the NEQ predicate on the "target_language" field.
func TargetLanguageNEQ(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTargetLanguage), v))
	})
}

// TargetLanguageIn applies the In predicate on the "target_language" field.
func TargetLanguageIn(vs ...string) predicate.Submission {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Submission(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTargetLanguage), v...))
	})
}

// TargetLanguageNotIn applies the NotIn predicate on the "target_language" field.
func TargetLanguageNotIn(vs ...string) predicate.Submission {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Submission(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTargetLanguage), v...))
	})
}

// TargetLanguageGT applies the GT predicate on the "target_language" field.
func TargetLanguageGT(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTargetLanguage), v))
	})
}

// TargetLanguageGTE applies the GTE predicate on the "target_language" field.
func TargetLanguageGTE(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTargetLanguage), v))
	})
}

// TargetLanguageLT applies the LT predicate on the "target_language" field.
func TargetLanguageLT(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTargetLanguage), v))
	})
}

// TargetLanguageLTE applies the LTE predicate on the "target_language" field.
func TargetLanguageLTE(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTargetLanguage), v))
	})
}

// TargetLanguageContains applies the Contains predicate on the "target_language" field.
func TargetLanguageContains(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldTargetLanguage), v))
	})
}

// TargetLanguageHasPrefix applies the HasPrefix predicate on the "target_language" field.
func TargetLanguageHasPrefix(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldTargetLanguage), v))
	})
}

// TargetLanguageHasSuffix applies the HasSuffix predicate on the "target_language" field.
func TargetLanguageHasSuffix(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldTargetLanguage), v))
	})
}

// TargetLanguageEqualFold applies the EqualFold predicate on the "target_language" field.
func TargetLanguageEqualFold(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldTargetLanguage), v))
	})
}

// TargetLanguageContainsFold applies the ContainsFold predicate on the "target_language" field.
func TargetLanguageContainsFold(v string) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldTargetLanguage), v))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Submission {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Submission(func(s *sql.Selector) {
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
func CreatedAtNotIn(vs ...time.Time) predicate.Submission {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Submission(func(s *sql.Selector) {
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
func CreatedAtGT(v time.Time) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Submission) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Submission) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
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
func Not(p predicate.Submission) predicate.Submission {
	return predicate.Submission(func(s *sql.Selector) {
		p(s.Not())
	})
}
