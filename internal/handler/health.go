package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// RegisterRoutes registers the health handler routes.
func (h *HealthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/healthz", h.HealthCheck)
}

// HealthCheckResponse represents the response for the health check endpoint.
type HealthCheckResponse struct {
	Status string `json:"status"`
}

func (resp *HealthCheckResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
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
	render.Render(w, r, &HealthCheckResponse{Status: "OK"})
}
