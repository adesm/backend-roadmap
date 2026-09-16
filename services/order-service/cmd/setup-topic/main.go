// services/order-service/cmd/setup-topic/main.go
package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

func main() {
	// tambahkan sementara di setup-topic/main.go, atau bikin script kecil terpisah
	conn, _ := kafkago.Dial("tcp", "localhost:9092")
	defer conn.Close()
	partitions, err := conn.ReadPartitions("order-events")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("jumlah partition: %d", len(partitions))
	for _, p := range partitions {
		log.Printf("partition=%d leader=%v", p.ID, p.Leader)
	}

	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("failed to get controller: %v", err)
	}

	controllerConn, err := kafkago.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatalf("failed to dial controller: %v", err)
	}
	defer controllerConn.Close()

	topicConfig := kafkago.TopicConfig{
		Topic:             "order-events",
		NumPartitions:     3,
		ReplicationFactor: 1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn.DeleteTopics("order-events") // abaikan error kalau topic belum ada

	err = controllerConn.CreateTopics(topicConfig)
	if err != nil {
		log.Fatalf("failed to create topic: %v", err)
	}
	_ = ctx

	log.Println("topic 'order-events' created successfully (or already exists)")
}
