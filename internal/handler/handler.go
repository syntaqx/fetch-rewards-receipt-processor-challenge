package handler

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid/v5"
	"github.com/syntaqx/fetch-rewards-receipt-processor-challenge/internal/model"
	"github.com/syntaqx/fetch-rewards-receipt-processor-challenge/internal/repository"
	"github.com/syntaqx/fetch-rewards-receipt-processor-challenge/pkg/points"
)

type Handler struct {
	validate *validator.Validate
	repo     repository.ReceiptRepository
}

func NewHandler(validate *validator.Validate, repo repository.ReceiptRepository) *Handler {
	// Register the custom validation function
	if err := validate.RegisterValidation("price", validatePrice); err != nil {
		panic(err)
	}
	return &Handler{validate: validate, repo: repo}
}

// Custom validation function for price format
func validatePrice(fl validator.FieldLevel) bool {
	priceRegex := `^\d+\.\d{2}$`
	return regexp.MustCompile(priceRegex).MatchString(fl.Field().String())
}

// ProcessReceiptResponse represents the response for processing a receipt
type ProcessReceiptResponse struct {
	ID uuid.UUID `json:"id"`
}

func (resp *ProcessReceiptResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// ProcessReceipt godoc
// @Summary Submit a receipt for processing
// @Description Submits a receipt for processing
// @Tags receipts
// @Accept  json
// @Produce  json
// @Param receipt body model.Receipt true "Receipt"
// @Success 200 {object} ProcessReceiptResponse
// @Failure 400 {string} string "Invalid receipt"
// @Failure 500 {string} string "Failed to generate ID"
// @Router /receipts/process [post]
func (h *Handler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt model.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid receipt", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
		return
	}
	points := points.CalculatePoints(receipt)

	h.repo.SaveReceipt(id.String(), receipt, points)

	response := &ProcessReceiptResponse{ID: id}
	if err := render.Render(w, r, response); err != nil {
		http.Error(w, "Failed to render response", http.StatusInternalServerError)
	}
}

// GetPointsResponse represents the response for getting receipt points
type GetPointsResponse struct {
	Points int64 `json:"points"`
}

func (resp *GetPointsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// GetReceiptPoints godoc
// @Summary Get the points awarded for a receipt
// @Description Returns the points awarded for the receipt
// @Tags receipts
// @Produce  json
// @Param id path string true "Receipt ID"
// @Success 200 {object} GetPointsResponse
// @Failure 404 {string} string "No receipt found for that id"
// @Router /receipts/{id}/points [get]
func (h *Handler) GetReceiptPoints(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	points, exists := h.repo.GetPoints(id)
	if !exists {
		http.Error(w, "No receipt found for that id", http.StatusNotFound)
		return
	}

	response := &GetPointsResponse{Points: points}
	if err := render.Render(w, r, response); err != nil {
		http.Error(w, "Failed to render response", http.StatusInternalServerError)
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns 200 if the server is running
// @Tags health
// @Produce  json
// @Success 200 {string} string "OK"
// @Router /healthz [get]
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
