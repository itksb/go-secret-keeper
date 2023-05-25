package migrate

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
	"io/fs"
)

// Migrations - virtual file system
//
//go:embed migrations/*.sql
var Migrations embed.FS

// IAppMigrator - migrator interface
type IAppMigrator interface {
	Migrate(dsn string, path fs.FS) error
}

// AppMigratorFunc - migrator adapter
type AppMigratorFunc func(dsn string, path fs.FS) error

// Migrate - executes the migration
func (f AppMigratorFunc) Migrate(dsn string, path fs.FS) error {
	return f(dsn, path)
}

// Migrate - executes the migration
func Migrate(dsn string, path fs.FS) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	goose.SetBaseFS(path)
	return goose.Up(db, "migrations")
}
