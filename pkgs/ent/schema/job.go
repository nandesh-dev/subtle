package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		field.Bool("running").
			Default(false),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return nil
}
