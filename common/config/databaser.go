package config

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// PostgresDriver list of supported drivers
const (
	PostgresDriver = "postgres"
)

type Databaser interface {
	DB() *sqlx.DB
	Driver() string
}

type databaser struct {
	db     *sqlx.DB
	driver string
}

type YamlDatabaseConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Port     string `yaml:"port"`
	SslMode  string `yaml:"ssl_mode"`
	Driver   string `yaml:"driver"`
}

func (c YamlDatabaseConfig) toPSQLPath() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host, c.User, c.Password, c.Name, c.Port, c.SslMode)
}

func NewDatabaser(dbConfig YamlDatabaseConfig) Databaser {
	var db *sqlx.DB
	var err error

	switch dbConfig.Driver {
	case PostgresDriver:
		db, err = sqlx.Connect(dbConfig.Driver, dbConfig.toPSQLPath())
		if err != nil {
			panic(errors.Wrapf(err, "failed to open database connection: %s", dbConfig.toPSQLPath()))
		}
	default:
		panic(errors.Errorf("provided driver unsupported: %s", dbConfig.Driver))
	}
	return &databaser{
		db:     db,
		driver: dbConfig.Driver,
	}
}

func (d *databaser) DB() *sqlx.DB {
	return d.db
}

func (d *databaser) Driver() string {
	return d.driver
}
