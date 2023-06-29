package main

import (
	"os"

	"configuration-bot/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(2)
	}
	os.Exit(0)
}
