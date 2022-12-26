package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Record holds the schema definition for the Record entity.
type Record struct {
	ent.Schema
}

// Fields of the Record.
func (Record) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.Int("uid"),
		field.String("type"), // oob type
		field.String("remote_addr"),
		// dns
		field.String("dns_name").
			Default(""),
		// http
		field.String("http_url").
			Default(""),
		field.String("http_method").
			Default(""),
		field.String("http_raw").
			Default(""),
		// jndi
		field.String("jndi_protocol").
			Default(""),
		field.String("jndi_path").
			Default(""),

		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Record.
func (Record) Edges() []ent.Edge {
	return nil
}
