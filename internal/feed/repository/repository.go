package repository

import (
	"github.com/mdcantarini/twitter-clone/internal/feed/model"
)

type Repository interface {
	GetUserTimeline(
		userID uint,
		limit int,
	) ([]model.FeedEntry, error)
	InsertUserTimeline(
		followerIds []uint,
		createdAt string,
		tweetId string,
		userId uint,
		tweetContent string,
	) error
}
