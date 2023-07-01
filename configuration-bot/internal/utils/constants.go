package utils

type Command string

const (
	StartCommand     Command = "start"
	WhitelistCommand Command = "whitelist"
	UnlistCommand    Command = "unlist"
	TokenCommand     Command = "token"
)

func (c Command) Command() string {
	return "/" + string(c)
}

func (c Command) String() string {
	return string(c)
}
