package model

import "time"

type Channel struct {
	ChannelID int64     `db:"channel_id,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	Priority  int32     `db:"priority,omitempty"`
}

func (n Channel) TableName() string {
	return CHANNELS
}
