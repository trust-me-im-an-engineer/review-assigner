package handlers

import (
	"net/http"

	"review-assigner/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddTeam(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) SetUserActivity(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) MergePullRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ReassignPullRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetUserAssignments(w http.ResponseWriter, r *http.Request) {

}
