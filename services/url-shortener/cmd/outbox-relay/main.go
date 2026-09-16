// services/url-shortener/cmd/outbox-relay/main.go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	godotenv.Load()

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	amqpConn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("failed to connect rabbitmq: %v", err)
	}
	defer amqpConn.Close()

	ch, err := amqpConn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}

	err = ch.ExchangeDeclare("short_url_events", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}

	log.Println("outbox relay started, polling every 2s...")

	for {
		relayPendingEvents(pool, ch)
		time.Sleep(2 * time.Second)
	}
}

func relayPendingEvents(pool *pgxpool.Pool, ch *amqp.Channel) {
	ctx := context.Background()

	rows, err := pool.Query(ctx,
		`SELECT id, payload FROM outbox_events WHERE published_at IS NULL ORDER BY created_at LIMIT 100`,
	)
	if err != nil {
		log.Printf("failed to query outbox: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var payload []byte
		if err := rows.Scan(&id, &payload); err != nil {
			log.Printf("failed to scan row: %v", err)
			continue
		}

		err = ch.Publish("short_url_events", "", false, false, amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         payload,
		})
		if err != nil {
			log.Printf("failed to publish event %s: %v", id, err)
			continue // biarkan published_at tetap NULL, akan dicoba lagi loop berikutnya
		}

		_, err = pool.Exec(ctx,
			`UPDATE outbox_events SET published_at = NOW() WHERE id = $1`, id,
		)
		if err != nil {
			log.Printf("failed to mark event %s as published: %v", id, err)
		}
	}
}
