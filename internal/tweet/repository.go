package tweet

import "github.com/gocql/gocql"

type Repository interface {
	InsertTweet(session *gocql.Session, tweet Tweet) error
	GetTweetsByUser(session *gocql.Session, userID uint, limit uint) ([]Tweet, error)
	GetTweetById(session *gocql.Session, tweetId gocql.UUID) (Tweet, error)
}

func InsertTweet(session *gocql.Session, tweet Tweet) error {
	batch := session.NewBatch(gocql.LoggedBatch)

	// Insert into tweets_by_user (for user’s own profile view)
	batch.Query(`
		INSERT INTO tweets_by_user (user_id, created_at, tweet_id, content)
		VALUES (?, ?, ?, ?)`,
		tweet.UserID, tweet.CreatedAt, tweet.TweetID, tweet.Content,
	)

	// Insert into tweets_by_id (for lookup by ID)
	batch.Query(`
		INSERT INTO tweets_by_id (tweet_id, user_id, content, created_at)
		VALUES (?, ?, ?, ?)`,
		tweet.TweetID, tweet.UserID, tweet.Content, tweet.CreatedAt,
	)

	return session.ExecuteBatch(batch)
}

func GetTweetsByUser(session *gocql.Session, userID uint, limit uint) ([]Tweet, error) {
	var tweets []Tweet

	iter := session.Query(`
		SELECT user_id, tweet_id, content, created_at
		FROM tweets_by_user
		WHERE user_id = ?
		LIMIT ?`,
		userID, limit,
	).Iter()

	var tweet Tweet
	for iter.Scan(&tweet.UserID, &tweet.TweetID, &tweet.Content, &tweet.CreatedAt) {
		tweets = append(tweets, tweet)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return tweets, nil
}

func GetTweetById(session *gocql.Session, tweetId gocql.UUID) (Tweet, error) {
	var tweet Tweet

	err := session.Query(`
		SELECT user_id, tweet_id, content, created_at
		FROM tweets_by_id
		WHERE tweet_id = ?`,
		tweetId,
	).Scan(&tweet.UserID, &tweet.TweetID, &tweet.Content, &tweet.CreatedAt)

	if err != nil {
		return Tweet{}, err
	}
	return tweet, nil
}
