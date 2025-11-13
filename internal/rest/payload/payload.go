// Package payload represents rest api specific data: requests and responses.
package payload

import "review-assigner/internal/model"

// TeamAddRequest corresponds to the /team/add POST request body.
// Validation is applied via embedded model.Team structure.
type TeamAddRequest model.Team

// SetIsActiveRequest corresponds to the /users/setIsActive POST request body.
type SetIsActiveRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	IsActive bool   `json:"is_active"`
}

// PullRequestCreateRequest corresponds to the /pullRequest/create POST request body.
type PullRequestCreateRequest struct {
	PullRequestID   string `json:"pull_request_id" validate:"required"`
	PullRequestName string `json:"pull_request_name" validate:"required"`
	AuthorID        string `json:"author_id" validate:"required"`
}

// PullRequestMergeRequest corresponds to the /pullRequest/merge POST request body.
type PullRequestMergeRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required"`
}

// PullRequestReassignRequest corresponds to the /pullRequest/reassign POST request body.
type PullRequestReassignRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required"`
	OldReviewerID string `json:"old_user_id" validate:"required"`
}

// GetUserReviewResponse corresponds to the /users/getReview GET response.
// As this is a response payload, validation tags are typically omitted.
type GetUserReviewResponse struct {
	UserID       string                   `json:"user_id"`
	PullRequests []model.PullRequestShort `json:"pull_requests"`
}

// ErrorResponse represents a standardized error response from the API.
// Corresponds to #/components/schemas/ErrorResponse.
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
