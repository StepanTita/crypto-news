package config

import (
	"fmt"

	"golang.org/x/oauth2"
)

type Twitter interface {
	OAuthConfig() *oauth2.Config

	AuthAddress() string
}

type YamlTwitterConfig struct {
	ApiKey        string `yaml:"api_key"`
	ApiSecret     string `yaml:"api_secret"`
	BearerToken   string `yaml:"bearer_token"`
	Authenticator struct {
		Address  string `yaml:"address"`
		TokenURL string `yaml:"token_url"`
		AuthURL  string `yaml:"auth_url"`
	} `json:"authenticator"`
}

type twitter struct {
	apiKey      string
	apiSecret   string
	bearerToken string
	authAddress string
	tokenURL    string
	authURL     string
}

func NewTwitter(posterConfig YamlTwitterConfig) Twitter {
	return &twitter{
		apiKey:      posterConfig.ApiKey,
		apiSecret:   posterConfig.ApiSecret,
		bearerToken: posterConfig.BearerToken,
		authAddress: posterConfig.Authenticator.Address,
		authURL:     posterConfig.Authenticator.AuthURL,
		tokenURL:    posterConfig.Authenticator.TokenURL,
	}
}

func (t twitter) OAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://127.0.0.1%s/oauth/callback", t.authAddress),
		ClientID:     t.apiKey,
		ClientSecret: t.apiSecret,
		Scopes:       []string{"tweet.read", "users.read", "tweet.write", "offline.access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   t.authURL,
			TokenURL:  t.tokenURL,
			AuthStyle: 0,
		},
	}
}

func (t twitter) AuthAddress() string {
	return t.authAddress
}
