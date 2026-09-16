package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/adesm/backend-roadmap/services/url-shortener/internal/domain"
)

type ShortenerRepository struct {
	db *pgxpool.Pool
}

// Compile-time assertion — samain prinsipnya kayak di mock tadi
var _ domain.ShortURLRepository = (*ShortenerRepository)(nil)

func NewShortenerRepository(db *pgxpool.Pool) *ShortenerRepository {
	return &ShortenerRepository{db: db}
}

func (r *ShortenerRepository) Save(url *domain.ShortURL) error {
	query := `
		INSERT INTO short_urls (id, original_url, short_code, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(context.Background(), query,
		url.ID, url.OriginalURL, url.ShortCode, url.CreatedAt)
	return err
}

func (r *ShortenerRepository) SaveWithEvent(url *domain.ShortURL, event domain.ShortURLCreatedEvent) error {
	ctx := context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO short_urls (id, original_url, short_code, created_at) VALUES ($1, $2, $3, $4)`,
		url.ID, url.OriginalURL, url.ShortCode, url.CreatedAt,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	payload, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO outbox_events (id, event_type, payload) VALUES ($1, $2, $3)`,
		event.EventID, "short_url.created", payload,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return tx.Commit(ctx)
}

func (r *ShortenerRepository) FindByCode(code string) (*domain.ShortURL, error) {
	query := `
		SELECT id, original_url, short_code, created_at
		FROM short_urls
		WHERE short_code = $1
	`
	row := r.db.QueryRow(context.Background(), query, code)

	var url domain.ShortURL
	err := row.Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrShortURLNotFound
	}
	if err != nil {
		return nil, err
	}

	return &url, nil
}
