package dao

// ReviewAssignment maps to 'review_assignments' junction table.
type ReviewAssignment struct {
	UserID        string `db:"user_id"`
	PullRequestID string `db:"pull_request_id"`
}
