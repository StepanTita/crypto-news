package main

import (
	"os"

	"common"
	"twitter-bot/internal/cli"
)

func main() {
	common.SetupWorkingDirectory()
	if !cli.Run(os.Args) {
		os.Exit(2)
	}
	os.Exit(0)
}
