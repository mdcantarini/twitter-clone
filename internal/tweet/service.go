package tweet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/google/uuid"

	"github.com/segmentio/kafka-go"
)

type Service struct {
	db                  *gocql.Session
	tweetsQueueProducer *kafka.Writer
}

func NewService(db *gocql.Session, tweetProducer *kafka.Writer) *Service {
	return &Service{
		db:                  db,
		tweetsQueueProducer: tweetProducer,
	}
}

type TweetEvent struct {
	TweetID gocql.UUID `json:"tweet_id"`
	UserID  uint       `json:"user_id"`
}

func (s *Service) CreateTweet(c *gin.Context) {
	var input struct {
		UserID  uint   `json:"user_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO - Limit content to 280 char

	tweetIdUUID := gocql.UUID(uuid.New())

	newTweet := Tweet{
		TweetID:   tweetIdUUID,
		UserID:    input.UserID,
		Content:   input.Content,
		CreatedAt: time.Now(),
	}

	err := InsertTweet(s.db, newTweet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tweet"})
		return
	}

	event := TweetEvent{
		TweetID: tweetIdUUID,
		UserID:  input.UserID,
	}
	payload, _ := json.Marshal(event)

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("tweet-id")),
		Value: payload,
	}

	// TODO - Find a way to ensure 100% the message post
	err = s.tweetsQueueProducer.WriteMessages(c.Request.Context(), msg)
	if err != nil {
		// TODO - Log error here!
		fmt.Println(err)
	}

	c.JSON(http.StatusCreated, newTweet)
}

const defaultLimit = 50

func (s *Service) GetTweetsByUser(c *gin.Context) {
	// Get user_id from URL parameter
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get limit from query parameter
	limitStr := c.DefaultQuery("limit", strconv.Itoa(defaultLimit))
	limit, err := strconv.ParseUint(limitStr, 10, 32)
	if err != nil {
		limit = defaultLimit
	}

	tweets, err := GetTweetsByUser(s.db, uint(userID), uint(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tweets"})
		return
	}

	c.JSON(http.StatusOK, tweets)
}

func (s *Service) GetTweetById(c *gin.Context) {
	// Get user_id from URL parameter
	tweetIdStr := c.Param("tweet_id")

	tweet, err := GetTweetById(s.db, gocql.UUID(uuid.MustParse(tweetIdStr)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tweet"})
		return
	}

	c.JSON(http.StatusOK, tweet)
}

func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/tweets", s.CreateTweet)
	router.GET("/tweet/:tweet_id", s.GetTweetById)
	router.GET("/users/:user_id/tweets", s.GetTweetsByUser)
}
