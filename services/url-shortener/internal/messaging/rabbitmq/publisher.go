package rabbitmq

import (
	"encoding/json"

	"github.com/adesm/backend-roadmap/services/url-shortener/internal/domain"
	"github.com/streadway/amqp"
)

type Publisher struct {
	channel *amqp.Channel
}

const exchangeName = "short_url_events"

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,  // durable — exchange tetap ada meski RabbitMQ restart
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{channel: ch}, nil
}

func (p *Publisher) PublishShortURLCreated(event domain.ShortURLCreatedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		exchangeName,
		"",    // routing key, diabaikan untuk fanout
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // pesan disimpan ke disk, bukan cuma memory
			Body:         body,
		},
	)
}
