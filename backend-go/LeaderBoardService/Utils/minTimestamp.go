package Utils

import (
	"time"
	"xxx/shared"
)

func GetEarliestTimestamp(answers []shared.UserAnswer) time.Time {
	earliest := answers[0].Timestamp

	for _, ans := range answers[1:] {
		if ans.Timestamp.Before(earliest) {
			earliest = ans.Timestamp
		}
	}

	return earliest
}

func GetLatestTimestamp(answers []shared.UserAnswer) time.Time {
	earliest := answers[0].Timestamp

	for _, ans := range answers[1:] {
		if ans.Timestamp.After(earliest) {
			earliest = ans.Timestamp
		}
	}

	return earliest
}
