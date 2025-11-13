package service

import (
	"context"
	"fmt"

	"review-assigner/internal/model"
	"review-assigner/internal/storage"
)

// Service is considered to be a core layer of witch only one could exist,
// therefore it doesn't use interface
type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) AddTeamAddUpdateUsers(ctx context.Context, team *model.Team) (*model.Team, error) {
	// schema requires user to always have team,
	// so user creation, updating and team assignment must be in one big method,
	// even tho it couples business logic to storage too tight.
	members, err := s.storage.AddTeamAddUpdateUsers(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("storage failed to add team add/update users: %w", err)
	}
	return members, nil
}

func (s *Service) GetTeam(ctx context.Context, name string) (*model.Team, error) {
	team, err := s.storage.GetTeam(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("storage failed to get team: %w", err)
	}
	return team, nil
}

func (s *Service) SetUserActivity(ctx context.Context, id string, active bool) (*model.User, error) {
	user, err := s.storage.SetUserActivity(ctx, id, active)
	if err != nil {
		return nil, fmt.Errorf("storage failed to set user activity: %w", err)
	}
	return user, nil
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
