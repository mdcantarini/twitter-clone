package main

import (
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mdcantarini/twitter-clone/internal/db/cassandra"

	"github.com/mdcantarini/twitter-clone/internal/tweet"
)

func main() {
	// Get Cassandra configuration from environment
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

	router := gin.Default()
	api := router.Group("/api/v1")

	// instantiate kafka producer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:29092"},
		Topic:   "tweets",
	})

	tweetService := tweet.NewService(session, writer)
	tweetService.RegisterRoutes(api)

	log.Println("tweet-api running on :8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
