package feed

import (
	"gorm.io/gorm"

	"github.com/mdcantarini/twitter-clone/internal/tweet"
	"github.com/mdcantarini/twitter-clone/internal/user"
)

// The Feed model represents the association between users and tweets in their feed.
type Feed struct {
	gorm.Model
	UserID  uint
	User    user.User
	TweetID uint
	Tweet   tweet.Tweet
}
