package mocks

import (
	"github.com/mdcantarini/twitter-clone/internal/tweet/model"
	"time"

	"github.com/gocql/gocql"
)

type FakeNoSqlRepository struct {
	InsertTweetFunc  func(tweet model.Tweet) error
	GetTweetByIdFunc func(tweetId gocql.UUID) (model.Tweet, error)
}

func (f *FakeNoSqlRepository) InsertTweet(tweet model.Tweet) error {
	if f.InsertTweetFunc != nil {
		return f.InsertTweetFunc(tweet)
	}
	return nil
}

func (f *FakeNoSqlRepository) GetTweetById(tweetId gocql.UUID) (model.Tweet, error) {
	if f.GetTweetByIdFunc != nil {
		return f.GetTweetByIdFunc(tweetId)
	}
	return model.Tweet{
		TweetID:   tweetId,
		UserID:    1,
		Content:   "Test tweet",
		CreatedAt: time.Now(),
	}, nil
}
