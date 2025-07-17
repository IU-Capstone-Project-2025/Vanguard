package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"xxx/real_time/models"
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

// Ping helper
func (c *Client) Ping() error { return c.rdb.Ping(c.ctx).Err() }

// Expire sets TTL for a key; exposed for tests/maintenance.
func (c *Client) Expire(key string, ttl time.Duration) error {
	return c.rdb.Expire(c.ctx, key, ttl).Err()
}

// SetSessionQuiz stores the ongoing quiz data as all questions and current question index for a given session with a TTL.
func (c *Client) SetSessionQuiz(sessionID string, quizData models.OngoingQuiz) error {
	key := fmt.Sprintf("session:%s:quiz_state", sessionID)
	data, err := json.Marshal(quizData)
	if err != nil {
		return err
	}
	return c.rdb.Set(c.ctx, key, data, 24*time.Hour).Err()
}

// GetSessionQuiz retrieves the stored quiz data as all questions and current question index for a session.
func (c *Client) GetSessionQuiz(sessionID string) (models.OngoingQuiz, error) {
	key := fmt.Sprintf("session:%s:quiz_state", sessionID)
	rawVal, err := c.rdb.Get(c.ctx, key).Bytes()
	if err != nil {
		return models.OngoingQuiz{}, err
	}

	var quiz models.OngoingQuiz
	if err = json.Unmarshal(rawVal, &quiz); err != nil {
		return models.OngoingQuiz{}, err
	}
	return quiz, nil
}

// DeleteSession clears all Redis keys related to a session.
func (c *Client) DeleteSession(sessionID string) error {
	keys := []string{
		fmt.Sprintf("session:%s:quiz_state", sessionID),
	}
	// Pattern for user-specific answer keys
	pattern := fmt.Sprintf("session:%s:user:*:answers", sessionID)

	// Use SCAN to find matching keys
	var cursor uint64
	for {
		matchedKeys, nextCursor, err := c.rdb.Scan(c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		keys = append(keys, matchedKeys...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	// Delete all collected keys
	if len(keys) > 0 {
		return c.rdb.Del(c.ctx, keys...).Err()
	}
	return nil
}

// GetAllSessions retrieves all sessions states.
func (c *Client) GetAllSessions() (map[string]models.OngoingQuiz, error) {
	sessions := make(map[string]models.OngoingQuiz)
	// Pattern for user-specific answer keys
	pattern := "session:*:quiz_state"

	// Use SCAN to find matching keys
	var cursor uint64
	for {
		keys, nextCursor, err := c.rdb.Scan(c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			// key format: session:<sessionID>:quiz_state
			// split once on ':' to extract sessionID at index 1
			parts := strings.Split(key, ":")
			if len(parts) >= 3 {
				sessionId := parts[1]
				quiz, _ := c.GetSessionQuiz(sessionId)
				sessions[sessionId] = quiz
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return sessions, nil
}

// SetQuestionIndex stores the current question index for a session with a TTL.
func (c *Client) SetQuestionIndex(sessionID string, idx int) error {
	quizState, err := c.GetSessionQuiz(sessionID)
	if err != nil {
		return err
	}
	quizState.CurrQuestionIdx = idx

	return c.SetSessionQuiz(sessionID, quizState)
}

// GetQuestionIndex retrieves the stored question index for a session.
func (c *Client) GetQuestionIndex(sessionID string) (int, error) {
	quizState, err := c.GetSessionQuiz(sessionID)
	if err != nil {
		return 0, err
	}
	return quizState.CurrQuestionIdx, nil
}

// RecordAnswer stores correctness of a user's answer in a Redis hash.
func (c *Client) RecordAnswer(sessionID, userID string, question int, answer models.UserAnswer) error {
	hash := fmt.Sprintf("session:%s:user:%s:answers", sessionID, userID)

	data, err := json.Marshal(answer)
	if err != nil {
		return err
	}

	return c.rdb.HSet(c.ctx, hash, question, data).Err()
}

// GetAllAnswers retrieves every user's recorded answers for a given session.
//
// Result format:
//
//	map[ userID ] -> []models.UserAnswer
//
// Each models.UserAnswer contains Question index and IsCorrect flag (or any
// extra fields you defined in that struct).
func (c *Client) GetAllAnswers(sessionID string) (map[string][]models.UserAnswer, error) {
	result := make(map[string][]models.UserAnswer)

	// Keys look like: session:<sessionID>:user:<userID>:answers
	pattern := fmt.Sprintf("session:%s:user:*:answers", sessionID)

	var cursor uint64
	for {
		keys, nextCursor, err := c.rdb.Scan(c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			parts := strings.Split(key, ":")
			if len(parts) < 5 {
				continue // unexpected format, skip
			}
			userID := parts[3] // session:<id>:user:<userID>:answers

			hashData, err := c.rdb.HGetAll(c.ctx, key).Result()
			if err != nil {
				return nil, err
			}

			for _, raw := range hashData {
				var ans models.UserAnswer
				if err := json.Unmarshal([]byte(raw), &ans); err != nil {
					// skip malformed answer but continue collecting others
					continue
				}
				result[userID] = append(result[userID], ans)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return result, nil
}
