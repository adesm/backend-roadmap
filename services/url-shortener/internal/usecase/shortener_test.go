// services/url-shortener/internal/usecase/shortener_test.go
package usecase_test

import (
	"errors"
	"testing"

	"github.com/adesm/backend-roadmap/services/url-shortener/internal/domain"
	"github.com/adesm/backend-roadmap/services/url-shortener/internal/usecase"
)

// mockRepository — implementasi PALSU dari domain.ShortURLRepository, khusus testing
type mockRepository struct {
	data map[string]*domain.ShortURL
}

func newMockRepository() *mockRepository {
	return &mockRepository{data: make(map[string]*domain.ShortURL)}
}

func (m *mockRepository) Save(url *domain.ShortURL) error {
	m.data[url.ShortCode] = url
	return nil
}

func (m *mockRepository) SaveWithEvent(url *domain.ShortURL, urlEvent domain.ShortURLCreatedEvent) error {
	m.data[url.ShortCode] = url
	return nil
}

func (m *mockRepository) FindByCode(code string) (*domain.ShortURL, error) {
	url, ok := m.data[code]
	if !ok {
		return nil, domain.ErrShortURLNotFound
	}
	return url, nil
}

func TestShorten_Success(t *testing.T) {
	repo := newMockRepository()
	u := usecase.NewShortenerUsecase(repo)

	result, err := u.Shorten("https://example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.OriginalURL != "https://example.com" {
		t.Errorf("expected OriginalURL to be preserved, got %q", result.OriginalURL)
	}
	if len(result.ShortCode) != 6 {
		t.Errorf("expected ShortCode length 6, got %d", len(result.ShortCode))
	}
}

func TestShorten_EmptyURL_ReturnsError(t *testing.T) {
	repo := newMockRepository()
	u := usecase.NewShortenerUsecase(repo)

	_, err := u.Shorten("")

	if !errors.Is(err, domain.ErrInvalidURL) {
		t.Errorf("expected ErrInvalidURL, got %v", err)
	}
}

func TestResolve_Found(t *testing.T) {
	repo := newMockRepository()
	u := usecase.NewShortenerUsecase(repo)

	saved, _ := u.Shorten("https://example.com")

	result, err := u.Resolve(saved.ShortCode)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.OriginalURL != "https://example.com" {
		t.Errorf("expected to resolve back to original URL, got %q", result.OriginalURL)
	}
}

func TestResolve_NotFound(t *testing.T) {
	repo := newMockRepository()
	u := usecase.NewShortenerUsecase(repo)

	_, err := u.Resolve("doesnotexist")

	if !errors.Is(err, domain.ErrShortURLNotFound) {
		t.Errorf("expected ErrShortURLNotFound, got %v", err)
	}
}
