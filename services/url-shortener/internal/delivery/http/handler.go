// services/url-shortener/internal/delivery/http/handler.go
package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/adesm/backend-roadmap/services/url-shortener/internal/domain"
)

// Interface milik delivery — consumer yang butuh, consumer yang definisikan.
// Cuma method yang beneran dipakai HTTP handler ini, gak wajib sama persis usecase.
type ShortenerService interface {
	Shorten(originalURL string) (*domain.ShortURL, error)
	Resolve(code string) (*domain.ShortURL, error)
}

type Handler struct {
	service ShortenerService
}

func NewHandler(service ShortenerService) *Handler {
	return &Handler{service: service}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortCode string `json:"short_code"`
	Original  string `json:"original_url"`
}

func (h *Handler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.Shorten(req.URL)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidURL) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := shortenResponse{
		ShortCode: result.ShortCode,
		Original:  result.OriginalURL,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(body); err != nil {
		return
	}
}

func (h *Handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	result, err := h.service.Resolve(code)
	if err != nil {
		if errors.Is(err, domain.ErrShortURLNotFound) {
			http.Error(w, "short url not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, result.OriginalURL, http.StatusMovedPermanently)
}
