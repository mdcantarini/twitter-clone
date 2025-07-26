package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/tweet"
	"github.com/mdcantarini/twitter-clone/internal/user"
)

func main() {
	db, err := gorm.Open(sqlite.Open("/data/twitter.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&tweet.Tweet{}, &user.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	router := gin.Default()
	api := router.Group("/api/v1")

	tweetHandler := tweet.NewHandler(db)
	tweetHandler.RegisterRoutes(api)

	log.Println("tweet-api running on :8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
