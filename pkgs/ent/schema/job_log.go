package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// JobLogSchema holds the schema definition for the JobSchema entity.
type JobLogSchema struct {
	ent.Schema
}

// Annotations of the JobLogSchema.
func (JobLogSchema) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "job_logs"},
	}
}

// Fields of the JobLogSchema.
func (JobLogSchema) Fields() []ent.Field {
	return []ent.Field{
		field.Time("start_timestamp"),
		field.Int("duration"),
	}
}

// Edges of the JobLogSchema.
func (JobLogSchema) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("job", JobSchema.Type).
			Ref("logs"),
	}
}
