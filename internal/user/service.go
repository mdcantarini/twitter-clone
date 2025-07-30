package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/user/repository"
)

type Service struct {
	db repository.SqlRepositoryImplementation
}

func NewService(db *gorm.DB) *Service {
	sqlImpl := repository.NewSqlRepositoryImplementation(db)

	return &Service{db: sqlImpl}
}

func (s *Service) CreateUser(c *gin.Context) {
	var input struct {
		Username    string `json:"username" binding:"required"`
		DisplayName string `json:"display_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := &User{
		Username:    input.Username,
		DisplayName: input.DisplayName,
	}

	user, err := s.db.InsertUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (s *Service) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := s.db.GetUser(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/users", s.CreateUser)
	router.GET("/users/:id", s.GetUser)
}
