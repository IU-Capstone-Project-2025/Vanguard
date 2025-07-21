package Storage

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"xxx/LeaderBoardService/models"
	"xxx/shared"
)

type Cache interface {
	LoadLeaderboard(quizID string, previousPlaces map[string]int) ([]shared.UserScore, error)
	AddScoresBatch(quizID string, updates []shared.UserCurrentPoint) error
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
