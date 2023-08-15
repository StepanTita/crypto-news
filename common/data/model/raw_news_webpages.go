package model

import (
	"time"

	"github.com/google/uuid"
)

type RawNewsWebpage struct {
	ID        uuid.UUID `db:"id,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	Body      *string   `db:"body"`
}

func (t RawNewsWebpage) TableName() string {
	return RAW_NEWS_WEBPAGES
}
