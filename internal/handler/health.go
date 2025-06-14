package handler

import (
	"net/http"
)

// HealthHandler handles the /healthz route.
type HealthHandler struct{}

// NewHealthHandler creates and returns a new HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// ServeHTTP implements the http.Handler interface for HealthHandler.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
