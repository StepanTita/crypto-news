package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	StatusPending   = "pending"
	StatusProcessed = "processed"
	StatusFailed    = "failed"
)

const (
	ResourceTypeSource = "source"
)

type News struct {
	ID          uuid.UUID  `db:"id,omitempty"`
	CreatedAt   time.Time  `db:"created_at,omitempty"`
	UpdatedAt   *time.Time `db:"updated_at"`
	PublishedAt *time.Time `db:"published_at"`
	Url         *string    `db:"url"`

	// Data
	Media *NewsMedia `db:"media"`

	Source         *string `db:"source"`
	OriginalSource *string `db:"original_source"`

	Status *string `db:"status"`

	Coins []Coin `db:"-"`
}

func (n News) TableName() string {
	return NEWS
}

// NewsMedia struct for media content, can be a nearly free form json
type NewsMedia struct {
	Title     *string             `json:"title"`
	Text      *string             `json:"text"`
	Resources []NewsMediaResource `json:"resources"`
}

type NewsMediaResource struct {
	Type *string         `json:"type"`
	URL  *string         `json:"url"`
	Meta json.RawMessage `json:"meta"`
}

func (a NewsMedia) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *NewsMedia) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type UpdateNewsParams struct {
	Status *string `db:"status"`

	UpdatedAt *time.Time `db:"updated_at"`
}

func (n UpdateNewsParams) TableName() string {
	return NEWS
}

type MetaLinksData struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}
