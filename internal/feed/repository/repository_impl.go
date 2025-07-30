package repository

import (
	"github.com/gocql/gocql"

	"github.com/mdcantarini/twitter-clone/internal/feed"
)

type NoSqlRepositoryImplementation struct {
	session *gocql.Session
}

func NewNoSqlRepositoryImplementation(session *gocql.Session) NoSqlRepositoryImplementation {
	return NoSqlRepositoryImplementation{session}
}

func (ci *NoSqlRepositoryImplementation) GetUserTimeline(
	userID uint,
	limit int,
) ([]feed.FeedEntry, error) {
	var feed []feed.FeedEntry
	iter := ci.session.Query(`
		SELECT tweet_id, author_id, created_at, content
		FROM user_timeline
		WHERE user_id = ?
		LIMIT ?`,
		userID, limit,
	).Iter()

	var entry feed.FeedEntry
	for iter.Scan(&entry.TweetID, &entry.AuthorID, &entry.CreatedAt, &entry.Content) {
		feed = append(feed, entry)
	}
	return feed, iter.Close()
}

func (ci *NoSqlRepositoryImplementation) InsertUserTimeline(
	followerIds []uint,
	createdAt string,
	tweetId string,
	userId uint,
	tweetContent string,
) error {
	batch := ci.session.NewBatch(gocql.LoggedBatch)

	for _, followerId := range followerIds {
		batch.Query(`
			INSERT INTO user_timeline (user_id, created_at, tweet_id, author_id, content)
			VALUES (?, ?, ?, ?, ?)`,
			followerId, createdAt, tweetId, userId, tweetContent,
		)
	}

	return ci.session.ExecuteBatch(batch)
}
