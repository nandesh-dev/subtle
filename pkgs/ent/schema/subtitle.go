package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/nandesh-dev/subtle/pkgs/language"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
)

// SubtitleSchema holds the schema definition for the SubtitleSchema entity.
type SubtitleSchema struct {
	ent.Schema
}

// Annotations of the SubtitleSchema.
func (SubtitleSchema) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "subtitles"},
	}
}

// Fields of the SubtitleSchema.
func (SubtitleSchema) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("language").
			GoType(language.English),
		field.Enum("stage").
			Values("detected", "extracted", "formated", "exported"),
		field.Bool("is_processing").
			Default(false),
		field.Bool("import_is_external").
			Default(false),
		field.String("import_format").
			GoType(subtitle.SRT),
		field.Int("import_video_stream_index").
			Optional().
			Nillable(),
		field.String("export_path").
			Optional().
			Nillable(),
		field.String("export_format").
			GoType(subtitle.SRT).
			Optional().
			Nillable(),
	}
}

// Edges of the SubtitleSchema.
func (SubtitleSchema) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("cues", SubtitleCueSchema.Type),
		edge.From("video", VideoSchema.Type).
			Ref("subtitles"),
	}
}
