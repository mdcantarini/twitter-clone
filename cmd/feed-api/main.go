package main

import (
	"github.com/mdcantarini/twitter-clone/db/cassandra"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"

	"github.com/mdcantarini/twitter-clone/internal/feed"
)

func main() {
	// get Cassandra configuration from environment
	cassandraHost := os.Getenv("CASSANDRA_HOST")
	if cassandraHost == "" {
		log.Fatal("failed to get CASSANDRA_HOST env value")
	}
	cassandraPort := os.Getenv("CASSANDRA_PORT")
	if cassandraPort == "" {
		log.Fatal("failed to get CASSANDRA_PORT env value")
	}
	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	if keyspace == "" {
		log.Fatal("failed to get CASSANDRA_KEYSPACE env value")
	}

	hosts := []string{cassandraHost + ":" + cassandraPort}
	if cassandraNodes := os.Getenv("CASSANDRA_NODES"); cassandraNodes != "" {
		hosts = strings.Split(cassandraNodes, ",")
	}

	session := cassandra.NewSession(hosts, keyspace)
	defer session.Close()

	// Get Kafka brokers from environment or use default
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("failed to get KAFKA_BROKERS env value")
	}
	brokers := strings.Split(kafkaBrokers, ",")

	// Initialize Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   "tweets",
		GroupID: "feed-service-v1",
	})

	router := gin.Default()
	api := router.Group("/api/v1")

	feedService := feed.NewService(session, reader)
	feedService.RegisterRoutes(api)

	// Start the Kafka consumer in a goroutine
	go feedService.RunTweetQueueConsumer()

	log.Println("feed-api running on :8084")
	if err := router.Run(":8084"); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
