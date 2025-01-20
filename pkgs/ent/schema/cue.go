package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Segment holds the schema definition for the Segment entity.
type Cue struct {
	ent.Schema
}

// Fields of the Segment.
func (Cue) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("timestamp_start").
			GoType(time.Second),
		field.Int64("timestamp_end").
			GoType(time.Second),
	}
}

// Edges of the Segment.
func (Cue) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("subtitle", Subtitle.Type).
			Ref("cues"),
		edge.To("cue_content_segments", CueContentSegment.Type),
		edge.To("cue_original_images", CueOriginalImage.Type),
	}
}
