package tweet

import (
	"time"

	"github.com/gocql/gocql"
)

// The Tweet model represents a single tweet/post.
// `Tweet` belongs to `User`, `UserID` is the foreign key
type Tweet struct {
	TweetID   gocql.UUID
	UserID    uint
	Content   string
	CreatedAt time.Time
}
