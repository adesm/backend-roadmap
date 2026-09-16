package domain

type ShortURLCreatedEvent struct {
	EventID     string
	ShortCode   string
	OriginalURL string
}

type EventPublisher interface {
	PublishShortURLCreated(event ShortURLCreatedEvent) error
}
