package rest

import (
	"net/http"

	"review-assigner/internal/rest/handlers"
	"review-assigner/internal/rest/middleware"
	"review-assigner/internal/service"
	"review-assigner/internal/validator"
)

func NewRouter(s *service.Service, v *validator.Validator, a *middleware.Auth) *http.ServeMux {
	h := handlers.NewHandler(s, v)
	mux := http.NewServeMux()

	// No auth
	mux.HandleFunc("POST /team/add", h.AddTeamAddUpdateUsers)

	// User or admin
	mux.HandleFunc("GET /team/get", a.UserOrAdmin(h.GetTeam))
	mux.HandleFunc("GET /users/getReview", a.UserOrAdmin(h.GetUserAssignments))

	// Admin
	mux.HandleFunc("POST /users/setIsActive", a.Admin(h.SetUserActivity))
	mux.HandleFunc("POST /pullRequest/create", a.Admin(h.CreatePullRequest))
	mux.HandleFunc("POST /pullRequest/merge", a.Admin(h.MergePullRequest))
	mux.HandleFunc("POST /pullRequest/reassign", a.Admin(h.ReassignPullRequest))

	return mux
}
