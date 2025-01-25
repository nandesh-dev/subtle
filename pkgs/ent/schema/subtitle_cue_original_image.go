package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SubtitleCueOriginalImageSchema holds the schema definition for the SubtitleCueOriginalImageSchema entity.
type SubtitleCueOriginalImageSchema struct {
	ent.Schema
}

// Annotations of the SubtitleCueOriginalImageSchema.
func (SubtitleCueOriginalImageSchema) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "subtitle_cue_original_images"},
    }
}

// Fields of the SubtitleCueOriginalImageSchema.
func (SubtitleCueOriginalImageSchema) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("position"),
		field.Bytes("data"),
	}
}

// Edges of the SubtitleCueOriginalImageSchema.
func (SubtitleCueOriginalImageSchema) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cue", SubtitleCueSchema.Type).
			Ref("original_images"),
	}
}
