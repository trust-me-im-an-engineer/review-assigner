package dao

import "time"

// PullRequest maps to 'pull_requests' table.
type PullRequest struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	AuthorID  string     `db:"author_id"` // Foreign key
	Status    string     `db:"status"`    // Stored as text (e.g., 'OPEN', 'MERGED')
	CreatedAt *time.Time `db:"created_at"`
	MergedAt  *time.Time `db:"merged_at"`
}
