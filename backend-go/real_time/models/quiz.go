package models

import "xxx/shared"

type QuizGame struct {
	CurrQuestionIdx int         // index of the current question
	QuizData        shared.Quiz // the questions and options of the quiz
}
