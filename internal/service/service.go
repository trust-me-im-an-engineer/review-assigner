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

func (s *Service) AddTeamAddUpdateUsers(context context.Context, m *model.Team) (*model.Team, error) {
	return nil, nil
}

func (s *Service) GetTeam(name string) (*model.Team, error) {
	return nil, nil
}

func NewService(storage *storage.Storage) *Service {
	return &Service{storage: storage}
}
