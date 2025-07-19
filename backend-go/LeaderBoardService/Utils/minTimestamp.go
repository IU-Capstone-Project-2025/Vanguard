package Utils

import (
	"time"
	"xxx/shared"
)

func GetEarliestTimestamp(answers1 []shared.Answer) time.Time {
	var answers []shared.Answer
	for _, u := range answers1 {
		if u.Correct {
			answers = append(answers, u)
		}
	}
	if answers == nil || len(answers) == 0 {
		return time.Now()
	}
	earliest := answers[0].Timestamp

	for _, ans := range answers[1:] {
		if ans.Timestamp.Before(earliest) {
			earliest = ans.Timestamp
		}
	}

	return earliest
}

func GetLatestTimestamp(answers []shared.Answer) time.Time {
	earliest := answers[0].Timestamp

	for _, ans := range answers[1:] {
		if ans.Timestamp.After(earliest) {
			earliest = ans.Timestamp
		}
	}

	return earliest
}
