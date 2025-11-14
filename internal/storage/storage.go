package storage

import (
	"context"

	"review-assigner/internal/model"
)

type Storage interface {
	AddTeam(ctx context.Context, name string) (string, error)
	AddUpdateUsers(ctx context.Context, users []model.User) ([]model.User, error)
	GetTeam(ctx context.Context, name string) (*model.Team, error)
	SetUserActivity(ctx context.Context, id string, active bool) (*model.User, error)
}
