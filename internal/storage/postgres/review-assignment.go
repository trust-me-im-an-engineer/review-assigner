package postgres

import (
	"context"

	"review-assigner/internal/model"
)

func (s *Storage) DeleteReviewAssignment(ctx context.Context, prID string, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) AddReviewAssignment(ctx context.Context, prID string, userID string) (reviewerID string, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetUserAssignments(ctx context.Context, userID string) ([]model.PullRequestShort, error) {
	//TODO implement me
	panic("implement me")
}
