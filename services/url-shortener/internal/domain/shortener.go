// services/url-shortener/internal/domain/shortener.go
package domain

import (
	"errors"
	"time"
)

// Entity: representasi "benda" nyata dalam bisnis ini
type ShortURL struct {
	ID          string
	OriginalURL string
	ShortCode   string
	CreatedAt   time.Time
}

// Sentinel errors — didefinisikan di domain, dipakai lintas layer
var (
	ErrInvalidURL       = errors.New("invalid url")
	ErrShortURLNotFound = errors.New("short url not found")
)

// Interface (kontrak) — usecase butuh ini, tapi TIDAK tau implementasinya pakai apa
type ShortURLRepository interface {
	Save(url *ShortURL) error
	FindByCode(code string) (*ShortURL, error)
}
