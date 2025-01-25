package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SubtitleCueSchema holds the schema definition for the SubtitleCueSchema entity.
type SubtitleCueSchema struct {
	ent.Schema
}

// Annotations of the SubtitleCueSchema.
func (SubtitleCueSchema) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "subtitle_cues"},
    }
}

// Fields of the SubtitleCueSchema.
func (SubtitleCueSchema) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("timestamp_start").
			GoType(time.Second),
		field.Int64("timestamp_end").
			GoType(time.Second),
	}
}

// Edges of the SubtitleCueSchema.
func (SubtitleCueSchema) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("subtitle", SubtitleSchema.Type).
			Ref("cues"),
		edge.To("content_segments", SubtitleCueContentSegmentSchema.Type),
		edge.To("original_images", SubtitleCueOriginalImageSchema.Type),
	}
}
