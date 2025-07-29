-- Create keyspace for Twitter clone
CREATE KEYSPACE IF NOT EXISTS twitter_clone
    WITH replication = {
        'class': 'SimpleStrategy',
        'replication_factor': 1
    };

USE twitter_clone;

-- Table for tweets organized by user (for user timeline queries)
CREATE TABLE IF NOT EXISTS tweets_by_user (
    user_id bigint,
    tweet_id uuid,
    content text,
    created_at timestamp,
    PRIMARY KEY (user_id, created_at, tweet_id)
) WITH CLUSTERING ORDER BY (created_at DESC, tweet_id DESC);

-- Table for tweets by ID (for single tweet lookups)
CREATE TABLE IF NOT EXISTS tweets_by_id (
    tweet_id uuid PRIMARY KEY,
    user_id bigint,
    content text,
    created_at timestamp
);

-- Create indexes if needed
CREATE INDEX IF NOT EXISTS ON tweets_by_user (tweet_id);

-- Optional: Table for user timelines (home feed)
-- This would store aggregated tweets from followed users
CREATE TABLE IF NOT EXISTS user_timeline (
    user_id bigint,
    tweet_id uuid,
    author_id bigint,
    content text,
    created_at timestamp,
    PRIMARY KEY (user_id, created_at, tweet_id)
) WITH CLUSTERING ORDER BY (created_at DESC, tweet_id DESC);