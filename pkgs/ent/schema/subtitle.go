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
		field.Bool("processing").
			Default(false),
		field.Bool("extracted").
			Default(false),
		field.Bool("formated").
			Default(false),
		field.Bool("exported").
			Default(false),
		field.Bool("import_is_external").
			Optional(),
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
		edge.To("segments", Segment.Type),
		edge.From("video", Video.Type).
			Ref("subtitles"),
	}
}
