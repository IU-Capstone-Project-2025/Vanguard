package LeaderBoard

import (
	"fmt"
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
			if timePenalty > 1 {
				timePenalty = 1
			}
			UserPoint := int(float64(MaxScore) * (1 - timePenalty))
			if UserPoint <= 0 {
				UserPoint = 0
			}
			if UserPoint >= MaxScore {
				UserPoint = MaxScore
			}
			CurrentPoints = append(CurrentPoints, shared.UserCurrentPoint{UserId: u.UserId, Score: UserPoint})
		} else {
			UserPoint := 0
			CurrentPoints = append(CurrentPoints, shared.UserCurrentPoint{UserId: u.UserId, Score: UserPoint})
		}
	}
	oldScores, err := l.Cache.LoadLeaderboard(ans.SessionCode, nil)
	if err != nil {
		return shared.ScoreTable{}, fmt.Errorf("error to load data from redis")
	}
	prevPlaces := make(map[string]int)
	for _, s := range oldScores {
		prevPlaces[s.UserId] = s.Place
	}
	err = l.Cache.AddScoresBatch(ans.SessionCode, CurrentPoints)
	if err != nil {
		return shared.ScoreTable{}, fmt.Errorf("error to add points in redis")
	}
	newScores, err := l.Cache.LoadLeaderboard(ans.SessionCode, prevPlaces)
	if err != nil {
		return shared.ScoreTable{}, fmt.Errorf("error to load data from redis redis")
	}
	table := shared.ScoreTable{
		SessionCode: ans.SessionCode,
		Users:       newScores,
	}
	return table, nil
}
