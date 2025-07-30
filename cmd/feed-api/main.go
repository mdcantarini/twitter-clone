package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"

	"github.com/mdcantarini/twitter-clone/internal/db/cassandra"
	"github.com/mdcantarini/twitter-clone/internal/feed"
)

func main() {
	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	if keyspace == "" {
		log.Fatal("failed to get CASSANDRA_KEYSPACE env value")
	}
	session := cassandra.NewSession([]string{"cassandra:9042"}, keyspace)
	defer session.Close()

	// Get Kafka brokers from environment or use default
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("failed to get KAFKA_BROKERS env value")
	}
	brokers := strings.Split(kafkaBrokers, ",")

	// Initialize Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       "tweets",
		GroupID:     "feed-service-v4",
		StartOffset: kafka.FirstOffset,
	})
	// Note: We don't close the reader here since it's used by the consumer goroutine

	router := gin.Default()
	api := router.Group("/api/v1")

	feedService := feed.NewService(session, reader)
	feedService.RegisterRoutes(api)

	// Start the Kafka consumer in a goroutine
	go feedService.RunTweetQueueConsumer()

	log.Println("feed-api running on :8084")
	if err := router.Run(":8084"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
