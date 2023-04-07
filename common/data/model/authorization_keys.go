package model

import "time"

type AuthorizationKeys struct {
	AuthorizationToken     string    `json:"authorization_token"`
	RefreshToken           string    `json:"refresh_token"`
	AuthorizationExpiresAt time.Time `json:"authorization_expires_at"`
	RefreshExpiresAt       time.Time `json:"refresh_expires_at"`
}

func (n AuthorizationKeys) TableName() string {
	panic("Shouldn't be called, NO TABLES in no-sql. Method is implemented only for the interface compliance!")
}
