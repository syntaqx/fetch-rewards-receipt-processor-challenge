package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/syntaqx/fetch-rewards-receipt-processor-challenge/docs"
	"github.com/syntaqx/fetch-rewards-receipt-processor-challenge/internal/handler"
	"github.com/syntaqx/fetch-rewards-receipt-processor-challenge/internal/repository"
)

// @title Receipt Processor API
// @version 1.0
// @description This is a simple receipt processor API.
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
func main() {
	r := chi.NewRouter()

	// Adding important chi base middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.AllowAll().Handler)

	// Initialize validator and repository
	validate := validator.New()
	repo := repository.NewReceiptRepository()

	// Initialize handlers with dependencies
	receiptHandler := handler.NewReceiptHandler(validate, repo)
	healthHandler := handler.NewHealthHandler()

	// Register routes
	receiptHandler.RegisterRoutes(r)
	healthHandler.RegisterRoutes(r)

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: r,
	}

	fmt.Printf("server listening at %s\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			fmt.Printf("server closed unexpectedly with error: %v\n", err)
		}
	}
}
