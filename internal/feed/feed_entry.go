package feed

import (
	"time"

	"github.com/gocql/gocql"
)

// The FeedEntry model represents the association between users and tweets in their feed.
type FeedEntry struct {
	TweetID   gocql.UUID
	AuthorID  uint
	Content   string
	CreatedAt time.Time
}
