package services

import chat_bot "github.com/StepanTita/go-EdgeGPT/chat-bot"

type generationsResponse struct {
	content string
	coins   []string

	sources   []chat_bot.ResponseLink
	resources []chat_bot.ResourceLink
}
