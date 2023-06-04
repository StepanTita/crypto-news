package services

import chat_bot "github.com/StepanTita/go-EdgeGPT/chat-bot"

type generationsResponse struct {
	content string
	links   []chat_bot.ResponseLink
	coins   []string
}
