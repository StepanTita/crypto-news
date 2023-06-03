package config

type Listener interface {
	TelegramApiToken() string
	Sources() []string
}

type listener struct {
	telegramApiToken string
	sources          []string
}

func NewListener(tgApiToken string, sources []string) Listener {
	return &listener{
		telegramApiToken: tgApiToken,
		sources:          sources,
	}
}

func (l listener) TelegramApiToken() string {
	return l.telegramApiToken
}

func (l listener) Sources() []string {
	return l.sources
}
