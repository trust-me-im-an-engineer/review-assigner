package service

import (
	"context"

	"review-assigner/internal/model"
	"review-assigner/internal/storage"
)

// Service is considered to be a core layer of witch only one could exist,
// therefore it doesn't use interface
type Service struct {
	storage storage.Storage
}

func (s *Service) AddTeamAddUpdateUsers(ctx context.Context, m *model.Team) (*model.Team, error) {
	return nil, nil
}

func (s *Service) GetTeam(ctx context.Context, name string) (*model.Team, error) {
	return nil, nil
}

func (s *Service) SetUserActivity(ctx context.Context, id string, active bool) (*model.User, error) {
	return nil, nil
}

// Use model.PullRequestShort ignoring status
func (s *Service) CreatePullRequest(ctx context.Context, pr *model.PullRequestShort) (*model.PullRequest, error) {
	return nil, nil
}

func (s *Service) MergePullRequest(ctx context.Context, id string) (*model.PullRequest, error) {
	return nil, nil
}

func (s *Service) ReassignPullRequest(ctx context.Context, pullRequestID, oldReviewerID string) (pr *model.PullRequest, newReviewerID string, err error) {
	return nil, "", nil
}

func (s *Service) GetUserAssignments(ctx context.Context, id string) ([]model.PullRequestShort, error) {
	return nil, nil
}

func NewService(storage *storage.Storage) *Service {
	return &Service{storage: storage}
}
