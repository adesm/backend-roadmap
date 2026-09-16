package usecase

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/adesm/backend-roadmap/services/url-shortener/internal/domain"
)

const shortCodeLength = 6
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ShortenerUsecase struct {
	repo domain.ShortURLRepository
}

// Constructor — inject dependency lewat interface, bukan concrete type
func NewShortenerUsecase(repo domain.ShortURLRepository) *ShortenerUsecase {
	return &ShortenerUsecase{repo: repo}
}

func (u *ShortenerUsecase) Shorten(originalURL string) (*domain.ShortURL, error) {
	if originalURL == "" {
		return nil, domain.ErrInvalidURL
	}

	code, err := generateShortCode(shortCodeLength)
	if err != nil {
		return nil, err
	}

	shortURL := &domain.ShortURL{
		ID:          generateID(),
		OriginalURL: originalURL,
		ShortCode:   code,
		CreatedAt:   time.Now(),
	}

	event := domain.ShortURLCreatedEvent{
		EventID:     shortURL.ID, // konsisten dengan diskusi idempotency sebelumnya
		ShortCode:   shortURL.ShortCode,
		OriginalURL: shortURL.OriginalURL,
	}

	if err := u.repo.SaveWithEvent(shortURL, event); err != nil {
		return nil, err
	}

	return shortURL, nil
}

func (u *ShortenerUsecase) Resolve(code string) (*domain.ShortURL, error) {
	return u.repo.FindByCode(code)
}

func generateShortCode(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}
	return string(result), nil
}

func generateID() string {
	code, _ := generateShortCode(16)
	return code
}
