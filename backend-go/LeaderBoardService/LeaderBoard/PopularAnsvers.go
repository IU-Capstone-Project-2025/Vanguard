package LeaderBoard

import (
	"strconv"
	"xxx/shared"
)

func (l *LeaderBoard) PopularAns(ans shared.SessionAnswers) (shared.PopularAns, error) {
	answers := shared.PopularAns{
		SessionCode: ans.SessionCode,
		Answers:     make(map[string]int),
	}
	for i := 1; i < ans.OptionsAmount+1; i++ {
		answers.Answers[strconv.Itoa(i)] = 0
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
