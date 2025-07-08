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

	_ "github.com/jackc/pgx/v5"
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
	stmt, err := s.db.PrepareContext(ctx, `SELECT name, type, counter, gauge FROM metrics;`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	mm := make([]metrics.Metric, 0) // TODO: it's not good
	for rows.Next() {
		var m metrics.Metric
		if err := rows.Scan(&m.ID, &m.Type, &m.Counter, &m.Gauge); err != nil {
			return nil, err
		}
		mm = append(mm, m)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return mm, nil
}

func (s *PostgresStorage) FetchOne(ctx context.Context, m metrics.Metric) (metrics.Metric, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT counter, gauge FROM metrics WHERE name = $1 AND type = $2;`)
	if err != nil {
		return metrics.Metric{}, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(ctx, m.ID, m.Type)
	if err := row.Scan(&m.Counter, &m.Gauge); err != nil {
		switch err {
		case sql.ErrNoRows:
			return metrics.Metric{}, metrics.ErrMetricNotFound
		}
		return metrics.Metric{}, err
	}
	return m, nil
}

func (s *PostgresStorage) UpdateOne(ctx context.Context, m metrics.Metric) error {
	query, err := selectUpsertQuery(m)
	if err != nil {
		return err
	}
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, m.ID, m.Type, m.Counter, m.Gauge)
	return err
}

func (s *PostgresStorage) Update(ctx context.Context, mm []metrics.Metric) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, m := range mm {
		query, err := selectUpsertQuery(m)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, query, m.ID, m.Type, m.Counter, m.Gauge); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func selectUpsertQuery(m metrics.Metric) (query string, err error) {
	switch m.Type {
	case metrics.Counter:
		query = `
			INSERT INTO metrics (name, type, counter, gauge) VALUES ($1, $2, $3, $4) 
			ON CONFLICT (name, type) DO 
				UPDATE SET counter = metrics.counter + EXCLUDED.counter;
		`
	case metrics.Gauge:
		query = `
			INSERT INTO metrics (name, type, counter, gauge) VALUES ($1, $2, $3, $4) 
			ON CONFLICT (name, type) DO 
				UPDATE SET gauge = EXCLUDED.gauge;
		`
	default:
		err = metrics.ErrMetricUnknownType
	}
	return
}
