package Redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	models "xxx/SessionService/models"
	"xxx/shared"
)

type Cache interface {
	SaveSession(session *shared.Session) error
	LoadSession(code string) (*shared.Session, error)
	DeleteSession(code string) error
	EditSessionState(sessionID string, state string) error
	CodeExist(code string) bool
	GetPlayersForSession(sessionCode string) ([]string, error)
	AddPlayerToSession(sessionCode string, playerName string) error
}

type Redis struct {
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, cfg models.Config) (*Redis, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return &Redis{}, err
	}
	r := &Redis{Client: db}
	return r, nil
}

// SaveSession store session in Redis
func (r *Redis) SaveSession(session *shared.Session) error {
	ctx := context.Background()
	key := "session:" + session.Code
	err := r.Client.HSet(ctx, key, map[string]interface{}{
		"id":    session.ID,
		"code":  session.Code,
		"state": session.State,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

// LoadSession Load session by its Id from Redis
func (r *Redis) LoadSession(code string) (*shared.Session, error) {
	key := "session:" + code
	ctx := context.Background()
	res, err := r.Client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("session with code %s not found in Redis", code)
	}
	session := &shared.Session{
		ID:    res["id"],
		Code:  res["code"],
		State: res["state"],
	}
	return session, nil
}

// DeleteSession Delete session from Redis
func (r *Redis) DeleteSession(code string) error {
	key := "session:" + code
	ctx := context.Background()
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) EditSessionState(code string, state string) error {
	key := "session:" + code
	ctx := context.Background()
	err := r.Client.HSet(ctx, key, "state", state).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) CodeExist(code string) bool {
	key := "session:" + code
	ctx := context.Background()
	exist, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	if exist > 0 {
		return true
	} else {
		return false
	}
}

func (r *Redis) AddPlayerToSession(sessionCode string, playerName string) error {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s:players", sessionCode)
	err := r.Client.SAdd(ctx, key, playerName).Err()
	if err != nil {
		return err
	}

	r.Client.Expire(ctx, key, time.Hour)

	return nil
}

func (r *Redis) GetPlayersForSession(sessionCode string) ([]string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s:players", sessionCode)

	players, err := r.Client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return players, nil
}
