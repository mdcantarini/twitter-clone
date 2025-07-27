package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mdcantarini/twitter-clone/internal/db/cassandra"
	"github.com/mdcantarini/twitter-clone/internal/feed"
	"log"
)

func main() {
	session := cassandra.NewSession([]string{"cassandra:9042"}, "twitter")
	defer session.Close()

	router := gin.Default()
	api := router.Group("/api/v1")

	feedHandler := feed.NewHandler(session)
	feedHandler.RegisterRoutes(api)

	log.Println("feed-api running on :8084")
	if err := router.Run(":8084"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
