package tweet

import "github.com/gocql/gocql"

func InsertTweet(session *gocql.Session, tweet Tweet) error {
	batch := session.NewBatch(gocql.LoggedBatch)
	batch.Query(`
		INSERT INTO tweets_by_user (user_id, created_at, tweet_id, content)
		VALUES (?, ?, ?, ?)`,
		tweet.UserID, tweet.CreatedAt, tweet.TweetID, tweet.Content,
	)
	// TODO - Insert in tweet_by_id
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
