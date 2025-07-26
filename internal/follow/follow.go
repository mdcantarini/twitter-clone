package follow

import "time"

// The Follow model represents the following relationship between users.
type Follow struct {
	FollowerID uint `gorm:"primaryKey"`
	FollowedID uint `gorm:"primaryKey"`
	CreatedAt  time.Time
}
