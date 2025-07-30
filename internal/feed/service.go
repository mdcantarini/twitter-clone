package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mdcantarini/twitter-clone/internal/feed/model"
	messagebroker "github.com/mdcantarini/twitter-clone/messagebroker/kafka"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/segmentio/kafka-go"

	"github.com/mdcantarini/twitter-clone/internal/client/follow"
	"github.com/mdcantarini/twitter-clone/internal/client/tweet"
	"github.com/mdcantarini/twitter-clone/internal/feed/repository"
)

type Service struct {
	db                  repository.Repository
	followClient        follow.Client
	tweetClient         tweet.Client
	tweetsMessageReader messagebroker.Consumer
}

func NewService(session *gocql.Session, reader *kafka.Reader) *Service {
	noSqlImpl := repository.NewNoSqlRepositoryImplementation(session)
	followClient := follow.NewFollowClient()
	tweetClient := tweet.NewTweetClient()

	return &Service{
		db:                  noSqlImpl,
		followClient:        followClient,
		tweetClient:         tweetClient,
		tweetsMessageReader: reader,
	}
}

const userFeedLimit = 50

// TODO - Improve! There is a limitation here in the number of entries we can return for a give user.
// We should use pagination to avoid overloaded database and return as many entries as the user request.
func (s *Service) GetUserFeed(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	repoFeedItems, err := s.db.GetUserTimeline(uint(id), userFeedLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feed"})
		return
	}

	// Convert repository.FeedEntry to feed.FeedEntry
	feedItems := make([]model.FeedEntry, len(repoFeedItems))
	for i, item := range repoFeedItems {
		feedItems[i] = model.FeedEntry{
			TweetID:   item.TweetID,
			AuthorID:  item.AuthorID,
			Content:   item.Content,
			CreatedAt: item.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, feedItems)
}

func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/feed/:user_id", s.GetUserFeed)
}

// ================================================
// Tweet events processor
// ================================================

type TweetEvent struct {
	TweetID gocql.UUID `json:"tweet_id"`
	UserID  uint       `json:"user_id"`
}

// TODO - Improve! If there is a failure when proccesing the message we can lose it.
func (s *Service) RunTweetQueueConsumer() {
	log.Println("Starting Kafka consumer for tweets topic...")
	for {
		msg, err := s.tweetsMessageReader.ReadMessage(context.Background())
		if err != nil {
			log.Println("kafka read error:", err)
			continue
		}

		var tweetEvent TweetEvent
		if err := json.Unmarshal(msg.Value, &tweetEvent); err != nil {
			log.Println("failed to unmarshal tweet:", err)
			continue
		}

		err = s.processTweetEvent(tweetEvent)
		if err != nil {
			log.Println("failed to process tweet event:", err)
			continue
		}
	}
}

func (s *Service) processTweetEvent(event TweetEvent) error {
	// 1. get tweet from tweet-api
	tweet, err := s.tweetClient.FetchTweet(event.TweetID.String())
	if err != nil {
		return fmt.Errorf("failed to get tweet: %s", err)
	}

	// 2. get followers from follow-api
	followerIds, err := s.followClient.FetchFollowerIds(tweet.UserID)
	if err != nil {
		return fmt.Errorf("failed to get followers: %s", err)
	}

	// 3. create batch for updating user timeline
	err = s.db.InsertUserTimeline(followerIds, tweet.CreatedAt, tweet.TweetID, tweet.UserID, tweet.Content)
	if err != nil {
		return fmt.Errorf("failed updating user timeline: %s", err)
	}

	return nil
}
