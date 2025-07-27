package follow

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

func (h *Handler) FollowUser(c *gin.Context) {
	var input struct {
		FollowerID uint `json:"follower_id" binding:"required"`
		FollowedID uint `json:"followed_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.FollowerID == input.FollowedID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	followData := &Follow{
		FollowerID: input.FollowerID,
		FollowedID: input.FollowedID,
	}

	if err := InsertFollow(h.db, followData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Followed successfully"})
}

func (h *Handler) UnfollowUser(c *gin.Context) {
	followerID, err := strconv.ParseUint(c.Param("follower_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid follower ID"})
		return
	}

	followedID, err := strconv.ParseUint(c.Param("followed_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid followed ID"})
		return
	}

	if err := RemoveFollow(h.db, uint(followerID), uint(followedID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/follow", h.FollowUser)
	router.DELETE("/follow/:follower_id/:followed_id", h.UnfollowUser)
}
