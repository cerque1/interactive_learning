package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type Migrator struct {
	srcDriver source.Driver
}

func NewMigrator(sqlFiles embed.FS, dirName string) *Migrator {
	driver, err := iofs.New(sqlFiles, dirName)
	if err != nil {
		log.Fatal(err)
	}

	return &Migrator{driver}
}

func (m *Migrator) ApplyMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("error create db instance " + err.Error())
	}

	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, "psql_db", driver)
	if err != nil {
		log.Fatal("error create migrator " + err.Error())
	}

	defer migrator.Close()

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("error apply migration: %v", err)
	}
}
