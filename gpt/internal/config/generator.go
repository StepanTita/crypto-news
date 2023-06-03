package config

import "time"

type Generator interface {
	GenerateEvery() time.Duration
}

type generator struct {
	generateEvery time.Duration
}

func NewGenerator(generateEvery time.Duration) Generator {
	return &generator{
		generateEvery: generateEvery,
	}
}

func (g generator) GenerateEvery() time.Duration {
	return g.generateEvery
}
