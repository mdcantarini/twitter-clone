package tweet

import (
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/user"
)

// The Tweet model represents a single tweet/post.
// `Tweet` belongs to `User`, `UserID` is the foreign key
type Tweet struct {
	gorm.Model
	UserID  uint
	User    user.User
	Content string
}
