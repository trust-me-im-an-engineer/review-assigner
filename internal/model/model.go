// Package model represents data in app.
package model

import (
	"time"
)

type PullRequestStatus string

const (
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
)

// TeamMember represents a user who is part of a team.
// Corresponds to #/components/schemas/TeamMember.
type TeamMember struct {
	UserID   string `json:"user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
	IsActive bool   `json:"is_active"`
}

// Team represents a collection of users.
// Corresponds to #/components/schemas/Team.
type Team struct {
	TeamName string       `json:"team_name" validate:"required"`
	Members  []TeamMember `json:"members" validate:"required,dive"`
}

// User represents an individual user with their team and activity status.
// Corresponds to #/components/schemas/User.
type User struct {
	UserID   string `json:"user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
	TeamName string `json:"team_name" validate:"required"`
	IsActive bool   `json:"is_active"`
}

// PullRequestShort provides a basic, short representation of a pull request.
// Corresponds to #/components/schemas/PullRequestShort.
type PullRequestShort struct {
	PullRequestID   string            `json:"pull_request_id" validate:"required"`
	PullRequestName string            `json:"pull_request_name" validate:"required"`
	AuthorID        string            `json:"author_id" validate:"required"`
	Status          PullRequestStatus `json:"status" validate:"required,oneof=OPEN MERGED"`
}

// PullRequest represents a full pull request object, including assigned reviewers
// and timestamps.
// Corresponds to #/components/schemas/PullRequest.
type PullRequest struct {
	PullRequestID   string            `json:"pull_request_id" validate:"required"`
	PullRequestName string            `json:"pull_request_name" validate:"required"`
	AuthorID        string            `json:"author_id" validate:"required"`
	Status          PullRequestStatus `json:"status" validate:"required,oneof=OPEN MERGED"`
	// Max 2 reviewers are assigned, as per API description/logic
	AssignedReviewers []string   `json:"assigned_reviewers" validate:"max=2"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}
