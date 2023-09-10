package model

import (
	"time"

	"github.com/google/uuid"
)

type RawNews struct {
	ID        uuid.UUID `db:"id,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	TitleID   uuid.UUID `db:"title_id"`
	Body      *string   `db:"body"`
}

func (t RawNews) TableName() string {
	return RAW_NEWS
}
