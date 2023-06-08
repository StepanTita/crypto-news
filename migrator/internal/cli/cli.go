package cli

import (
	"os"

	"github.com/urfave/cli/v2"

	"common/config"
	"migrator/migrate"
)

func Run(args []string) bool {
	cfg := config.NewFromFile(os.Getenv("CONFIG"))
	log := cfg.Logging()

	log.WithField("version", cfg.Version()).Info("Running version")

	defer func() {
		if rvr := recover(); rvr != nil {
			log.Error("internal panicked: ", rvr)
		}
	}()

	app := &cli.App{
		Commands: cli.Commands{
			{
				Name:  "migrate",
				Usage: "migrate service database",
				Subcommands: cli.Commands{
					{
						Name:  "up",
						Usage: "migrate database up",
						Action: func(c *cli.Context) error {
							log.Debug("Migrating up")
							if err := migrate.Migrate(cfg, migrate.Up); err != nil {
								return err
							}
							return nil
						},
					},
					{
						Name:  "down",
						Usage: "migrate database down",
						Action: func(c *cli.Context) error {
							log.Debug("Migrating down")
							if err := migrate.Migrate(cfg, migrate.Down); err != nil {
								return err
							}
							return nil
						},
					},
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
