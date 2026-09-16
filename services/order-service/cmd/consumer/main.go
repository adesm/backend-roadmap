// services/order-service/cmd/consumer/main.go
package main

import (
	"context"
	"encoding/json"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

type OrderCreatedEvent struct {
	OrderID string  `json:"order_id"`
	UserID  string  `json:"user_id"`
	Amount  float64 `json:"amount"`
}

func main() {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "order-events",
		GroupID: "order-processor",
	})
	defer reader.Close()

	log.Println("consumer started, waiting for messages...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}

		var event OrderCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			continue
		}

		log.Printf("partition=%d offset=%d order_id=%s amount=%.2f",
			msg.Partition, msg.Offset, event.OrderID, event.Amount)
	}
}
