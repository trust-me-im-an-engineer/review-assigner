package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"review-assigner/internal/errs"
	"review-assigner/internal/model"
)

// AddTeam inserts a new team into the database.
func (s *Storage) AddTeam(ctx context.Context, name string) (string, error) {
	q := `INSERT INTO teams (name) VALUES ($1)`

	var insertedName string
	err := s.getExecutor(ctx).QueryRow(ctx, q, name).Scan(&insertedName)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) && pgxError.Code == UniqueViolationErr {
			return "", errs.TeamExistsError{TeamName: name}
		}
		return "", fmt.Errorf("postgres failed to execute insert query for team: %w", err)
	}

	return insertedName, nil
}

func (s *Storage) GetTeam(ctx context.Context, name string) (*model.Team, error) {
	//TODO implement me
	panic("implement me")
}
