// services/url-shortener/cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	deliveryhttp "github.com/adesm/backend-roadmap/services/url-shortener/internal/delivery/http"
	"github.com/adesm/backend-roadmap/services/url-shortener/internal/repository/postgres"
	"github.com/adesm/backend-roadmap/services/url-shortener/internal/usecase"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// --- Dependency wiring, dari dalam ke luar ---
	repo := postgres.NewShortenerRepository(pool)
	uc := usecase.NewShortenerUsecase(repo)
	handler := deliveryhttp.NewHandler(uc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", handler.HandleShorten)
	mux.HandleFunc("GET /{code}", handler.HandleRedirect)

	log.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
