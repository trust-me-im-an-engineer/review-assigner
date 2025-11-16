package middleware

import (
	"net/http"

	"review-assigner/internal/rest/payload"
	"review-assigner/internal/rest/wjson"
)

type Auth struct {
	adminToken string
}

func NewAuth(adminToken string) *Auth {
	return &Auth{adminToken: adminToken}
}

// UserOrAdmin is a middleware that allows any user token and specified admin token
func (a *Auth) UserOrAdmin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminToken := r.Header.Get("X-Admin-Token")
		userToken := r.Header.Get("X-User-Token")
		if userToken != "" || adminToken != "" && adminToken == a.adminToken {
			h(w, r)
			return
		}

		wjson.Error(w, "Authentication required (Admin or User)", http.StatusUnauthorized, payload.ErrCodeNotFound)
		return
	}
}

// Admin is a middleware that allows only specified admin token
func (a *Auth) Admin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Admin-Token")
		if token == "" || token != a.adminToken {
			wjson.Error(w, "Authentication required: Admin", http.StatusUnauthorized, payload.ErrCodeNotFound)
			return
		}
		h(w, r)
	}
}
