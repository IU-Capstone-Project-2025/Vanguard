package LeaderBoard

import (
	"xxx/shared"
)

func (l *LeaderBoard) PopularAns(ans shared.SessionAnswers) (shared.PopularAns, error) {
	answers := shared.PopularAns{
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
