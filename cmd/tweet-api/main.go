package main

import (
	"github.com/mdcantarini/twitter-clone/db/cassandra"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mdcantarini/twitter-clone/internal/tweet"
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

	// Initialize Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   "tweets",
	})

	router := gin.Default()
	api := router.Group("/api/v1")

	tweetService := tweet.NewService(session, writer)
	tweetService.RegisterRoutes(api)

	log.Println("tweet-api running on :8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
