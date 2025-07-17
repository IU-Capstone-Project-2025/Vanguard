package Storage

import (
	"fmt"
	"golang.org/x/net/context"
	"xxx/shared"
)

func (r *Redis) LoadLeaderboard(quizID string) ([]shared.UserScore, error) {
	key := "leaderboard:" + quizID
	ctx := context.Background()
	zs, err := r.Client.ZRevRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var scores []shared.UserScore
	for _, z := range zs {
		userID := fmt.Sprintf("%v", z.Member)
		scores = append(scores, shared.UserScore{
			UserId:     userID,
			TotalScore: int(z.Score),
		})
	}

	return scores, nil
}
