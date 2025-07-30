package repository

import (
	"github.com/mdcantarini/twitter-clone/internal/follow/model"
)

type Repository interface {
	InsertFollow(follow *model.Follow) error
	GetFollowers(followedID uint) ([]model.Follow, error)
}
