package feed

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/follow"
	"github.com/mdcantarini/twitter-clone/internal/tweet"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetUserFeed(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var feedItems []Feed
	if err := h.db.Preload("Tweet").
		Preload("Tweet.User").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&feedItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feed"})
		return
	}

	tweets := make([]tweet.Tweet, len(feedItems))
	for i, item := range feedItems {
		tweets[i] = item.Tweet
	}

	c.JSON(http.StatusOK, tweets)
}

func (h *Handler) GetUserTimeline(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var followedUserIDs []uint
	if err := h.db.Model(&follow.Follow{}).
		Where("follower_id = ?", userID).
		Pluck("followed_id", &followedUserIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve followed users"})
		return
	}

	followedUserIDs = append(followedUserIDs, uint(userID))

	var tweets []tweet.Tweet
	if err := h.db.Preload("User").
		Where("user_id IN ?", followedUserIDs).
		Order("created_at desc").
		Find(&tweets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve timeline"})
		return
	}

	c.JSON(http.StatusOK, tweets)
}

func (h *Handler) AddToFeed(c *gin.Context) {
	var input struct {
		UserID  uint `json:"user_id" binding:"required"`
		TweetID uint `json:"tweet_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedItem := Feed{
		UserID:  input.UserID,
		TweetID: input.TweetID,
	}

	if err := h.db.Create(&feedItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to feed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Added to feed successfully"})
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/feed/:user_id", h.GetUserFeed)
	router.GET("/timeline/:user_id", h.GetUserTimeline)
	router.POST("/feed", h.AddToFeed)
}
