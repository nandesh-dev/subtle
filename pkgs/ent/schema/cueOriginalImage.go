package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CueOriginalImage holds the schema definition for the CueOriginalImage entity.
type CueOriginalImage struct {
	ent.Schema
}

// Fields of the CueOriginalImage.
func (CueOriginalImage) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("position"),
		field.Bytes("data"),
	}
}

// Edges of the CueOriginalImage.
func (CueOriginalImage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cue", Cue.Type).
			Ref("cue_original_images"),
	}
}
