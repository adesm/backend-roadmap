// services/url-shortener/internal/repository/postgres/shortener_test.go
package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/adesm/backend-roadmap/services/url-shortener/internal/domain"
	"github.com/adesm/backend-roadmap/services/url-shortener/internal/repository/postgres"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	t.Cleanup(pool.Close)

	_, err = pool.Exec(ctx, `
		CREATE TABLE short_urls (
			id VARCHAR(32) PRIMARY KEY,
			original_url TEXT NOT NULL,
			short_code VARCHAR(16) UNIQUE NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return pool
}

func TestShortenerRepository_SaveAndFind(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewShortenerRepository(pool)

	url := &domain.ShortURL{
		ID:          "test-id-1",
		OriginalURL: "https://example.com",
		ShortCode:   "abc123",
		CreatedAt:   time.Now().Truncate(time.Second),
	}

	err := repo.Save(url)
	if err != nil {
		t.Fatalf("expected no error saving, got %v", err)
	}

	found, err := repo.FindByCode("abc123")
	if err != nil {
		t.Fatalf("expected no error finding, got %v", err)
	}
	if found.OriginalURL != url.OriginalURL {
		t.Errorf("expected OriginalURL %q, got %q", url.OriginalURL, found.OriginalURL)
	}
}

func TestShortenerRepository_FindByCode_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewShortenerRepository(pool)

	_, err := repo.FindByCode("doesnotexist")

	if err != domain.ErrShortURLNotFound {
		t.Errorf("expected ErrShortURLNotFound, got %v", err)
	}
}
