package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Token хранит PAT и привязку к user_id.
type Token struct {
	ent.Schema
}

func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.UUID("user_id", uuid.UUID{}).
			Comment("UUID пользователя из UserService"),
		field.String("token").
			Unique().
			Sensitive(),
		field.Time("issued_at").
			Default(time.Now),
		field.Time("expires_at").
			Optional().
			Nillable(),
		field.Bool("revoked").
			Default(false),
	}
}

func (Token) Edges() []ent.Edge {
	// У AuthService нет реальных связей на уровне БД к users.
	return nil
}
