package Storage

import (
	"context"
	"xxx/LeaderBoardService/models"
)

func (r *Redis) AddScoresBatch(quizID string, updates []models.UserCurrentPoint) error {
	key := "leaderboard:" + quizID
	pipe := r.Client.Pipeline()
	ctx := context.Background()
	for _, update := range updates {
		pipe.ZIncrBy(ctx, key, float64(update.Score), update.UserId)
	}
	_, err := pipe.Exec(ctx)
	return err
}
