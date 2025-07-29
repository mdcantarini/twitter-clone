package feed

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/segmentio/kafka-go"

	"github.com/mdcantarini/twitter-clone/internal/client"
)

type Service struct {
	db                  *gocql.Session
	tweetsQueueConsumer *kafka.Reader
}

func NewService(db *gocql.Session, reader *kafka.Reader) *Service {
	return &Service{
		db:                  db,
		tweetsQueueConsumer: reader,
	}
}

const userFeedLimit = 50

func (s *Service) GetUserFeed(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	feedItems, err := GetUserFeed(s.db, uint(id), userFeedLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feed"})
		return
	}

	c.JSON(http.StatusOK, feedItems)
}

func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/feed/:user_id", s.GetUserFeed)
}

type TweetEvent struct {
	TweetID gocql.UUID `json:"tweet_id"`
	UserID  uint       `json:"user_id"`
}

func (s *Service) RunTweetQueueConsumer() {
	for {
		msg, err := s.tweetsQueueConsumer.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		var tweetEvent TweetEvent
		if err := json.Unmarshal(msg.Value, &tweetEvent); err != nil {
			log.Println("Failed to unmarshal tweet:", err)
			continue
		}

		// 1. Get tweet from tweet-api
		tweet, err := client.FetchTweet(tweetEvent.TweetID.String())
		if err != nil {
			log.Println("Failed to get tweet:", err)
			continue
		}

		// 2. Get followers from follow-api
		followers, err := client.FetchFollowers(tweet.UserID)
		if err != nil {
			log.Println("Failed to get followers:", err)
			continue
		}

		// 3. Create batch for updating user timeline -FanOut Write pattern
		batch := s.db.NewBatch(gocql.UnloggedBatch)
		for _, followerID := range followers {
			// TODO - move this function to repository
			batch.Query(`
			INSERT INTO user_timeline (user_id, created_at, tweet_id, author_id, content)
			VALUES (?, ?, ?, ?, ?)`,
				followerID, tweet.CreatedAt, tweetEvent.TweetID, tweet.UserID, tweet.Content,
			)
		}

		err = s.db.ExecuteBatch(batch)
		if err != nil {
			log.Println("Failed updating user timeline:", err)
			continue
		}
	}
}
