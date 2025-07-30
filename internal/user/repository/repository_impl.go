package repository

import (
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/user"
)

type SqlRepositoryImplementation struct {
	db *gorm.DB
}

func NewSqlRepositoryImplementation(db *gorm.DB) SqlRepositoryImplementation {
	return SqlRepositoryImplementation{db}
}

func (si *SqlRepositoryImplementation) InsertUser(userData *user.User) (*user.User, error) {
	if err := si.db.Create(userData).Error; err != nil {
		return nil, err
	}

	return userData, nil
}

func (si *SqlRepositoryImplementation) GetUser(id uint) (*user.User, error) {
	userData := &user.User{}
	if err := si.db.First(&userData, id).Error; err != nil {
		return nil, err
	}

	return userData, nil
}