package mocks

import (
	"github.com/mdcantarini/twitter-clone/internal/follow/model"
)

type FakeSqlRepository struct {
	InsertFollowFunc func(follow *model.Follow) error
	GetFollowersFunc func(followedID uint) ([]model.Follow, error)
}

func (f *FakeSqlRepository) InsertFollow(follow *model.Follow) error {
	if f.InsertFollowFunc != nil {
		return f.InsertFollowFunc(follow)
	}
	return nil
}

func (f *FakeSqlRepository) GetFollowers(followedID uint) ([]model.Follow, error) {
	if f.GetFollowersFunc != nil {
		return f.GetFollowersFunc(followedID)
	}
	return []model.Follow{
		{FollowerID: 2, FollowedID: followedID},
		{FollowerID: 3, FollowedID: followedID},
	}, nil
}
