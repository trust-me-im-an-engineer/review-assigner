package postgres

import (
	"context"

	"review-assigner/internal/model"
)

func (s *Storage) AddTeam(ctx context.Context, name string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetTeam(ctx context.Context, name string) (*model.Team, error) {
	//TODO implement me
	panic("implement me")
}
