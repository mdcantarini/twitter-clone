package user

import (
	"gorm.io/gorm"
)

// The User model represents a user account in the system.
type User struct {
	gorm.Model
	Username    string `gorm:"uniqueIndex"`
	DisplayName string
}
