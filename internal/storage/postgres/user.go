package postgres

import (
	"context"

	"review-assigner/internal/model"
)

func (s *Storage) AddUpdateUsers(ctx context.Context, users []model.User) ([]model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SetUserActivity(ctx context.Context, id string, active bool) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetActiveColleges(ctx context.Context, userID string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}
