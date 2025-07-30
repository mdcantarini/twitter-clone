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

	"github.com/mdcantarini/twitter-clone/internal/client/followapi"
	"github.com/mdcantarini/twitter-clone/internal/client/tweetapi"
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

	feedItems, err := GetUserTimeline(s.db, uint(id), userFeedLimit)
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
	log.Println("Starting Kafka consumer for tweets topic...")
	for {
		msg, err := s.tweetsQueueConsumer.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		log.Printf("Received message from Kafka: %s", string(msg.Value))

		var tweetEvent TweetEvent
		if err := json.Unmarshal(msg.Value, &tweetEvent); err != nil {
			log.Println("Failed to unmarshal tweet:", err)
			continue
		}

		// 1. Get tweet from tweet-api
		log.Printf("Fetching tweet %s from tweet-api", tweetEvent.TweetID.String())
		tweet, err := tweetapi.FetchTweet(tweetEvent.TweetID.String())
		if err != nil {
			log.Println("Failed to get tweet:", err)
			continue
		}
		log.Printf("Successfully fetched tweet: %+v", tweet)

		// 2. Get followers from follow-api
		log.Printf("Fetching followers for user %d from follow-api", tweet.UserID)
		followerIds, err := followapi.FetchFollowerIds(tweet.UserID)
		if err != nil {
			log.Println("Failed to get followers:", err)
			continue
		}
		log.Printf("Found %d followers", len(followerIds))

		// 3. Create batch for updating user timeline - FanOut Write pattern
		err = InsertUserTimeline(s.db, followerIds, tweet.CreatedAt, tweet.TweetID, tweet.UserID, tweet.Content)
		if err != nil {
			log.Println("Failed updating user timeline:", err)
			continue
		}

		log.Printf("Successfully fan-out tweet %s to %d followers", tweetEvent.TweetID.String(), len(followerIds))
	}
}
