// Package model represents data in app.
package model

import "time"

type PullRequestStatus string

const (
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
)

// TeamMember represents a user who is part of a team.
// Corresponds to #/components/schemas/TeamMember.
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// Team represents a collection of users.
// Corresponds to #/components/schemas/Team.
type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

// User represents an individual user with their team and activity status.
// Corresponds to #/components/schemas/User.
type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// PullRequestShort provides a basic, short representation of a pull request.
// Corresponds to #/components/schemas/PullRequestShort.
type PullRequestShort struct {
	PullRequestID   string            `json:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name"`
	AuthorID        string            `json:"author_id"`
	Status          PullRequestStatus `json:"status"`
}

// PullRequest represents a full pull request object, including assigned reviewers
// and timestamps.
// Corresponds to #/components/schemas/PullRequest.
type PullRequest struct {
	PullRequestID     string            `json:"pull_request_id"`
	PullRequestName   string            `json:"pull_request_name"`
	AuthorID          string            `json:"author_id"`
	Status            PullRequestStatus `json:"status"`
	AssignedReviewers []string          `json:"assigned_reviewers"`
	CreatedAt         *time.Time        `json:"createdAt,omitempty"`
	MergedAt          *time.Time        `json:"mergedAt,omitempty"`
}
