// services/order-service/cmd/api/main.go
package main

import (
	"fmt"
	"log"

	"github.com/adesm/backend-roadmap/services/order-service/internal/domain"
	"github.com/adesm/backend-roadmap/services/order-service/internal/messaging/kafka"
)

func main() {
	producer := kafka.NewProducer("localhost:9092", "order-events")

	// 5 event untuk order_id SAMA — harus tetap berurutan
	for i := 1; i <= 5; i++ {
		event := domain.OrderCreatedEvent{
			OrderID: "ORD-SAME",
			UserID:  "user-1",
			Amount:  float64(i * 1000),
		}
		if err := producer.PublishOrderCreated(event); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("published event %d for ORD-SAME\n", i)
	}

	// 5 event untuk order_id BERBEDA-BEDA
	for i := 1; i <= 5; i++ {
		event := domain.OrderCreatedEvent{
			OrderID: fmt.Sprintf("ORD-%d", i),
			UserID:  "user-1",
			Amount:  float64(i * 500),
		}
		if err := producer.PublishOrderCreated(event); err != nil {
			log.Fatal(err)
		}
	}
}
