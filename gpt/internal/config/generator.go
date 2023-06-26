package config

import "time"

type Generator interface {
	GenerateEvery() time.Duration
	ShortSummaryPrompt() string
	ImagesPrompt() string
}

type generator struct {
	generateEvery      time.Duration
	shortSummaryPrompt string
	imagesPrompt       string
}

func NewGenerator(generateEvery time.Duration, shortSummaryPrompt, imagesPrompt string) Generator {
	return &generator{
		generateEvery:      generateEvery,
		shortSummaryPrompt: shortSummaryPrompt,
		imagesPrompt:       imagesPrompt,
	}
}

func (g generator) GenerateEvery() time.Duration {
	return g.generateEvery
}

func (g generator) ShortSummaryPrompt() string {
	return g.shortSummaryPrompt
}

func (g generator) ImagesPrompt() string {
	return g.imagesPrompt
}
