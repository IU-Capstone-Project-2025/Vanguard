package Storage

import (
	"fmt"
	"golang.org/x/net/context"
	"xxx/shared"
)

func (r *Redis) LoadLeaderboard(quizID string, previousPlaces map[string]int) ([]shared.UserScore, error) {
	key := "leaderboard:" + quizID
	ctx := context.Background()

	zs, err := r.Client.ZRevRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var scores []shared.UserScore
	for i, z := range zs {
		userID := fmt.Sprintf("%v", z.Member)
		newPlace := i + 1
		prevPlace := previousPlaces[userID]

		progress := shared.ProgressSame
		if prevPlace > 0 && newPlace < prevPlace {
			progress = shared.ProgressUp
		} else if prevPlace > 0 && newPlace == prevPlace {
			progress = shared.ProgressSame
		} else if prevPlace > 0 && newPlace > prevPlace {
			progress = shared.ProgressDown
		}
		scores = append(scores, shared.UserScore{
			UserId:        userID,
			TotalScore:    int(z.Score),
			Place:         newPlace,
			PreviousPlace: prevPlace,
			Progress:      progress,
		})
	}

	return scores, nil
}
