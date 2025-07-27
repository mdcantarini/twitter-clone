package tweet

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type Handler struct {
	db *gocql.Session
}

func NewHandler(db *gocql.Session) *Handler {
	return &Handler{db: db}
}

func (h *Handler) CreateTweet(c *gin.Context) {
	var input struct {
		UserID  uint   `json:"user_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweetIdUUID := gocql.UUID(uuid.New())

	newTweet := Tweet{
		TweetID:   tweetIdUUID,
		UserID:    input.UserID,
		Content:   input.Content,
		CreatedAt: time.Now(),
	}

	if err := InsertTweet(h.db, newTweet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tweet"})
		return
	}

	c.JSON(http.StatusCreated, newTweet)
}

const defaultLimit = 50

func (h *Handler) GetTweetsByUser(c *gin.Context) {
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

	tweets, err := GetTweetsByUser(h.db, uint(userID), uint(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tweets"})
		return
	}

	c.JSON(http.StatusOK, tweets)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/tweets", h.CreateTweet)
	router.GET("/users/:user_id/tweets", h.GetTweetsByUser)
}
