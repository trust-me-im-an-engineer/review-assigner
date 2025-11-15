package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"review-assigner/internal/model"
	"review-assigner/internal/storage/postgres/dao"
)

func (s *Storage) DeleteReviewAssignment(ctx context.Context, prID string, userID string) error {
	q := `DELETE FROM review_assignments WHERE pull_request_id = $1 AND user_id = $2`
	_, err := s.getExecutor(ctx).Exec(ctx, q, prID, userID)
	if err != nil {
		return fmt.Errorf("postgres failed to delete review assignment: %w", err)
	}
	return nil
}

func (s *Storage) AddReviewAssignment(ctx context.Context, prID string, userID string) (reviewerID string, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetUserAssignments(ctx context.Context, userID string) ([]model.PullRequestShort, error) {
	q := `SELECT id, name, author_id, status FROM pull_requests
		  WHERE id = 
		        (SELECT pull_request_id FROM review_assignments WHERE user_id = $1)`
	rows, err := s.getExecutor(ctx).Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("postgres failed to execute get user assignments query: %w", err)
	}
	defer rows.Close()

	daoPRs, err := pgx.CollectRows(rows, pgx.RowToStructByName[dao.PullRequestShort])
	if err != nil {
		return nil, fmt.Errorf("postgres failed to collect rows: %w", err)
	}

	prs := make([]model.PullRequestShort, len(daoPRs))
	for i, daoPR := range daoPRs {
		prs[i] = model.PullRequestShort{
			Id:       daoPR.ID,
			Name:     daoPR.Name,
			AuthorID: daoPR.AuthorID,
			Status:   daoPR.Status,
		}
	}

	return prs, nil
}
