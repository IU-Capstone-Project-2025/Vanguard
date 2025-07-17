package Storage

import (
	"context"
	"xxx/shared"
)

func (r *Redis) AddScoresBatch(quizID string, updates []shared.UserCurrentPoint) error {
	key := "leaderboard:" + quizID
	pipe := r.Client.Pipeline()
	ctx := context.Background()
	for _, update := range updates {
		pipe.ZIncrBy(ctx, key, float64(update.Score), update.UserId)
	}
	_, err := pipe.Exec(ctx)
	return err
}
