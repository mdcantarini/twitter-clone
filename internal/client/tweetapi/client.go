package tweetapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Client interface {
	FetchTweet(tweetID string) (*Tweet, error)
}

type Tweet struct {
	TweetID   string `json:"TweetID"`
	UserID    uint   `json:"UserID"`
	Content   string `json:"Content"`
	CreatedAt string `json:"CreatedAt"`
}

func FetchTweet(tweetID string) (*Tweet, error) {
	url := fmt.Sprintf("http://%s/api/v1/tweet/%s", os.Getenv("TWEET_API_URL"), tweetID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call tweet-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tweet-api responded with status %d", resp.StatusCode)
	}

	var tweet Tweet
	if err := json.NewDecoder(resp.Body).Decode(&tweet); err != nil {
		return nil, fmt.Errorf("failed to decode tweet response: %w", err)
	}
	return &tweet, nil
}
