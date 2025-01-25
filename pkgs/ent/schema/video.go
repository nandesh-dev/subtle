package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// VideoSchema holds the schema definition for the VideoSchema entity.
type VideoSchema struct {
	ent.Schema
}

// Annotations of the VideoSchema.
func (VideoSchema) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "videos"},
    }
}

// Fields of the VideoSchema.
func (VideoSchema) Fields() []ent.Field {
	return []ent.Field{
		field.String("filepath"),
	}
}

// Edges of the VideoSchema.
func (VideoSchema) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subtitles", SubtitleSchema.Type),
	}
}
