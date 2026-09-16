// services/email-service/cmd/consumer/main.go
package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ShortURLCreatedEvent struct {
	EventID     string `json:"EventID"`
	ShortCode   string `json:"ShortCode"`
	OriginalURL string `json:"OriginalURL"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare exchange yang SAMA persis dengan producer (idempotent, aman dipanggil berkali-kali)
	err = ch.ExchangeDeclare("short_url_events", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}

	// Consumer declare QUEUE-nya sendiri, lalu BIND ke exchange
	q, err := ch.QueueDeclare(
		"email_service_queue", // nama queue spesifik untuk consumer ini
		true,                  // durable
		false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	err = ch.QueueBind(q.Name, "", "short_url_events", false, nil)
	if err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer tag
		false, // auto-ack — SENGAJA false, dibahas di bawah
		false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("failed to register consumer: %v", err)
	}

	log.Println("email-service waiting for events...")

	for msg := range msgs {
		var event ShortURLCreatedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("failed to unmarshal event: %v", err)
			msg.Nack(false, false) // buang pesan yang rusak, jangan requeue
			continue
		}

		log.Printf("sending email notification for short_code=%s original_url=%s",
			event.ShortCode, event.OriginalURL)

		// simulasi kirim email (di sini nanti bisa tambah idempotency check)

		msg.Ack(false) // konfirmasi pesan udah selesai diproses
	}
}
