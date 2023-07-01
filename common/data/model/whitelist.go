package model

import "github.com/google/uuid"

type Whitelist struct {
	ID       uuid.UUID  `db:"id,omitempty"`
	Username *string    `db:"username"`
	Token    *uuid.UUID `db:"token"`
}

func (l Whitelist) TableName() string {
	return WHITELIST
}
