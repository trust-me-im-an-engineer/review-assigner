package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"review-assigner/internal/errs"
	"review-assigner/internal/model"
	"review-assigner/internal/storage/postgres/dao"
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

// GetTeam retrieves a team by name, including all its members.
func (s *Storage) GetTeam(ctx context.Context, name string) (*model.Team, error) {
	var team model.Team
	err := s.InTransaction(ctx, func(ctx context.Context) error {
		e := s.getExecutor(ctx)

		qTeam := `SELECT name FROM teams WHERE name = $1`
		var teamName string
		err := e.QueryRow(ctx, qTeam, name).Scan(&teamName)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errs.NotFoundErr
			}
			return fmt.Errorf("postgres failed to query team: %w", err)
		}

		qMembers := `SELECT id, username, is_active FROM users WHERE team_name = $1`
		rows, err := e.Query(ctx, qMembers, teamName)
		if err != nil {
			return fmt.Errorf("postgres failed to query team members: %w", err)
		}
		defer rows.Close()

		daoMembers, err := pgx.CollectRows(rows, pgx.RowToStructByName[dao.Member])
		if err != nil {
			return fmt.Errorf("pgx failed to collect team daoMember rows: %w", err)
		}

		members := make([]model.TeamMember, len(daoMembers))
		for i, daoMember := range daoMembers {
			members[i] = daoMember.ToModel()
		}

		team = model.Team{
			Name:    teamName,
			Members: members,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &team, nil
}
