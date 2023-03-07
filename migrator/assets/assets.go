package assets

import (
	"embed"

	"github.com/rubenv/sql-migrate"
)

//go:embed migrations/postgres/*.sql
var migrationsSQL embed.FS

var Migrations = migrate.EmbedFileSystemMigrationSource{
	FileSystem: migrationsSQL,
	Root:       "migrations/postgres",
}
