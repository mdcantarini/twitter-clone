package feed

import "time"

// The FeedEntry model represents the association between users and tweets in their feed.
type FeedEntry struct {
	TweetID   uint
	AuthorID  uint
	CreatedAt time.Time
}
