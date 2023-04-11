package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTwitter_Port(t *testing.T) {
	tw := NewTwitter(YamlTwitterConfig{
		Authenticator: struct {
			Address  string `yaml:"address"`
			TokenURL string `yaml:"token_url"`
			AuthURL  string `yaml:"auth_url"`
		}{
			Address: "some.url.com:8080",
		},
	},
	)
	require.Equal(t, ":8080", tw.Port())
}
