package model

import (
	"time"

	"github.com/google/uuid"
)

// User TODO: use for preferences
type User struct {
	ID        uuid.UUID  `db:"id,omitempty"`
	CreatedAt time.Time  `db:"created_at,omitempty"`
	UpdatedAt *time.Time `db:"updated_at"`

	Username  *string `db:"username"`
	FirstName *string `db:"first_name"`
	LastName  *string `db:"last_name"`

	Platform *string `db:"platform"`
}

func (u User) TableName() string {
	return USERS
}
