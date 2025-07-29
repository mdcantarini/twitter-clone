package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Tweet struct {
	TweetID   string `json:"tweet_id"`
	UserID    uint   `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
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
