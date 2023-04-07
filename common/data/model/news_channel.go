package model

import (
	"github.com/google/uuid"
)

type NewsChannel struct {
	ID        uuid.UUID `db:"id,omitempty"`
	ChannelID int64     `db:"channel_id,omitempty"`
	NewsID    uuid.UUID `db:"news_id,omitempty"`
}

func (n NewsChannel) TableName() string {
	return NEWS_CHANNELS
}
