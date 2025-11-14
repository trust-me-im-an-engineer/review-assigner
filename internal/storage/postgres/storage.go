package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"review-assigner/internal/config"
	"review-assigner/internal/storage"
)

var _ storage.Storage = (*Storage)(nil)

type activeTxKey int

var txContextKey activeTxKey = 0

type Storage struct {
	pool *pgxpool.Pool
}

// New initializes new postgres storage ready to work.
func New(ctx context.Context, cfg config.DBConfig) (*Storage, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name,
	)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	slog.Info("connected to postgres database")

	return &Storage{pool: pool}, nil
}

// Close must be called for graceful shutdown
func (s *Storage) Close() {
	if s.pool != nil {
		slog.Info("closing postgres pool connection...")
		s.pool.Close()
		slog.Info("postgres pool connection closed")
	}
}

func (s *Storage) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("postgres pool failed to begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, txContextKey, tx)

	var commitErr error
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		} else if commitErr != nil {
			_ = tx.Rollback(ctx)
		} else {
			if commitErr = tx.Commit(ctx); commitErr != nil {
				commitErr = fmt.Errorf("failed to commit transaction: %w", commitErr)
			}
		}
	}()

	commitErr = fn(txCtx)

	return commitErr
}
