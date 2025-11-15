package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"review-assigner/internal/errs"
	"review-assigner/internal/model"
	"review-assigner/internal/rest/payload"
	"review-assigner/internal/service"
	"review-assigner/internal/validator"
)

const (
	invalidJsonBodyMsg     = "invalid JSON body"
	internalServerErrorMsg = "internal server error"
)

// Handler contains handlers for rest api.
// Note that handler is forced to return semantically incorrect error codes to meet openapi specs.
type Handler struct {
	service  *service.Service
	validate *validator.Validator
}

func NewHandler(service *service.Service, validate *validator.Validator) *Handler {
	return &Handler{service: service, validate: validate}
}

// AddTeamAddUpdateUsers handles POST /team/add
func (h *Handler) AddTeamAddUpdateUsers(w http.ResponseWriter, r *http.Request) {
	var req payload.TeamAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	team, err := h.service.AddTeamAddUpdateUsers(r.Context(), &model.Team{
		Name:    req.Name,
		Members: req.Members,
	})
	if err != nil {
		var teamErr errs.TeamExistsError
		if errors.As(err, &teamErr) {
			slog.Warn("team already exists on add", "team_name", req.Name, "error", teamErr)
			writeJSONError(w, teamErr.Error(), http.StatusConflict, payload.ErrCodeTEAMExists)
			return
		}
		slog.Error("service failed to add team and add/update users", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	writeJSONResponse(w, map[string]*model.Team{"team": team}, http.StatusCreated)
}

// GetTeam handles GET /team/get
func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeJSONError(w, "missing query parameter 'team_name'", http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if len(teamName) > 255 {
		writeJSONError(w, "team_name cannot be longer than 255 symbols", http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	team, err := h.service.GetTeam(r.Context(), teamName)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			writeJSONError(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		slog.Error("service failed to get team", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	writeJSONResponse(w, team, http.StatusOK)
}

// SetUserActivity handles POST /users/setIsActive
func (h *Handler) SetUserActivity(w http.ResponseWriter, r *http.Request) {
	var req payload.SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	user, err := h.service.SetUserActivity(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			writeJSONError(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		slog.Error("service failed to set user activity", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	writeJSONResponse(w, map[string]*model.User{"user": user}, http.StatusOK)
}

// CreatePullRequest handles POST /pullRequest/create
func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var req payload.PullRequestCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	// pull request status is ignored
	pr, err := h.service.CreatePullRequest(r.Context(), &model.PullRequestShort{
		Id:       req.PullRequestID,
		Name:     req.PullRequestName,
		AuthorID: req.AuthorID,
	})
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			writeJSONError(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		var prErr errs.PullRequestExistsError
		if errors.As(err, &prErr) {
			writeJSONError(w, prErr.Error(), http.StatusConflict, payload.ErrCodePRExists)
			return
		}
		slog.Error("service failed to create pull request", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	writeJSONResponse(w, map[string]*model.PullRequest{"pr": pr}, http.StatusCreated)
}

// MergePullRequest handles POST /pullRequest/merge
func (h *Handler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	var req payload.PullRequestMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	pr, err := h.service.MergePullRequest(r.Context(), req.PullRequestID)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			writeJSONError(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		slog.Error("service failed to merge pull request", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	writeJSONResponse(w, map[string]*model.PullRequest{"pr": pr}, http.StatusOK)
}

// ReassignPullRequest handles POST /pullRequest/reassign
func (h *Handler) ReassignPullRequest(w http.ResponseWriter, r *http.Request) {
	var req payload.PullRequestReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	pr, newReviewerID, err := h.service.ReassignPullRequest(r.Context(), req.PullRequestID, req.OldReviewerID)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			writeJSONError(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		if errors.Is(err, errs.PullRequestMergedErr) {
			writeJSONError(w, errs.PullRequestMergedErr.Error(), http.StatusConflict, payload.ErrCodePRMerged)
			return
		}
		if errors.Is(err, errs.NotAssignedErr) {
			writeJSONError(w, errs.NotAssignedErr.Error(), http.StatusConflict, payload.ErrCodeNotAssigned)
			return
		} else if errors.Is(err, errs.NoCandidateErr) {
			writeJSONError(w, errs.NoCandidateErr.Error(), http.StatusConflict, payload.ErrCodeNoCandidate)
			return
		}
		slog.Error("service failed to reassign pull request", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	writeJSONResponse(w, map[string]any{"pr": pr, "replaced_by": newReviewerID}, http.StatusOK)
}

// GetUserAssignments handles GET /users/getReview
func (h *Handler) GetUserAssignments(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeJSONError(w, "missing query parameter 'user_id'", http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if len(userID) > 255 {
		writeJSONError(w, "user_id cannot be longer than 255 symbols", http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	pullRequests, err := h.service.GetUserAssignments(r.Context(), userID)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			writeJSONError(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		slog.Error("service failed to get user assignments", "user_id", userID, "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	response := payload.GetUserReviewResponse{
		UserID:       userID,
		PullRequests: pullRequests,
	}

	writeJSONResponse(w, response, http.StatusOK)
}

func writeJSONError(w http.ResponseWriter, msg string, statusCode int, apiCode payload.ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := payload.ErrorResponse{
		Error: payload.InnerError{
			Code:    apiCode,
			Message: msg,
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to write JSON response", "error", err)
	}
}

func writeJSONResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to write JSON response", "error", err)
	}
}
