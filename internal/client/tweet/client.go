package tweet

type Client interface {
	FetchTweet(tweetID string) (*Tweet, error)
}
