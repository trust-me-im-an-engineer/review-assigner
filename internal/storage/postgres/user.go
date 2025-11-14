package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"review-assigner/internal/errs"
	"review-assigner/internal/model"
	"review-assigner/internal/storage/postgres/dao"
)

func (s *Storage) AddUpdateUsers(ctx context.Context, users []model.User) ([]model.User, error) {
	vals := make([]any, 0, len(users))
	for _, user := range users {
		vals = append(vals, []any{user.UserID, user.Username, user.TeamName, user.IsActive})
	}

	builder := squirrelBuilder.Insert("users").
		Columns("id", "username", "team_name", "is_active").
		Values(vals...).
		Suffix(`ON CONFLICT (id) DO UPDATE SET 
            username = EXCLUDED.username,
            team_name = EXCLUDED.team_name,
            is_active = EXCLUDED.is_active
			RETURNING *`)

	query, vals, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("squirrel failed to build query: %w", err)
	}

	rows, err := s.getExecutor(ctx).Query(ctx, query, vals...)
	if err != nil {
		return nil, fmt.Errorf("postgres failed to execute query: %w", err)
	}
	defer rows.Close()

	daoUsers, err := pgx.CollectRows(rows, pgx.RowToStructByName[dao.User])
	if err != nil {
		return nil, fmt.Errorf("pgx failed to collect rows: %w", err)
	}

	result := make([]model.User, len(daoUsers))
	for i, daoUser := range daoUsers {
		result[i] = daoUser.ToModel()
	}

	return result, nil
}

func (s *Storage) SetUserActivity(ctx context.Context, id string, active bool) (*model.User, error) {
	q := `UPDATE users SET is_active = $1 WHERE id = $2 RETURNING *`
	rows, err := s.getExecutor(ctx).Query(ctx, q, active, id)
	if err != nil {
		return nil, fmt.Errorf("postgres failed to esecute query: %w", err)
	}
	defer rows.Close()

	daoUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dao.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFoundErr
		}
		return nil, fmt.Errorf("pgx failed to collect one row: %w", err)
	}

	user := daoUser.ToModel()

	return &user, nil
}

func (s *Storage) GetActiveColleges(ctx context.Context, userID string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}
