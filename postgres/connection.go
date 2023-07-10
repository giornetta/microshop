package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, dbUrl string) (*pgxpool.Pool, error) {
	if err := runMigrations(dbUrl); err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// runMigrations looks for a 'migrations' folder next to the application entrypoint.
func runMigrations(dbUrl string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	migration, err := migrate.New(fmt.Sprintf("file://%s/migrations", dir), dbUrl)
	if err != nil {
		return err
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
