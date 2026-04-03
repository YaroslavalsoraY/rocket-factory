package pg

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type PgMigrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *PgMigrator {
	return &PgMigrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *PgMigrator) Up() error {
	err := goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return err
	}

	return nil
}
