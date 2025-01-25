package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SubtitleCueContentSegmentSchema holds the schema definition for the SubtitleCueContentSegmentSchema entity.
type SubtitleCueContentSegmentSchema struct {
	ent.Schema
}

// Annotations of the SubtitleCueContentSchema.
func (SubtitleCueContentSegmentSchema) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "subtitle_cue_content_segments"},
    }
}

// Fields of the SubtitleCueContentSegmentSchema.
func (SubtitleCueContentSegmentSchema) Fields() []ent.Field {
	return []ent.Field{
		field.Int("position"),
		field.String("text"),
	}
}

// Edges of the SubtitleCueContentSegmentSchema.
func (SubtitleCueContentSegmentSchema) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cue", SubtitleCueSchema.Type).
			Ref("content_segments"),
	}
}
