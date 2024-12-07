package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Routine holds the schema definition for the Routine entity.
type Routine struct {
	ent.Schema
}

// Fields of the Routine.
func (Routine) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		field.Bool("running"),
	}
}

// Edges of the Routine.
func (Routine) Edges() []ent.Edge {
	return nil
}
