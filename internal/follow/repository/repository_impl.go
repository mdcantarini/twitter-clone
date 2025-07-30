package repository

import (
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/follow"
)

type SqlRepositoryImplementation struct {
	db *gorm.DB
}

func NewSqlRepositoryImplementation(db *gorm.DB) SqlRepositoryImplementation {
	return SqlRepositoryImplementation{db}
}

func (si SqlRepositoryImplementation) InsertFollow(follow *follow.Follow) error {
	return si.db.Create(follow).Error
}

func (si SqlRepositoryImplementation) RemoveFollow(followerID, followedID uint) error {
	return si.db.Where("follower_id = ? AND followed_id = ?", followerID, followedID).
		Delete(&follow.Follow{}).Error
}

func (si SqlRepositoryImplementation) GetFollowers(followedID uint) ([]follow.Follow, error) {
	var followers []follow.Follow
	if err := si.db.Where("followed_id = ?", followedID).Find(&followers).Error; err != nil {
		return nil, err
	}

	return followers, nil
}
