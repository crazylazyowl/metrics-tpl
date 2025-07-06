package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/ping"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Options struct {
	DSN        string
	Migrations string
}

type PostgresStorage struct {
	db   *sql.DB
	opts Options
}

var _ ping.Pinger = (*PostgresStorage)(nil)
var _ metrics.MetricRegistry = (*PostgresStorage)(nil)

func NewPostgresStorage(opts Options) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", opts.DSN)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(opts.Migrations, "postgres", driver)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	return &PostgresStorage{db: db, opts: opts}, nil
}

func (s *PostgresStorage) Close(ctx context.Context) error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) Fetch(ctx context.Context) ([]metrics.Metric, error) {
	return nil, nil
}
func (s *PostgresStorage) FetchOne(ctx context.Context, m metrics.Metric) (metrics.Metric, error) {
	return metrics.Metric{}, nil
}
func (s *PostgresStorage) Update(ctx context.Context, m metrics.Metric) error {
	return nil
}
