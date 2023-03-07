package main

import (
	"os"

	"common"
	"crypto-news/internal/cli"
)

func main() {
	common.SetupWorkingDirectory()
	if !cli.Run(os.Args) {
		os.Exit(2)
	}
	os.Exit(0)
}
