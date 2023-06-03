package main

import (
	"os"

	"gpt/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(2)
	}
	os.Exit(0)
}
