package migrate

import (
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"

	commoncfg "common/config"
	"migrator/assets"
)

const (
	Up   = "up"
	Down = "down"
)

// migrateUp Migrates database up
func migrateUp(cfg commoncfg.Config) (int, error) {
	applied, err := migrate.Exec(cfg.DB().DB, cfg.Driver(), assets.Migrations, migrate.Up)

	if err != nil {
		return 0, errors.Wrap(err, "failed to apply migrations")
	}

	return applied, nil
}

// migrateDown Migrates database down
func migrateDown(cfg commoncfg.Config) (int, error) {
	applied, err := migrate.Exec(cfg.DB().DB, cfg.Driver(), assets.Migrations, migrate.Down)
	if err != nil {
		return 0, errors.Wrap(err, "failed to apply migrations")
	}
	return applied, nil
}

func Migrate(cfg commoncfg.Config, direction string) error {
	migrator := migrateUp
	switch direction {
	case Down:
		migrator = migrateDown
	}
	migrationsCount, err := migrator(cfg)
	if err != nil {
		return err
	}
	cfg.Logging().WithField("applied", migrationsCount).Info("Migrations applied")
	return nil
}
