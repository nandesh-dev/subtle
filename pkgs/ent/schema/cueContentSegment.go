package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CueContentSegment holds the schema definition for the CueContentSegment entity.
type CueContentSegment struct {
	ent.Schema
}

// Fields of the CueContentSegment.
func (CueContentSegment) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("position"),
		field.String("text"),
	}
}

// Edges of the CueContentSegment.
func (CueContentSegment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cue", Cue.Type).
			Ref("cue_content_segments"),
	}
}
