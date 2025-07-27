package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/follow"
	"github.com/mdcantarini/twitter-clone/internal/user"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "twitter.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&follow.Follow{}, &user.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	router := gin.Default()
	api := router.Group("/api/v1")

	followHandler := follow.NewHandler(db)
	followHandler.RegisterRoutes(api)

	log.Println("follow-api running on :8083")
	if err := router.Run(":8083"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
