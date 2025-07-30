package feed

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	followMocks "github.com/mdcantarini/twitter-clone/internal/client/follow/mocks"
	"github.com/mdcantarini/twitter-clone/internal/client/tweet"
	tweetMocks "github.com/mdcantarini/twitter-clone/internal/client/tweet/mocks"
	"github.com/mdcantarini/twitter-clone/internal/feed/model"
	feedMocks "github.com/mdcantarini/twitter-clone/internal/feed/repository/mocks"
)

func TestGetUserFeed_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeRepo := &feedMocks.FakeNoSqlRepository{
		GetUserTimelineFunc: func(userID uint, limit int) ([]model.FeedEntry, error) {
			return []model.FeedEntry{
				{
					TweetID:   gocql.UUID(uuid.New()),
					AuthorID:  2,
					Content:   "Hello world",
					CreatedAt: time.Now(),
				},
				{
					TweetID:   gocql.UUID(uuid.New()),
					AuthorID:  3,
					Content:   "Another tweet",
					CreatedAt: time.Now().Add(-1 * time.Hour),
				},
			}, nil
		},
	}

	service := &Service{
		db: fakeRepo,
	}

	router := gin.New()
	v1 := router.Group("/api/v1")
	service.RegisterRoutes(v1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/feed/1", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response []model.FeedEntry
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Len(t, response, 2)
}

func TestProcessTweetEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFollowClient := followMocks.NewMockClient(ctrl)
	mockTweetClient := tweetMocks.NewMockClient(ctrl)

	tweetID := gocql.UUID(uuid.New())
	userID := uint(1)
	
	tweetEvent := TweetEvent{
		TweetID: tweetID,
		UserID:  userID,
	}

	mockTweetClient.EXPECT().FetchTweet(tweetID.String()).Return(&tweet.Tweet{
		TweetID:   tweetID.String(),
		UserID:    userID,
		Content:   "Test tweet content",
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil)

	followerIDs := []uint{2, 3, 4}
	mockFollowClient.EXPECT().FetchFollowerIds(userID).Return(followerIDs, nil)

	fakeRepo := &feedMocks.FakeNoSqlRepository{
		InsertUserTimelineFunc: func(followerIds []uint, createdAt string, tweetId string, userId uint, tweetContent string) error {
			require.Len(t, followerIds, 3)
			require.Equal(t, userID, userId)
			require.Equal(t, "Test tweet content", tweetContent)
			return nil
		},
	}

	service := &Service{
		db:           fakeRepo,
		followClient: mockFollowClient,
		tweetClient:  mockTweetClient,
	}

	err := service.processTweetEvent(tweetEvent)
	require.NoError(t, err)
}
