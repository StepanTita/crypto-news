package model

import (
	"time"

	"github.com/google/uuid"
)

type Title struct {
	ID          uuid.UUID  `db:"id,omitempty"`
	CreatedAt   time.Time  `db:"created_at,omitempty"`
	UpdatedAt   *time.Time `db:"updated_at"`
	Title       *string    `db:"title"`
	Summary     *string    `db:"summary"`
	Hash        *string    `db:"hash"`
	URL         *string    `db:"url"`
	Status      *string    `db:"status"`
	ReleaseDate *time.Time `db:"release_date"`
}

func (t Title) TableName() string {
	return TITLES
}

type UpdateTitleParams struct {
	UpdatedAt *time.Time `db:"updated_at"`
	Status    *string    `db:"status"`
}

func (t UpdateTitleParams) TableName() string {
	return TITLES
}
