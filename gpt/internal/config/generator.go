package config

import "time"

type Generator interface {
	GenerateEvery() time.Duration
	ImagesPrompt() string
	QueryContext() string
}

type generator struct {
	generateEvery time.Duration
	imagesPrompt  string
	queryContext  string
}

func NewGenerator(generateEvery time.Duration, imagesPrompt, queryContext string) Generator {
	return &generator{
		generateEvery: generateEvery,
		imagesPrompt:  imagesPrompt,
		queryContext:  queryContext,
	}
}

func (g generator) GenerateEvery() time.Duration {
	return g.generateEvery
}

func (g generator) ImagesPrompt() string {
	return g.imagesPrompt
}

func (g generator) QueryContext() string {
	return g.queryContext
}
