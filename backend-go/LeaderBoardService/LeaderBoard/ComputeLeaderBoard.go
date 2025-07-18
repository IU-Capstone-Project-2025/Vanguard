package LeaderBoard

import (
	"xxx/LeaderBoardService/Utils"
	"xxx/shared"
)

func (l *LeaderBoard) ComputeLeaderBoard(ans shared.SessionAnswers) (shared.ScoreTable, error) {
	SessionAnswers := ans.Answers
	var CurrentPoints []shared.UserCurrentPoint
	BestTime := Utils.GetEarliestTimestamp(SessionAnswers)
	//WorstTime := Utils.GetLatestTimestamp(SessionAnswers)
	//duration := WorstTime.Sub(BestTime).Seconds()
	MaxScore := 1000
	for _, u := range SessionAnswers {
		if u.Correct {
			elapsed := u.Timestamp.Sub(BestTime).Seconds()
			if elapsed <= 0 {
				elapsed = 0
			}
			timePenalty := float64(elapsed) / 20
			UserPoint := int(float64(MaxScore) * (1 - timePenalty))
			if UserPoint >= MaxScore {
				UserPoint = MaxScore
			}
			CurrentPoints = append(CurrentPoints, shared.UserCurrentPoint{UserId: u.UserId, Score: UserPoint})
		} else {
			UserPoint := 0
			CurrentPoints = append(CurrentPoints, shared.UserCurrentPoint{UserId: u.UserId, Score: UserPoint})
		}
	}
	err := l.Cache.AddScoresBatch(ans.SessionCode, CurrentPoints)
	if err != nil {
		return shared.ScoreTable{}, err
	}
	UserScores, err := l.Cache.LoadLeaderboard(ans.SessionCode)
	SortedUserScore := Utils.SortUserScoresByScoreDesc(UserScores)
	if err != nil {
		return shared.ScoreTable{}, err
	}
	table := shared.ScoreTable{
		SessionCode: ans.SessionCode,
		Users:       SortedUserScore,
	}
	return table, nil
}
