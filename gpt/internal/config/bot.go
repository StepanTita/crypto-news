package config

type BotConfig interface {
	AuthToken() string
}

type botConfig struct {
	authToken string
}

func NewBotConfig(authToken string) BotConfig {
	return &botConfig{
		authToken: authToken,
	}
}

func (b botConfig) AuthToken() string {
	return b.authToken
}
