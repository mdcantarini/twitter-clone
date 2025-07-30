package repository

import (
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/follow"
)

type Repository interface {
	InsertFollow(db *gorm.DB, follow *follow.Follow) error
	RemoveFollow(db *gorm.DB, followerID, followedID uint) error
	GetFollowers(db *gorm.DB, followedID uint) ([]follow.Follow, error)
}
