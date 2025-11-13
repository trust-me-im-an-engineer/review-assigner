package storage

import (
	"context"

	"review-assigner/internal/model"
)

type Storage interface {
	AddTeamAddUpdateUsers(ctx context.Context, team *model.Team) (*model.Team, error)
	GetTeam(name string) (*model.Team, error)
}
