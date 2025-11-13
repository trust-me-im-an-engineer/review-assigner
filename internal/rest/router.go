package rest

import (
	"net/http"

	"review-assigner/internal/rest/handlers"
)

func NewRouter(s *service.Service) *http.ServeMux {
	h := handlers.NewHandler(s)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /team/add", h.AddTeamAddUpdateUsers)
	mux.HandleFunc("GET /team/get", h.GetTeam)
	mux.HandleFunc("POST /users/setIsActive", h.SetUserActivity)
	mux.HandleFunc("POST /pullRequest/create", h.CreatePullRequest)
	mux.HandleFunc("POST /pullRequest/merge", h.MergePullRequest)
	mux.HandleFunc("POST /pullRequest/reassign", h.ReassignPullRequest)
	mux.HandleFunc("GET /users/getReview", h.GetUserAssignments)

	return mux
}
