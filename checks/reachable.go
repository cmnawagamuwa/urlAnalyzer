package checks

import (
	"fmt"
	"net/http"
	"time"
)

// CheckReachable sends a GET request and returns the status, or an error if unreachable.
func CheckReachable(url string) (string, error) {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("Cannot reach URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return fmt.Sprintf("Success! Status: %d", resp.StatusCode), nil
	}
	return "", fmt.Errorf("Failed! Status: %d", resp.StatusCode)
}
