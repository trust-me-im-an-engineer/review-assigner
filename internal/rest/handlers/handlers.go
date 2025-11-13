package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	"review-assigner/internal/errs"
	"review-assigner/internal/model"
	"review-assigner/internal/rest/payload"
	"review-assigner/internal/service"
)

const (
	invalidJsonBodyMsg     = "invalid JSON body"
	internalServerErrorMsg = "internal server error"
)

// Handler contains handlers for rest api.
// Note that handler is forced to return semantically incorrect error codes to meet openapi specs.
type Handler struct {
	service  *service.Service
	validate *validator.Validate
}

func NewHandler(service *service.Service, validate *validator.Validate) *Handler {
	return &Handler{service: service, validate: validate}
}

// AddTeamAddUpdateUsers handles POST /team/add
func (h *Handler) AddTeamAddUpdateUsers(w http.ResponseWriter, r *http.Request) {
	var req payload.TeamAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNOT_FOUND)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNOT_FOUND)
		return
	}

	team, err := h.service.AddTeamAddUpdateUsers(r.Context(), &model.Team{
		TeamName: req.TeamName,
		Members:  req.Members,
	})
	if err != nil {
		var teamErr errs.TeamExistsError
		if errors.As(err, &teamErr) {
			slog.Warn("team already exists on add", "team_name", req.TeamName, "error", teamErr)
			writeJSONError(w, teamErr.Error(), http.StatusConflict, payload.ErrCodeTEAM_EXISTS)
			return
		}
		slog.Error("service failed to add team and add/update users", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNOT_FOUND)
		return
	}

	writeJSONResponse(w, map[string]*model.Team{"team": team}, http.StatusCreated)
}

// GetTeam handles GET /team/get
func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeJSONError(w, "missing query parameter 'team_name'", http.StatusBadRequest, payload.ErrCodeNOT_FOUND)
		return
	}

	team, err := h.service.GetTeam(teamName)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			writeJSONError(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNOT_FOUND)
			return
		}
		slog.Error("service failed to get team", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNOT_FOUND)
		return
	}

	writeJSONResponse(w, team, http.StatusOK)
}

// SetUserActivity handles POST /users/setIsActive
func (h *Handler) SetUserActivity(w http.ResponseWriter, r *http.Request) {
	var req payload.SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNOT_FOUND)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNOT_FOUND)
		return
	}

	user, err := h.service.SetUserActivity(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			writeJSONError(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNOT_FOUND)
			return
		}
		slog.Error("service failed to set user activity", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNOT_FOUND)
		return
	}

	writeJSONResponse(w, map[string]*model.User{"user": user}, http.StatusOK)
}

// CreatePullRequest handles POST /pullRequest/create
func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var req payload.PullRequestCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNOT_FOUND)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNOT_FOUND)
		return
	}

	// pull request status is ignored
	pr, err := h.service.CreatePullRequest(r.Context(), &model.PullRequestShort{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
	})
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			writeJSONError(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNOT_FOUND)
			return
		}
		var prErr errs.PullRequestExistsError
		if errors.As(err, &prErr) {
			writeJSONError(w, prErr.Error(), http.StatusConflict, payload.ErrCodePR_EXISTS)
			return
		}
		slog.Error("service failed to create pull request", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNOT_FOUND)
		return
	}

	writeJSONResponse(w, map[string]*model.PullRequest{"pr": pr}, http.StatusCreated)
}

// MergePullRequest handles POST /pullRequest/merge
func (h *Handler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	var req payload.PullRequestMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNOT_FOUND)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNOT_FOUND)
		return
	}

	pr, err := h.service.MergePullRequest(r.Context(), req.PullRequestID)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			writeJSONError(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNOT_FOUND)
			return
		}
		slog.Error("service failed to merge pull request", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNOT_FOUND)
		return
	}

	writeJSONResponse(w, map[string]*model.PullRequest{"pr": pr}, http.StatusOK)
}

// ReassignPullRequest handles POST /pullRequest/reassign
func (h *Handler) ReassignPullRequest(w http.ResponseWriter, r *http.Request) {
	var req payload.PullRequestReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, invalidJsonBodyMsg, http.StatusBadRequest, payload.ErrCodeNOT_FOUND)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest, payload.ErrCodeNOT_FOUND)
		return
	}

	pr, newReviewerID, err := h.service.ReassignPullRequest(r.Context(), req.PullRequestID, req.OldReviewerID)
	if err != nil {
		if errors.Is(err, errs.NotFoundErr) {
			writeJSONError(w, errs.NotFoundErr.Error(), http.StatusNotFound, payload.ErrCodeNOT_FOUND)
			return
		}
		if errors.Is(err, errs.PullRequestMergedErr) {
			writeJSONError(w, errs.PullRequestMergedErr.Error(), http.StatusConflict, payload.ErrCodePR_MERGED)
			return
		}
		if errors.Is(err, errs.NotAssignedErr) {
			writeJSONError(w, errs.NotAssignedErr.Error(), http.StatusConflict, payload.ErrCodeNOT_ASSIGNED)
			return
		} else if errors.Is(err, errs.NoCandidateErr) {
			writeJSONError(w, errs.NoCandidateErr.Error(), http.StatusConflict, payload.ErrCodeNO_CANDIDATE)
			return
		}
		slog.Error("service failed to reassign pull request", "error", err)
		writeJSONError(w, internalServerErrorMsg, http.StatusInternalServerError, payload.ErrCodeNOT_FOUND)
		return
	}

	writeJSONResponse(w, map[string]any{"pr": pr, "replaced_by": newReviewerID}, http.StatusOK)
}

func (h *Handler) GetUserAssignments(w http.ResponseWriter, r *http.Request) {

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
