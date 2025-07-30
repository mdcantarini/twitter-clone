package feed

import (
	"github.com/gocql/gocql"
)

type Repository interface {
	GetUserTimeline(
		session *gocql.Session,
		userID uint,
		limit int,
	) ([]FeedEntry, error)
	InsertUserTimeline(
		session *gocql.Session,
		followerIds []uint,
		createdAt string,
		tweetId string,
		userId uint,
		tweetContent string,
	) error
}

func GetUserTimeline(
	session *gocql.Session,
	userID uint,
	limit int,
) ([]FeedEntry, error) {
	var feed []FeedEntry
	iter := session.Query(`
		SELECT tweet_id, author_id, created_at, content
		FROM user_timeline
		WHERE user_id = ?
		LIMIT ?`,
		userID, limit,
	).Iter()

	var entry FeedEntry
	for iter.Scan(&entry.TweetID, &entry.AuthorID, &entry.CreatedAt, &entry.Content) {
		feed = append(feed, entry)
	}
	return feed, iter.Close()
}

func InsertUserTimeline(
	session *gocql.Session,
	followerIds []uint,
	createdAt string,
	tweetId string,
	userId uint,
	tweetContent string,
) error {
	batch := session.NewBatch(gocql.LoggedBatch)

	for _, followerId := range followerIds {
		batch.Query(`
			INSERT INTO user_timeline (user_id, created_at, tweet_id, author_id, content)
			VALUES (?, ?, ?, ?, ?)`,
			followerId, createdAt, tweetId, userId, tweetContent,
		)
	}

	return session.ExecuteBatch(batch)
}
