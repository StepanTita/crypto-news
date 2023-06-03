package cli

import (
	"os"

	"github.com/urfave/cli/v2"

	"gpt/internal/config"
	"gpt/internal/services"
)

func Run(args []string) bool {
	cfg := config.NewFromFile(os.Getenv("CONFIG"))
	log := cfg.Logging()

	defer func() {
		if rvr := recover(); rvr != nil {
			log.Error("internal panicked: ", rvr)
		}
	}()

	svc := services.New(cfg)

	app := &cli.App{
		Commands: cli.Commands{
			{
				Name:  "run",
				Usage: "run gpt daemon",
				Action: func(c *cli.Context) error {
					return svc.Run(c.Context)
				},
			},
		},
	}

	if err := app.Run(args); err != nil {
		log.Fatal(err, ": service initialization failed")
		return false
	}

	return true
}