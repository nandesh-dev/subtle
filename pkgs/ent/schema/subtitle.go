package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Subtitle holds the schema definition for the Subtitle entity.
type Subtitle struct {
	ent.Schema
}

// Fields of the Subtitle.
func (Subtitle) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("language"),
		field.Enum("stage").
			Values("detected", "extracted", "formated", "exported"),
		field.Bool("is_processing").
			Default(false),
		field.Bool("import_is_external").
			Default(false),
		field.String("import_format").
			Optional(),
		field.Int32("import_video_stream_index").
			Optional(),
		field.String("export_path").
			Optional(),
		field.String("export_format").
			Optional(),
	}
}

// Edges of the Subtitle.
func (Subtitle) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("cues", Cue.Type),
		edge.From("video", Video.Type).
			Ref("subtitles"),
	}
}
