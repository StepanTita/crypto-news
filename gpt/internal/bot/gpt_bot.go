package bot

import "context"

type Message struct {
	Text string
}

type Bot interface {
	Ask(ctx context.Context, prompt, context string, language string) (*Message, error)
}
