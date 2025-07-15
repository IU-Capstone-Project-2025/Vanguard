package LeaderBoard

import (
	"xxx/LeaderBoardService/models"
	"xxx/shared"
)

func (l *LeaderBoard) PopularAns(ans shared.SessionAnswers) (models.PopularAns, error) {
	answers := models.PopularAns{
		SessionCode: ans.SessionCode,
		Answers:     make(map[string]int),
	}

	UserAns := ans.Answers
	for _, UserAn := range UserAns {
		if !UserAn.Answered {
			continue
		}
		answers.Answers[UserAn.Option] += 1
	}
	return answers, nil
}
