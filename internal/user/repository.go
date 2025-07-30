package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	InsertUser(db *gorm.DB, user *User) (*User, error)
	GetUser(db *gorm.DB, id uint) (*User, error)
}

func InsertUser(db *gorm.DB, user *User) (*User, error) {
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetUser(db *gorm.DB, id uint) (*User, error) {
	user := &User{}
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}
