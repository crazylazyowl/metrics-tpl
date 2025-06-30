package postgres

import (
	"context"
	"database/sql"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/ping"

	_ "github.com/lib/pq"
)

type Options struct {
	DNS string
}

type PostgresStorage struct {
	db   *sql.DB
	opts Options
}

var _ ping.Pinger = (*PostgresStorage)(nil)

func NewPostgresStorage(opts Options) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", opts.DNS)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Close() error {
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
