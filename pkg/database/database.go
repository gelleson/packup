package database

import (
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	DSN string
}

type Database struct {
	config    Config
	connected bool
	db        *gorm.DB
}

func NewDatabase(config Config) *Database {
	return &Database{config: config}
}

func (d *Database) Connect() error {

	d.connected = false

	db, err := gorm.Open(sqlite.Open(d.config.DSN))

	if err != nil {
		return errors.Wrap(err, "database")
	}

	d.db = db

	d.connected = true

	return nil
}

func (d *Database) Conn() *gorm.DB {

	return d.db
}
