package follow

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/mdcantarini/twitter-clone/internal/follow/model"
	followMocks "github.com/mdcantarini/twitter-clone/internal/follow/repository/mocks"
)

func TestFollowUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeRepo := &followMocks.FakeSqlRepository{
		InsertFollowFunc: func(follow *model.Follow) error {
			require.Equal(t, uint(1), follow.FollowerID)
			require.Equal(t, uint(2), follow.FollowedID)
			return nil
		},
	}

	service := &Service{
		db: fakeRepo,
	}

	router := gin.New()
	v1 := router.Group("/api/v1")
	service.RegisterRoutes(v1)

	body := map[string]uint{
		"follower_id": 1,
		"followed_id": 2,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/follow", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "Followed successfully", response["message"])
}

func TestGetFollowerIds_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expectedFollowers := []model.Follow{
		{FollowerID: 2, FollowedID: 1},
		{FollowerID: 3, FollowedID: 1},
		{FollowerID: 4, FollowedID: 1},
	}

	fakeRepo := &followMocks.FakeSqlRepository{
		GetFollowersFunc: func(followedID uint) ([]model.Follow, error) {
			require.Equal(t, uint(1), followedID)
			return expectedFollowers, nil
		},
	}

	service := &Service{
		db: fakeRepo,
	}

	router := gin.New()
	v1 := router.Group("/api/v1")
	service.RegisterRoutes(v1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/1/follower_ids", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response map[string][]uint
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	followerIDs := response["follower_ids"]
	require.Len(t, followerIDs, 3)
	require.Equal(t, []uint{2, 3, 4}, followerIDs)
}
