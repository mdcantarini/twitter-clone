package tweet

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"

	"github.com/mdcantarini/twitter-clone/internal/tweet/model"
	tweetMocks "github.com/mdcantarini/twitter-clone/internal/tweet/repository/mocks"
)

func TestCreateTweet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeRepo := &tweetMocks.FakeNoSqlRepository{
		InsertTweetFunc: func(tweet model.Tweet) error {
			require.Equal(t, uint(1), tweet.UserID)
			require.Equal(t, "Hello, Twitter!", tweet.Content)
			return nil
		},
	}

	service := &Service{
		db:                    fakeRepo,
		tweetsMessageProducer: &mockKafkaWriter{},
	}

	router := gin.New()
	v1 := router.Group("/api/v1")
	service.RegisterRoutes(v1)

	body := CreateTweetRequest{
		UserID:  1,
		Content: "Hello, Twitter!",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tweets", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, uint(1), response.UserID)
	require.Equal(t, "Hello, Twitter!", response.Content)
}

func TestCreateTweet_Failure_TweetLongerThanAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeRepo := &tweetMocks.FakeNoSqlRepository{}

	service := &Service{
		db:                    fakeRepo,
		tweetsMessageProducer: &mockKafkaWriter{},
	}

	router := gin.New()
	v1 := router.Group("/api/v1")
	service.RegisterRoutes(v1)

	body := CreateTweetRequest{
		UserID:  1,
		Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit sollicitudin, ultricies nullam euismod. Neque blandit egestas quis dignissim ullamcorper nam suscipit id nisl pellentesque, suspendisse pulvinar sociosqu diam malesuada nibh faucibus mattis habitant, rutrum nulla maecenas alique.",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tweets", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "tweet content cannot be longer than 280 characters", response["error"])
}

func TestGetTweetById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tweetID := gocql.UUID(uuid.New())
	expectedTweet := model.Tweet{
		TweetID:   tweetID,
		UserID:    1,
		Content:   "Test tweet",
		CreatedAt: time.Now(),
	}

	fakeRepo := &tweetMocks.FakeNoSqlRepository{
		GetTweetByIdFunc: func(id gocql.UUID) (model.Tweet, error) {
			require.Equal(t, tweetID, id)
			return expectedTweet, nil
		},
	}

	service := &Service{
		db: fakeRepo,
	}

	router := gin.New()
	v1 := router.Group("/api/v1")
	service.RegisterRoutes(v1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tweet/"+tweetID.String(), nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, tweetID, response.TweetID)
	require.Equal(t, uint(1), response.UserID)
	require.Equal(t, "Test tweet", response.Content)
}

type mockKafkaWriter struct{}

func (m *mockKafkaWriter) WriteMessages(_ context.Context, _ ...kafka.Message) error {
	return nil
}
