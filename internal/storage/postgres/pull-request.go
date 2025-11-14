package postgres

import (
	"context"

	"review-assigner/internal/model"
)

func (s *Storage) CreatePullRequest(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetPullRequest(ctx context.Context, id string) (*model.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) UpdatePullRequest(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetUserAssignments(ctx context.Context, userID string) ([]model.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}
