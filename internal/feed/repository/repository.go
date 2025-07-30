package repository

import (
	"github.com/gocql/gocql"

	"github.com/mdcantarini/twitter-clone/internal/feed"
)

type Repository interface {
	GetUserTimeline(
		session *gocql.Session,
		userID uint,
		limit int,
	) ([]feed.FeedEntry, error)
	InsertUserTimeline(
		session *gocql.Session,
		followerIds []uint,
		createdAt string,
		tweetId string,
		userId uint,
		tweetContent string,
	) error
}
