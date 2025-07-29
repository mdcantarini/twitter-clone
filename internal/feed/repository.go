package feed

import (
	"github.com/gocql/gocql"
)

func GetUserFeed(session *gocql.Session, userID uint, limit int) ([]FeedEntry, error) {
	var feed []FeedEntry
	iter := session.Query(`
		SELECT tweet_id, author_id, created_at
		FROM user_timeline
		WHERE user_id = ?
		LIMIT ?`,
		userID, limit,
	).Iter()

	var entry FeedEntry
	for iter.Scan(&entry.TweetID, &entry.AuthorID, &entry.CreatedAt) {
		feed = append(feed, entry)
	}
	return feed, iter.Close()
}
