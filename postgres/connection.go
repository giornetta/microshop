package postgres

import (
	"context"

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

func runMigrations(dbUrl string) error {
	migration, err := migrate.New("file:///home/michele/dev/microshop/cmd/products-service/migrations", dbUrl)
	if err != nil {
		return err
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
