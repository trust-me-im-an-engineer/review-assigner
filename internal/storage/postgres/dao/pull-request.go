package dao

import (
	"time"

	"review-assigner/internal/model"
)

// PullRequest maps to 'pull_requests' table.
type PullRequest struct {
	ID        string                  `db:"id"`
	Name      string                  `db:"name"`
	AuthorID  string                  `db:"author_id"`
	Status    model.PullRequestStatus `db:"status"`
	CreatedAt *time.Time              `db:"created_at"`
	MergedAt  *time.Time              `db:"merged_at"`
}
