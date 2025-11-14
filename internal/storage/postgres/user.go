package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"review-assigner/internal/model"
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

	daoUsers, err := pgx.CollectRows(rows, pgx.RowToStructByName[userDAO])

	result := make([]model.User, len(daoUsers))
	for i, daoUser := range daoUsers {
		result[i] = model.User{
			UserID:   daoUser.ID,
			Username: daoUser.Username,
			TeamName: daoUser.TeamName,
			IsActive: daoUser.IsActive,
		}
	}

	return result, nil
}

func (s *Storage) SetUserActivity(ctx context.Context, id string, active bool) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetActiveColleges(ctx context.Context, userID string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}
