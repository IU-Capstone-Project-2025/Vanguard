package leaderboard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
	"xxx/shared"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient returns a leaderboard client with sane timeouts.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   3 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConnsPerHost: 10,
			},
		},
	}
}

// GetResults sends answers for one session and returns the leaderboard.
func (c *Client) GetResults(ctx context.Context, sessionCode string, answers []shared.Answer) (shared.BoardResponse, error) {

	reqBody, err := json.Marshal(shared.SessionAnswers{
		SessionCode: sessionCode,
		Answers:     answers,
	})
	if err != nil {
		return shared.BoardResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	u := url.URL{
		Scheme: "http",
		Host:   c.baseURL,
		Path:   "/get-results",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(reqBody))
	if err != nil {
		return shared.BoardResponse{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return shared.BoardResponse{}, fmt.Errorf("post leaderboard: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return shared.BoardResponse{}, fmt.Errorf("leaderboard error %d: %s", resp.StatusCode, string(b))
	}

	var board shared.BoardResponse
	if err := json.NewDecoder(resp.Body).Decode(&board); err != nil {
		return shared.BoardResponse{}, fmt.Errorf("decode response: %w", err)
	}

	return board, nil
}
