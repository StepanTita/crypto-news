package bot

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"

	"gpt/internal/config"
)

type openAIBot struct {
	log *logrus.Entry

	cfg config.Config

	client *openai.Client
}

func NewOpenAI(cfg config.Config) Bot {
	return &openAIBot{
		log: cfg.Logging().WithField("[BOT]", "openai-gpt4"),
		cfg: cfg,

		client: openai.NewClientWithConfig(openai.DefaultConfig(cfg.AuthToken())),
	}
}

func (b *openAIBot) Ask(ctx context.Context, prompt, context string, language string) (*Message, error) {
	resp, err := b.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo16K,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: context,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: fmt.Sprintf("Your entire reply should be translated to the following language: %s", language),
				},
			},
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create chat completion request")
	}

	return &Message{Text: resp.Choices[0].Message.Content}, nil
}
