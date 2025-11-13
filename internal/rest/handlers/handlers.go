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

	writeJSONResponse(w, map[string]interface{}{"team": team}, http.StatusCreated)
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

	writeJSONResponse(w, map[string]interface{}{"user": user}, http.StatusOK)
}
func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) MergePullRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ReassignPullRequest(w http.ResponseWriter, r *http.Request) {

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
