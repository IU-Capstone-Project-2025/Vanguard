package models

import (
	"time"
	"xxx/shared"
)

// OngoingQuiz stores data of the quiz process: Quiz payload, index of the current question
type OngoingQuiz struct {
	CurrQuestionIdx int         // index of the current question
	QuizData        shared.Quiz // the questions and options of the quiz
}

// UserAnswer stores the information about the answer given by a user: its correctness and timestamp, when answer was arrived
type UserAnswer struct {
	Answered  bool      `json:"answered"`  // indicates if user even have answered; if it is false, other fields are not matter
	Correct   bool      `json:"correct"`   // correctness of user's answer
	Timestamp time.Time `json:"timestamp"` // time when user has answered
}
