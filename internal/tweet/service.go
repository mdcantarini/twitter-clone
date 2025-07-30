package tweet

import (
	"encoding/json"
	"fmt"
	"github.com/mdcantarini/twitter-clone/internal/tweet/model"
	messagebroker "github.com/mdcantarini/twitter-clone/messagebroker/kafka"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"github.com/mdcantarini/twitter-clone/internal/tweet/repository"
)

type Service struct {
	db                    repository.Repository
	tweetsMessageProducer messagebroker.Producer
}

func NewService(session *gocql.Session, tweetProducer *kafka.Writer) *Service {
	sqlImpl := repository.NewNoSqlRepositoryImplementation(session)

	return &Service{
		db:                    sqlImpl,
		tweetsMessageProducer: tweetProducer,
	}
}

type CreateTweetRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type TweetEvent struct {
	TweetID gocql.UUID `json:"tweet_id"`
	UserID  uint       `json:"user_id"`
}

func (s *Service) CreateTweet(c *gin.Context) {
	req := CreateTweetRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateCreateTweetRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweetIdUUID := gocql.UUID(uuid.New())
	newTweet := model.Tweet{
		TweetID:   tweetIdUUID,
		UserID:    req.UserID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	err := s.db.InsertTweet(newTweet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tweet"})
		return
	}

	event := TweetEvent{
		TweetID: tweetIdUUID,
		UserID:  req.UserID,
	}
	payload, _ := json.Marshal(event)

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("tweet-id")),
		Value: payload,
	}

	// TODO - Improve! Tweet is saved but fan-out can fail.
	err = s.tweetsMessageProducer.WriteMessages(c.Request.Context(), msg)
	if err != nil {
		log.Printf("Failed to publish tweet event to Kafka: %v", err)
	}

	c.JSON(http.StatusCreated, newTweet)
}

func validateCreateTweetRequest(req CreateTweetRequest) error {
	if len(req.Content) > 280 {
		return fmt.Errorf("tweet content cannot be longer than 280 characters")
	}

	return nil
}

// TODO - Improve! Add mechanism to cache hot tweets
func (s *Service) GetTweetById(c *gin.Context) {
	// Get user_id from URL parameter
	tweetIdStr := c.Param("tweet_id")

	tweet, err := s.db.GetTweetById(gocql.UUID(uuid.MustParse(tweetIdStr)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tweet"})
		return
	}

	c.JSON(http.StatusOK, tweet)
}

func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/tweets", s.CreateTweet)
	router.GET("/tweet/:tweet_id", s.GetTweetById)
}
