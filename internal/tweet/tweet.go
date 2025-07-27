package tweet

import (
	"time"

	"github.com/gocql/gocql"
)

// The Tweet model represents a single tweet/post.
// `Tweet` belongs to `User`, `UserID` is the foreign key
// UserID is provided by user-api, it ensures a uniqueness by ID
// TweetID instead is generated in runtime each time a new tweet is created
type Tweet struct {
	TweetID   gocql.UUID
	UserID    uint
	Content   string
	CreatedAt time.Time
}
