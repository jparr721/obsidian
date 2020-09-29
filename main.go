package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/joho/godotenv"
	kafka "github.com/segmentio/kafka-go"
)

func initializeInterruptHandler() {
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		<-signalChannel
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func configureTopics() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	topic := os.Getenv("KAFKA_TOPIC")
	kafkaURL := os.Getenv("KAFKA_URL")

	conn, err := (&kafka.Dialer{
		Resolver: &net.Resolver{},
	}).DialLeader(ctx, "tcp", kafkaURL, topic, 0)

	if err != nil {
		log.Fatal("Failed to connect to kafka cluster")
	}
	defer conn.Close()

	conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	fmt.Println("Topics Configured")
}

func initiateProducersAndConsumers() {
	go consumer()
	go producer()
}

func main() {
	var routineWaitGroup sync.WaitGroup
	routineWaitGroup.Add(1)

	initializeInterruptHandler()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configureTopics()
	initiateProducersAndConsumers()

	routineWaitGroup.Wait()
}
