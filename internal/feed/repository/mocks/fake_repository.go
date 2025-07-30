package mocks

import (
	"github.com/mdcantarini/twitter-clone/internal/feed/model"
)

type FakeNoSqlRepository struct {
	GetUserTimelineFunc    func(userID uint, limit int) ([]model.FeedEntry, error)
	InsertUserTimelineFunc func(followerIds []uint, createdAt string, tweetId string, userId uint, tweetContent string) error
}

func (f *FakeNoSqlRepository) GetUserTimeline(userID uint, limit int) ([]model.FeedEntry, error) {
	if f.GetUserTimelineFunc != nil {
		return f.GetUserTimelineFunc(userID, limit)
	}
	return []model.FeedEntry{}, nil
}

func (f *FakeNoSqlRepository) InsertUserTimeline(followerIds []uint, createdAt string, tweetId string, userId uint, tweetContent string) error {
	if f.InsertUserTimelineFunc != nil {
		return f.InsertUserTimelineFunc(followerIds, createdAt, tweetId, userId, tweetContent)
	}
	return nil
}
