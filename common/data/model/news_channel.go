package model

import (
	"github.com/google/uuid"
)

type NewsChannel struct {
	ID        uuid.UUID `db:"id,omitempty"`
	ChannelID int64     `db:"channel_id,omitempty"`
	Status    *string   `db:"status"`
	NewsID    uuid.UUID `db:"news_id,omitempty"`
}

func (n NewsChannel) TableName() string {
	return NEWS_CHANNELS
}

type UpdateNewsChannelParams struct {
	Status *string `db:"status"`
}

func (n UpdateNewsChannelParams) TableName() string {
	return NEWS_CHANNELS
}
