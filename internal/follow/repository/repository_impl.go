package repository

import (
	"github.com/mdcantarini/twitter-clone/internal/follow/model"
	"gorm.io/gorm"
)

// TODO - Improve! Add test cases for real implementation
type SqlRepositoryImplementation struct {
	db *gorm.DB
}

func NewSqlRepositoryImplementation(db *gorm.DB) SqlRepositoryImplementation {
	return SqlRepositoryImplementation{db}
}

func (si SqlRepositoryImplementation) InsertFollow(follow *model.Follow) error {
	return si.db.Create(follow).Error
}

func (si SqlRepositoryImplementation) GetFollowers(followedID uint) ([]model.Follow, error) {
	var followers []model.Follow
	if err := si.db.Where("followed_id = ?", followedID).Find(&followers).Error; err != nil {
		return nil, err
	}

	return followers, nil
}
