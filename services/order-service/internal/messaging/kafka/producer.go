// services/order-service/internal/messaging/kafka/producer.go
package kafka

import (
	"context"
	"encoding/json"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/adesm/backend-roadmap/services/order-service/internal/domain"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(brokerAddr, topic string) *Producer {
	return &Producer{
		writer: &kafkago.Writer{
			Addr:     kafkago.TCP(brokerAddr),
			Topic:    topic,
			Balancer: &kafkago.Hash{},
		},
	}
}

func (p *Producer) PublishOrderCreated(event domain.OrderCreatedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(context.Background(), kafkago.Message{
		Key:   []byte(event.OrderID),
		Value: body,
	})
}
