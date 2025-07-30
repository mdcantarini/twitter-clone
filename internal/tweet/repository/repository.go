package repository

import (
	"github.com/gocql/gocql"
	"github.com/mdcantarini/twitter-clone/internal/tweet/model"
)

type Repository interface {
	InsertTweet(tweet model.Tweet) error
	GetTweetById(tweetId gocql.UUID) (model.Tweet, error)
}
