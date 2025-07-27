package user

import "time"

// The User model represents a user account in the system.
type User struct {
	ID          uint   `gorm:"primaryKey"`
	Username    string `gorm:"uniqueIndex"`
	DisplayName string
	CreatedAt   time.Time
}
