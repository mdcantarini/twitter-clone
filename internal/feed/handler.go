package feed

import (
	"github.com/gocql/gocql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *gocql.Session
}

func NewHandler(db *gocql.Session) *Handler {
	return &Handler{db: db}
}

const userFeedLimit = 50

func (h *Handler) GetUserFeed(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	feedItems, err := GetUserFeed(h.db, uint(id), userFeedLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feed"})
		return
	}

	c.JSON(http.StatusOK, feedItems)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/feed/:user_id", h.GetUserFeed)
}
