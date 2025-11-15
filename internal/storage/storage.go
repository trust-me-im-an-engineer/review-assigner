package storage

import (
	"context"

	"review-assigner/internal/model"
)

type Storage interface {
	Team
	User
	PullRequest

	// InTransaction executes given function in a transaction.
	// The transaction will be committed if fn returns nil, or rolled back otherwise.
	// The context passed to fn should be used for all underlying database operations.
	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Team interface {
	AddTeam(ctx context.Context, name string) (string, error)
	GetTeam(ctx context.Context, name string) (*model.Team, error)
}

type User interface {
	AddUpdateUsers(ctx context.Context, users []model.User) ([]model.User, error)
	SetUserActivity(ctx context.Context, id string, active bool) (*model.User, error)

	// GetActiveColleges returns userIDs of users in the same team as userID excluding userID itself.
	GetActiveColleges(ctx context.Context, userID string) ([]string, error)
}

type PullRequest interface {
	CreatePullRequestWithAssignments(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error)
	GetPullRequest(ctx context.Context, id string) (*model.PullRequest, error)
	UpdatePullRequest(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error)

	// GetUserAssignments returns pull requests where user is one of reviewers
	GetUserAssignments(ctx context.Context, userID string) ([]model.PullRequestShort, error)
}
