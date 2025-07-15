package Storage

import (
	"fmt"
	"golang.org/x/net/context"
	"xxx/LeaderBoardService/models"
)

func (r *Redis) LoadLeaderboard(quizID string) ([]models.UserScore, error) {
	key := "leaderboard:" + quizID
	ctx := context.Background()
	zs, err := r.Client.ZRevRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var scores []models.UserScore
	for _, z := range zs {
		userID := fmt.Sprintf("%v", z.Member)
		scores = append(scores, models.UserScore{
			UserId:     userID,
			TotalScore: int(z.Score),
		})
	}

	return scores, nil
}
