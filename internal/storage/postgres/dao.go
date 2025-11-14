package postgres

import (
	"time"
)

// teamDAO maps to 'teams' table.
type teamDAO struct {
	Name string `db:"name"`
}

// userDAO maps to 'users' table.
type userDAO struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	TeamName string `db:"team_name"`
	IsActive bool   `db:"is_active"`
}

// pullRequestDAO maps to 'pull_requests' table.
type pullRequestDAO struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	AuthorID  string     `db:"author_id"` // Foreign key
	Status    string     `db:"status"`    // Stored as text (e.g., 'OPEN', 'MERGED')
	CreatedAt *time.Time `db:"created_at"`
	MergedAt  *time.Time `db:"merged_at"`
}

// ReviewAssignmentDAO maps to 'review_assignments' junction table.
type ReviewAssignmentDAO struct {
	UserID        string `db:"user_id"`
	PullRequestID string `db:"pull_request_id"`
}
