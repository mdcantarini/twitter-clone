package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func FetchFollowers(userID uint) ([]uint, error) {
	url := fmt.Sprintf("http://%s/api/v1/users/%d/followers", os.Getenv("FOLLOW_API_URL"), userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call follow-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("follow-api responded with status %d", resp.StatusCode)
	}

	var result struct {
		FollowerIDs []uint `json:"follower_ids"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode follower response: %w", err)
	}
	return result.FollowerIDs, nil
}
