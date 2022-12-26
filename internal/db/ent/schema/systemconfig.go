package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// SystemConfig holds the schema definition for the SystemConfig entity.
type SystemConfig struct {
	ent.Schema
}

// Fields of the SystemConfig.
func (SystemConfig) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").
			Unique(),
		field.String("value").
			Default(""),
	}
}

// Edges of the SystemConfig.
func (SystemConfig) Edges() []ent.Edge {
	return nil
}
