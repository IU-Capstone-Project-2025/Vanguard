package Redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	models "xxx/SessionService/models"
)

type Cache interface {
	SaveSession(session *models.Session) error
	LoadSession(code string, sessionID string) (*models.Session, error)
	DeleteSession(sessionID string) error
	EditSessionState(sessionID string, state string) error
	CodeExist(code string) bool
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
func (r *Redis) SaveSession(session *models.Session) error {
	ctx := context.Background()
	key := "session:" + session.ID
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
func (r *Redis) LoadSession(code string, sessionID string) (*models.Session, error) {
	key := "session:" + code
	ctx := context.Background()
	res, err := r.Client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	session := &models.Session{
		ID:    sessionID,
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
