package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/user/model"
	userMocks "github.com/mdcantarini/twitter-clone/internal/user/repository/mocks"
)

func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeRepo := &userMocks.FakeSqlRepository{
		InsertUserFunc: func(user *model.User) (*model.User, error) {
			require.Equal(t, "johndoe", user.Username)
			require.Equal(t, "John Doe", user.DisplayName)
			user.ID = 1
			return user, nil
		},
	}

	service := &Service{
		db: fakeRepo,
	}

	router := gin.New()
	v1 := router.Group("/api/v1")
	service.RegisterRoutes(v1)

	body := map[string]string{
		"username":     "johndoe",
		"display_name": "John Doe",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, uint(1), response.ID)
	require.Equal(t, "johndoe", response.Username)
	require.Equal(t, "John Doe", response.DisplayName)
}

func TestGetUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expectedUser := &model.User{
		ID:          1,
		Username:    "johndoe",
		DisplayName: "John Doe",
	}

	fakeRepo := &userMocks.FakeSqlRepository{
		GetUserFunc: func(id uint) (*model.User, error) {
			require.Equal(t, uint(1), id)
			return expectedUser, nil
		},
	}

	service := &Service{
		db: fakeRepo,
	}

	router := gin.New()
	v1 := router.Group("/api/v1")
	service.RegisterRoutes(v1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, uint(1), response.ID)
	require.Equal(t, "johndoe", response.Username)
	require.Equal(t, "John Doe", response.DisplayName)
}

func TestGetUser_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeRepo := &userMocks.FakeSqlRepository{
		GetUserFunc: func(id uint) (*model.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := &Service{
		db: fakeRepo,
	}

	router := gin.New()
	v1 := router.Group("/api/v1")
	service.RegisterRoutes(v1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "User not found", response["error"])
}
