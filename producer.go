package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

func kafkaWriter(kafkaURL, topic string) *kafka.Writer {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

func producer() {
	kafkaURL := os.Getenv("KAFKA_URL")
	topic := os.Getenv("KAFKA_TOPIC")

	writer := kafkaWriter(kafkaURL, topic)
	defer writer.Close()

	fmt.Println("PRODUCE")
	i := 0
	for {
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("key: %d", i)),
			Value: []byte(fmt.Sprint(uuid.New())),
		}
		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
    i++
	}
}
