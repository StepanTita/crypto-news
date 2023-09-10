package config

type Listener interface {
	TelegramApiToken() string
}

type listener struct {
	telegramApiToken string
}

func NewListener(tgApiToken string) Listener {
	return &listener{
		telegramApiToken: tgApiToken,
	}
}

func (l listener) TelegramApiToken() string {
	return l.telegramApiToken
}
