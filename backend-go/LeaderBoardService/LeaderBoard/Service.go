package LeaderBoard

import (
	"golang.org/x/net/context"
	"log/slog"
	"xxx/LeaderBoardService/Storage"
	models2 "xxx/LeaderBoardService/models"
	"xxx/shared"
)

type Service interface {
	ComputeLeaderBoard(ans shared.SessionAnswers) (shared.ScoreTable, error)
	PopularAns(ans shared.SessionAnswers) (shared.PopularAns, error)
}

type LeaderBoard struct {
	log   *slog.Logger
	Cache Storage.Cache
}

func NewLeaderBoard(log *slog.Logger, redisConn string) (*LeaderBoard, error) {
	RedisConfig := models2.Config{
		Addr:        redisConn,
		Password:    "",
		DB:          0,
		MaxRetries:  0,
		DialTimeout: 0,
		Timeout:     0,
	}
	Cache, err := Storage.NewRedisClient(context.Background(), RedisConfig)
	if err != nil {
		return nil, err
	}
	return &LeaderBoard{log: log, Cache: Cache}, nil
}
