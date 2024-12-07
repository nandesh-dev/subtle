package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Segment holds the schema definition for the Segment entity.
type Segment struct {
	ent.Schema
}

// Fields of the Segment.
func (Segment) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("start_time").
			GoType(time.Second),
		field.Int64("end_time").
			GoType(time.Second),
		field.String("text").
			Optional(),
		field.String("original_text").
			Optional(),
		field.Bytes("original_image").
			Optional(),
	}
}

// Edges of the Segment.
func (Segment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("subtitle", Subtitle.Type).
			Ref("segments"),
	}
}
