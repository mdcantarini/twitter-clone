package repository

import (
	"github.com/gocql/gocql"

	"github.com/mdcantarini/twitter-clone/internal/tweet"
)

type NoSqlRepositoryImplementation struct {
	session *gocql.Session
}

func NewNoSqlRepositoryImplementation(session *gocql.Session) NoSqlRepositoryImplementation {
	return NoSqlRepositoryImplementation{session}
}

func (ni *NoSqlRepositoryImplementation) InsertTweet(tweetData tweet.Tweet) error {
	batch := ni.session.NewBatch(gocql.LoggedBatch)

	// Insert into tweets_by_user (for user's own profile view)
	batch.Query(`
		INSERT INTO tweets_by_user (user_id, created_at, tweet_id, content)
		VALUES (?, ?, ?, ?)`,
		tweetData.UserID, tweetData.CreatedAt, tweetData.TweetID, tweetData.Content,
	)

	// Insert into tweets_by_id (for lookup by ID)
	batch.Query(`
		INSERT INTO tweets_by_id (tweet_id, user_id, content, created_at)
		VALUES (?, ?, ?, ?)`,
		tweetData.TweetID, tweetData.UserID, tweetData.Content, tweetData.CreatedAt,
	)

	return ni.session.ExecuteBatch(batch)
}

func (ni *NoSqlRepositoryImplementation) GetTweetsByUser(userID uint, limit uint) ([]tweet.Tweet, error) {
	var tweets []tweet.Tweet

	iter := ni.session.Query(`
		SELECT user_id, tweet_id, content, created_at
		FROM tweets_by_user
		WHERE user_id = ?
		LIMIT ?`,
		userID, limit,
	).Iter()

	var tweetData tweet.Tweet
	for iter.Scan(&tweetData.UserID, &tweetData.TweetID, &tweetData.Content, &tweetData.CreatedAt) {
		tweets = append(tweets, tweetData)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return tweets, nil
}

func (ni *NoSqlRepositoryImplementation) GetTweetById(tweetId gocql.UUID) (tweet.Tweet, error) {
	var tweetData tweet.Tweet

	err := ni.session.Query(`
		SELECT user_id, tweet_id, content, created_at
		FROM tweets_by_id
		WHERE tweet_id = ?`,
		tweetId,
	).Scan(&tweetData.UserID, &tweetData.TweetID, &tweetData.Content, &tweetData.CreatedAt)

	if err != nil {
		return tweet.Tweet{}, err
	}
	return tweetData, nil
}