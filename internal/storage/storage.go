package storage

import (
	"context"

	"review-assigner/internal/model"
)

type Storage interface {
	AddTeamAddUpdateUsers(ctx context.Context, team *model.Team) (*model.Team, error)
	GetTeam(ctx context.Context, name string) (*model.Team, error)
	SetUserActivity(ctx context.Context, id string, active bool) (*model.User, error)
}
