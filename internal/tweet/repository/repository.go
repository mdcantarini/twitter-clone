package repository

import (
	"github.com/gocql/gocql"

	"github.com/mdcantarini/twitter-clone/internal/tweet"
)

type Repository interface {
	InsertTweet(tweet tweet.Tweet) error
	GetTweetsByUser(userID uint, limit uint) ([]tweet.Tweet, error)
	GetTweetById(tweetId gocql.UUID) (tweet.Tweet, error)
}