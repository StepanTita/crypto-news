package config

type Poster interface {
	BearerToken() string
	ApiKey() string
	ApiSecret() string
}

type YamlPosterConfig struct {
	ApiKey      string `yaml:"api_key"`
	ApiSecret   string `yaml:"api_secret"`
	BearerToken string `yaml:"bearer_token"`
}

type poster struct {
	apiKey      string
	apiSecret   string
	bearerToken string
}

func NewPoster(posterConfig YamlPosterConfig) Poster {
	return &poster{
		apiKey:      posterConfig.ApiKey,
		apiSecret:   posterConfig.ApiSecret,
		bearerToken: posterConfig.BearerToken,
	}
}

func (p poster) BearerToken() string {
	return p.bearerToken
}

func (p poster) ApiKey() string {
	return p.apiKey
}

func (p poster) ApiSecret() string {
	return p.apiSecret
}
