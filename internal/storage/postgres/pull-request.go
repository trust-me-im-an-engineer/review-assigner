package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"review-assigner/internal/errs"
	"review-assigner/internal/model"
	"review-assigner/internal/storage/postgres/dao"
)

func (s *Storage) CreatePullRequestWithAssignments(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error) {
	var createdPR model.PullRequest

	err := s.InTransaction(ctx, func(ctx context.Context) error {
		e := s.getExecutor(ctx)

		// Insert PullRequest
		qPR := `
			INSERT INTO pull_requests (id, name, author_id, status, created_at, merged_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING *;
		`

		rowsPR, err := e.Query(
			ctx,
			qPR,
			pr.Id,
			pr.Name,
			pr.AuthorID,
			pr.Status,
			pr.CreatedAt,
			pr.MergedAt,
		)
		if err != nil {
			return fmt.Errorf("postgres failed to execute insert query for pull request: %w", err)
		}
		defer rowsPR.Close()

		daoPR, err := pgx.CollectOneRow(rowsPR, pgx.RowToStructByName[dao.PullRequest])
		if err != nil {
			var pgxError *pgconn.PgError
			if errors.As(err, &pgxError) && pgxError.Code == UniqueViolationErr {
				return errs.PullRequestExistsError{PullRequestID: pr.Id}
			}
			return fmt.Errorf("postgres failed to collect dao pull request row: %w", err)
		}

		// Build VALUES for review_assignments
		if len(pr.AssignedReviewers) == 0 {
			// no assignments, return PR as-is
			createdPR = model.PullRequest{
				Id:                daoPR.ID,
				Name:              daoPR.Name,
				AuthorID:          daoPR.AuthorID,
				Status:            daoPR.Status,
				AssignedReviewers: []string{},
				CreatedAt:         daoPR.CreatedAt,
				MergedAt:          daoPR.MergedAt,
			}
			return nil
		}

		valuePlaceholders := make([]string, len(pr.AssignedReviewers))
		args := make([]any, 0, len(pr.AssignedReviewers)*2)

		argIndex := 1
		for i, reviewer := range pr.AssignedReviewers {
			// each row: (user_id, pull_request_id)
			valuePlaceholders[i] = fmt.Sprintf("($%d, $%d)", argIndex, argIndex+1)

			args = append(args, reviewer, pr.Id)
			argIndex += 2
		}

		qAssignments := fmt.Sprintf(`
			INSERT INTO review_assignments (user_id, pull_request_id)
			VALUES %s
			RETURNING *;
		`, strings.Join(valuePlaceholders, ","))

		rowsAssignments, err := e.Query(ctx, qAssignments, args...)
		if err != nil {
			return fmt.Errorf("postgres failed to execute insert query for review assignments: %w", err)
		}
		defer rowsAssignments.Close()

		daoAssignments, err := pgx.CollectRows(rowsAssignments, pgx.RowToStructByName[dao.ReviewAssignment])
		if err != nil {
			return fmt.Errorf("postgres failed to collect dao assignments: %w", err)
		}

		assignedReviewers := make([]string, 0, len(daoAssignments))
		for _, assignment := range daoAssignments {
			assignedReviewers = append(assignedReviewers, assignment.UserID)
		}

		createdPR = model.PullRequest{
			Id:                daoPR.ID,
			Name:              daoPR.Name,
			AuthorID:          daoPR.AuthorID,
			Status:            daoPR.Status,
			AssignedReviewers: assignedReviewers,
			CreatedAt:         daoPR.CreatedAt,
			MergedAt:          daoPR.MergedAt,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &createdPR, nil
}

func (s *Storage) GetPullRequest(ctx context.Context, id string) (*model.PullRequest, error) {
	var result model.PullRequest
	err := s.InTransaction(ctx, func(ctx context.Context) error {
		e := s.getExecutor(ctx)

		qPR := `SELECT * FROM pull_requests WHERE id = $1`
		rowsPR, err := e.Query(ctx, qPR, id)
		if err != nil {

			return fmt.Errorf("postgres failed to get pull request: %w", err)
		}
		defer rowsPR.Close()

		daoPR, err := pgx.CollectOneRow(rowsPR, pgx.RowToStructByName[dao.PullRequest])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errs.NotFoundErr
			}
			return fmt.Errorf("postgres failed to collect one row: %w", err)
		}

		result = model.PullRequest{
			Id:        daoPR.ID,
			Name:      daoPR.Name,
			AuthorID:  daoPR.AuthorID,
			Status:    daoPR.Status,
			CreatedAt: daoPR.CreatedAt,
			MergedAt:  daoPR.MergedAt,
		}

		qAssignments := `SELECT * FROM review_assignments WHERE pull_request_id = $1`
		rowsAssignments, err := e.Query(ctx, qAssignments, id)
		if err != nil {
			return fmt.Errorf("postgres failed to get review assignments: %w", err)
		}
		defer rowsAssignments.Close()

		daoAssignments, err := pgx.CollectRows(rowsAssignments, pgx.RowToStructByName[dao.ReviewAssignment])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errs.NotFoundErr
			}
			return fmt.Errorf("postgres failed to collect rows: %w", err)
		}

		assignments := make([]string, len(daoAssignments))
		for i, daoAssignment := range daoAssignments {
			assignments[i] = daoAssignment.UserID
		}

		result.AssignedReviewers = assignments

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Storage) UpdatePullRequest(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error) {
	q := `UPDATE pull_requests
		  SET name = $2, author_id = $3, status = $4, created_at = $5, merged_at = $6
		  WHERE id = $1 RETURNING *`
	rows, err := s.getExecutor(ctx).Query(ctx, q, pr.Id, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt)
	if err != nil {
		return nil, fmt.Errorf("postgres failed to execute update pull request query: %w", err)
	}
	defer rows.Close()

	daoPR, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dao.PullRequest])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFoundErr
		}
		return nil, fmt.Errorf("postgres failed to collect one row: %w", err)
	}

	updatedPR := &model.PullRequest{
		Id:                daoPR.ID,
		Name:              daoPR.Name,
		AuthorID:          daoPR.AuthorID,
		Status:            daoPR.Status,
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         daoPR.CreatedAt,
		MergedAt:          daoPR.MergedAt,
	}

	return updatedPR, nil
}
