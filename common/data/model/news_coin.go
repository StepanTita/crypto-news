package model

import "github.com/google/uuid"

type NewsCoin struct {
	ID     uuid.UUID `db:"id,omitempty"`
	Code   string    `db:"code"`
	NewsID uuid.UUID `db:"news_id,omitempty"`
}

func (n NewsCoin) TableName() string {
	return NEWS_COINS
}
