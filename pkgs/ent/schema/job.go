package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// JobSchema holds the schema definition for the JobSchema entity.
type JobSchema struct {
	ent.Schema
}

// Annotations of the JobSchema.
func (JobSchema) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "jobs"},
	}
}

// Fields of the JobSchema.
func (JobSchema) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Bool("is_running"),
	}
}

// Edges of the JobSchema.
func (JobSchema) Edges() []ent.Edge {
	return []ent.Edge{}
}
