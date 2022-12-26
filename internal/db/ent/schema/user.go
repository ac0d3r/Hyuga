package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("sid").
			Unique(), // 子域名 & 短ID
		field.String("username").
			Unique(),
		field.String("password").
			Sensitive(),
		field.String("token").
			Default(""),

		field.String("from_user").
			Default(""),
		field.UUID("invite_code", uuid.UUID{}).
			Default(uuid.New).Unique(),

		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
