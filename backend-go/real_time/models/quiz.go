package models

import "xxx/shared"

// QuizGame stores data of the quiz process: Quiz payload, index of the current question
type QuizGame struct {
	CurrQuestionIdx int         // index of the current question
	QuizData        shared.Quiz // the questions and options of the quiz
}
