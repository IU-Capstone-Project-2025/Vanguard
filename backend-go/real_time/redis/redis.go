package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// Client wraps a Redis client with helper methods for RealTime Service.
type Client struct {
	rdb *redis.Client
	ctx context.Context
}

// NewClient initializes a new Redis client.
func NewClient(addr, password string, db int) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Client{rdb: rdb, ctx: context.Background()}
}

// SetQuestionIndex stores the current question index for a session with a TTL.
func (c *Client) SetQuestionIndex(sessionID string, idx int) error {
	key := fmt.Sprintf("session:%s:question_idx", sessionID)
	return c.rdb.Set(c.ctx, key, idx, 24*time.Hour).Err()
}

// GetQuestionIndex retrieves the stored question index for a session.
func (c *Client) GetQuestionIndex(sessionID string) (int, error) {
	key := fmt.Sprintf("session:%s:question_idx", sessionID)
	val, err := c.rdb.Get(c.ctx, key).Int()
	if err == redis.Nil {
		return 0, nil // default to 0 if not set
	}
	return val, err
}

// RecordAnswer stores correctness of a user's answer in a Redis hash.
func (c *Client) RecordAnswer(sessionID, userID string, correct bool) error {
	hash := fmt.Sprintf("session:%s:answers", sessionID)
	return c.rdb.HSet(c.ctx, hash, userID, correct).Err()
}

// GetAllAnswers retrieves all recorded answers for a session.
func (c *Client) GetAllAnswers(sessionID string) (map[string]bool, error) {
	hash := fmt.Sprintf("session:%s:answers", sessionID)
	data, err := c.rdb.HGetAll(c.ctx, hash).Result()
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool, len(data))
	for user, str := range data {
		// Redis stores "true"/"false"
		val := str == "1" || str == "true"
		result[user] = val
	}
	return result, nil
}

// DeleteSession clears all Redis keys related to a session.
func (c *Client) DeleteSession(sessionID string) error {
	keys := []string{
		fmt.Sprintf("session:%s:question_idx", sessionID),
		fmt.Sprintf("session:%s:answers", sessionID),
	}
	return c.rdb.Del(c.ctx, keys...).Err()
}
