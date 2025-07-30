package repository

import (
	"github.com/gocql/gocql"
	"github.com/mdcantarini/twitter-clone/internal/tweet/model"
)

// TODO - Improve! Add test cases for real implementation
type NoSqlRepositoryImplementation struct {
	session *gocql.Session
}

func NewNoSqlRepositoryImplementation(session *gocql.Session) NoSqlRepositoryImplementation {
	return NoSqlRepositoryImplementation{session}
}

func (ni NoSqlRepositoryImplementation) InsertTweet(tweetData model.Tweet) error {
	batch := ni.session.NewBatch(gocql.LoggedBatch)

	// Insert into tweets_by_id (for lookup by ID)
	batch.Query(`
		INSERT INTO tweets_by_id (tweet_id, user_id, content, created_at)
		VALUES (?, ?, ?, ?)`,
		tweetData.TweetID, tweetData.UserID, tweetData.Content, tweetData.CreatedAt,
	)

	return ni.session.ExecuteBatch(batch)
}

func (ni NoSqlRepositoryImplementation) GetTweetById(tweetId gocql.UUID) (model.Tweet, error) {
	var tweetData model.Tweet

	err := ni.session.Query(`
		SELECT user_id, tweet_id, content, created_at
		FROM tweets_by_id
		WHERE tweet_id = ?`,
		tweetId,
	).Scan(&tweetData.UserID, &tweetData.TweetID, &tweetData.Content, &tweetData.CreatedAt)

	if err != nil {
		return model.Tweet{}, err
	}
	return tweetData, nil
}
