package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	databaseURL string
	pool        *pgxpool.Pool
}

func New(databaseURL string) (*Storage, error) {
	const op = "storage.postgres.New"

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pingContext, cancelPing := context.WithTimeout(ctx, time.Second*2)
	defer cancelPing()

	if err := pool.Ping(pingContext); err != nil {
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	storage := &Storage{
		databaseURL: databaseURL,
		pool:        pool,
	}

	return storage, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) Migrate(migrationsFolder string) error {
	const op = "storage.postgres.Migrate"

	databaseURL := strings.Replace(s.databaseURL, "postgres", "pgx5", 1)

	sourceURL := fmt.Sprintf("file:/%s", migrationsFolder)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No new migrations to apply.")
		} else {
			return fmt.Errorf("%s: failed to run migrations: %w", op, err)
		}
	} else {
		fmt.Println("Migration applied successfully.")
	}

	return nil
}
