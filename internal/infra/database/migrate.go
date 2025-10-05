package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/tern/v2/migrate"
)

//go:embed migrations
var migrationFS embed.FS

func virtualFs() (fs.FS, error) {
	vs, err := fs.Sub(migrationFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("virtualfs: %w", err)
	}
	return vs, nil
}

func (db *DBClient) MigrateUp() error {
	virtualFileSys, err := virtualFs()
	if err != nil {
		return fmt.Errorf("create virtual fs: %w", err)
	}

	conn, err := db.Pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("acquire pgx connection: %w", err)
	}
	defer conn.Release()

	m, err := migrate.NewMigrator(context.Background(), conn.Conn(), "schema_version")
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	if err := m.LoadMigrations(virtualFileSys); err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}

	if err := m.Migrate(context.Background()); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
