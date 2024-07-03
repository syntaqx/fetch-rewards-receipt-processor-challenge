package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// RegisterRoutes registers the health handler routes.
func (h *HealthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/healthz", h.HealthCheck)
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns 200 if the server is running
// @Tags health
// @Produce  json
// @Success 200 {string} string "OK"
// @Router /healthz [get]
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
