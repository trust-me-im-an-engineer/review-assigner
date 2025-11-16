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
	"review-assigner/internal/rest/wjson"
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
		wjson.Error(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		wjson.Error(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNotFound)
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
			wjson.Error(w, teamErr.Error(), http.StatusConflict, payload.ErrCodeTEAMExists)
			return
		}
		slog.Error("service failed to add team and add/update users", "error", err)
		wjson.Error(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	wjson.Response(w, map[string]*model.Team{"team": team}, http.StatusCreated)
}

// GetTeam handles GET /team/get
func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		wjson.Error(w, "missing query parameter 'team_name'", http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if len(teamName) > 255 {
		wjson.Error(w, "team_name cannot be longer than 255 symbols", http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	team, err := h.service.GetTeam(r.Context(), teamName)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			wjson.Error(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		slog.Error("service failed to get team", "error", err)
		wjson.Error(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	wjson.Response(w, team, http.StatusOK)
}

// SetUserActivity handles POST /users/setIsActive
func (h *Handler) SetUserActivity(w http.ResponseWriter, r *http.Request) {
	var req payload.SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		wjson.Error(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		wjson.Error(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	user, err := h.service.SetUserActivity(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			wjson.Error(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		slog.Error("service failed to set user activity", "error", err)
		wjson.Error(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	wjson.Response(w, map[string]*model.User{"user": user}, http.StatusOK)
}

// CreatePullRequest handles POST /pullRequest/create
func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var req payload.PullRequestCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		wjson.Error(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		wjson.Error(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNotFound)
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
			wjson.Error(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		var prErr errs.PullRequestExistsError
		if errors.As(err, &prErr) {
			wjson.Error(w, prErr.Error(), http.StatusConflict, payload.ErrCodePRExists)
			return
		}
		slog.Error("service failed to create pull request", "error", err)
		wjson.Error(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	wjson.Response(w, map[string]*model.PullRequest{"pr": pr}, http.StatusCreated)
}

// MergePullRequest handles POST /pullRequest/merge
func (h *Handler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	var req payload.PullRequestMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		wjson.Error(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		wjson.Error(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	pr, err := h.service.MergePullRequest(r.Context(), req.PullRequestID)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			wjson.Error(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		slog.Error("service failed to merge pull request", "error", err)
		wjson.Error(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	wjson.Response(w, map[string]*model.PullRequest{"pr": pr}, http.StatusOK)
}

// ReassignPullRequest handles POST /pullRequest/reassign
func (h *Handler) ReassignPullRequest(w http.ResponseWriter, r *http.Request) {
	var req payload.PullRequestReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		wjson.Error(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		wjson.Error(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	pr, newReviewerID, err := h.service.ReassignPullRequest(r.Context(), req.PullRequestID, req.OldReviewerID)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			wjson.Error(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		if errors.Is(err, errs.PullRequestMergedErr) {
			wjson.Error(w, errs.PullRequestMergedErr.Error(), http.StatusConflict, payload.ErrCodePRMerged)
			return
		}
		if errors.Is(err, errs.NotAssignedErr) {
			wjson.Error(w, errs.NotAssignedErr.Error(), http.StatusConflict, payload.ErrCodeNotAssigned)
			return
		} else if errors.Is(err, errs.NoCandidateErr) {
			wjson.Error(w, errs.NoCandidateErr.Error(), http.StatusConflict, payload.ErrCodeNoCandidate)
			return
		}
		slog.Error("service failed to reassign pull request", "error", err)
		wjson.Error(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	wjson.Response(w, map[string]any{"pr": pr, "replaced_by": newReviewerID}, http.StatusOK)
}

// GetUserAssignments handles GET /users/getReview
func (h *Handler) GetUserAssignments(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		wjson.Error(w, "missing query parameter 'user_id'", http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}
	if len(userID) > 255 {
		wjson.Error(w, "user_id cannot be longer than 255 symbols", http.StatusBadRequest, payload.ErrCodeNotFound)
		return
	}

	pullRequests, err := h.service.GetUserAssignments(r.Context(), userID)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			wjson.Error(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNotFound)
			return
		}
		slog.Error("service failed to get user assignments", "user_id", userID, "error", err)
		wjson.Error(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNotFound)
		return
	}

	response := payload.GetUserReviewResponse{
		UserID:       userID,
		PullRequests: pullRequests,
	}

	wjson.Response(w, response, http.StatusOK)
}
