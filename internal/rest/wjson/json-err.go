// Package wjson writes json errors or responses.
package wjson

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"review-assigner/internal/rest/payload"
)

// Error writes error to response writer in format described in openapi.
func Error(w http.ResponseWriter, msg string, statusCode int, apiCode payload.ErrorCode) {
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

// Response writes response in json format.
func Response(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to write JSON response", "error", err)
	}
}
