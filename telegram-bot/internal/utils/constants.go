package utils

type Command string

const (
	StartCommand     Command = "start"
	SubscribeCommand Command = "subscribe"
)

func (c Command) Command() string {
	return "/" + string(c)
}

func (c Command) String() string {
	return string(c)
}
