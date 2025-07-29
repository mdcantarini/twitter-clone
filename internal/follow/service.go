package follow

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) FollowUser(c *gin.Context) {
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

	if err := InsertFollow(s.db, followData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Followed successfully"})
}

func (s *Service) UnfollowUser(c *gin.Context) {
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

	if err := RemoveFollow(s.db, uint(followerID), uint(followedID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

func (s *Service) GetFollowers(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var followers []Follow
	if err := s.db.Where("followed_id = ?", userID).Find(&followers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get followers"})
		return
	}

	followerIDs := make([]uint, len(followers))
	for i, f := range followers {
		followerIDs[i] = f.FollowerID
	}

	c.JSON(http.StatusOK, gin.H{"follower_ids": followerIDs})
}

func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/follow", s.FollowUser)
	router.DELETE("/follow/:follower_id/:followed_id", s.UnfollowUser)
	router.GET("/users/:user_id/followers", s.GetFollowers)
}
