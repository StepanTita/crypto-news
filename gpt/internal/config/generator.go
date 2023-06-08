package config

import "time"

type Generator interface {
	GenerateEvery() time.Duration
	ShortSummaryPrompt() string
}

type generator struct {
	generateEvery      time.Duration
	shortSummaryPrompt string
}

func NewGenerator(generateEvery time.Duration, shortSummaryPrompt string) Generator {
	return &generator{
		generateEvery:      generateEvery,
		shortSummaryPrompt: shortSummaryPrompt,
	}
}

func (g generator) GenerateEvery() time.Duration {
	return g.generateEvery
}

func (g generator) ShortSummaryPrompt() string {
	return g.shortSummaryPrompt
}
