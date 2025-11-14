package dao

import "time"

// PullRequest maps to 'pull_requests' table.
type PullRequest struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	AuthorID  string     `db:"author_id"`
	Status    string     `db:"status"`
	CreatedAt *time.Time `db:"created_at"`
	MergedAt  *time.Time `db:"merged_at"`
}
