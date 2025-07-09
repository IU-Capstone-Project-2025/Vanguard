package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"xxx/real_time/models"
)

type Cache interface {
	SetSessionQuiz(sessionId string, quizData models.OngoingQuiz) error
	GetSessionQuiz(sessionId string) (models.OngoingQuiz, error)
	DeleteSession(sessionId string) error

	SetQuestionIndex(sessionId string, questionIdx int) error
	GetQuestionIndex(sessionId string) (int, error)

	RecordAnswer(sessionID, userID string, isCorrect bool) error
	GetAllAnswer(sessionId string) (map[string]bool, error)
}

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
