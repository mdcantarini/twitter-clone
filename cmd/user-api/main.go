package main

import (
	"github.com/mdcantarini/twitter-clone/internal/user/model"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/user"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("failed to get DB_PATH env value")
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	router := gin.Default()
	api := router.Group("/api/v1")

	userService := user.NewService(db)
	userService.RegisterRoutes(api)

	log.Println("user-api running on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
