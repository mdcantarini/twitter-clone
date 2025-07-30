package repository

import (
	"github.com/mdcantarini/twitter-clone/internal/user/model"
	"gorm.io/gorm"
)

// TODO - Improve! Add test cases for real implementation
type SqlRepositoryImplementation struct {
	db *gorm.DB
}

func NewSqlRepositoryImplementation(db *gorm.DB) SqlRepositoryImplementation {
	return SqlRepositoryImplementation{db}
}

func (si SqlRepositoryImplementation) InsertUser(userData *model.User) (*model.User, error) {
	if err := si.db.Create(userData).Error; err != nil {
		return nil, err
	}

	return userData, nil
}

func (si SqlRepositoryImplementation) GetUser(id uint) (*model.User, error) {
	userData := &model.User{}
	if err := si.db.First(&userData, id).Error; err != nil {
		return nil, err
	}

	return userData, nil
}
