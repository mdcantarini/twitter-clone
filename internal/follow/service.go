package follow

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/follow/repository"
)

type Service struct {
	db repository.SqlRepositoryImplementation
}

func NewService(db *gorm.DB) *Service {
	sqlImpl := repository.NewSqlRepositoryImplementation(db)

	return &Service{db: sqlImpl}
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

	if err := s.db.InsertFollow(followData); err != nil {
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

	if err := s.db.RemoveFollow(uint(followerID), uint(followedID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

func (s *Service) GetFollowerIds(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	followers, err := s.db.GetFollowers(uint(userID))
	if err != nil {
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
	router.GET("/users/:user_id/follower_ids", s.GetFollowerIds)
}
