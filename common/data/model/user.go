package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleAdmin  = "admin"
	RoleReader = "reader"
)

type User struct {
	ID        uuid.UUID  `db:"id,omitempty"`
	CreatedAt time.Time  `db:"created_at,omitempty"`
	UpdatedAt *time.Time `db:"updated_at"`
	FirstName *string    `db:"first_name"`
	LastName  *string    `db:"last_name"`
	Username  *string    `db:"username"`
	Platform  *string    `db:"platform"`
	Role      *string    `db:"role"`
}

func (u User) TableName() string {
	return USERS
}
