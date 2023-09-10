package cli

import (
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"

	"gpt/internal/config"
	"gpt/internal/services"
)

func Run(args []string) bool {
	cfg := config.NewFromFile(os.Getenv("CONFIG"))
	log := cfg.Logging()

	log.WithField("version", cfg.Version()).Info("Running version")

	defer func() {
		if rvr := recover(); rvr != nil {
			log.Error("internal panicked: ", rvr, string(debug.Stack()))
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
