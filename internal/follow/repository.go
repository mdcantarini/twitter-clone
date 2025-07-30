package follow

import (
	"gorm.io/gorm"
)

type Repository interface {
	InsertFollow(db *gorm.DB, follow *Follow) error
	RemoveFollow(db *gorm.DB, followerID, followedID uint) error
}

func InsertFollow(db *gorm.DB, follow *Follow) error {
	return db.Create(follow).Error
}

func RemoveFollow(db *gorm.DB, followerID, followedID uint) error {
	return db.Where("follower_id = ? AND followed_id = ?", followerID, followedID).
		Delete(&Follow{}).Error
}
