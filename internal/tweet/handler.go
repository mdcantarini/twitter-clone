package tweet

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
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

	newTweet := Tweet{
		UserID:  input.UserID,
		Content: input.Content,
	}

	if err := h.db.Create(&newTweet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tweet"})
		return
	}

	if err := h.db.Preload("User").First(&newTweet, newTweet.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load tweet details"})
		return
	}

	c.JSON(http.StatusCreated, newTweet)
}

func (h *Handler) GetTweet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tweet ID"})
		return
	}

	var tweetData Tweet
	if err := h.db.Preload("User").First(&tweetData, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tweet"})
		return
	}

	c.JSON(http.StatusOK, tweetData)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/tweets", h.CreateTweet)
	router.GET("/tweets/:id", h.GetTweet)
}
