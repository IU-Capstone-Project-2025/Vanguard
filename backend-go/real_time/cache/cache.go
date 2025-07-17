package cache

import "xxx/real_time/models"

type Cache interface {
	SetSessionQuiz(sessionId string, quizData models.OngoingQuiz) error
	GetSessionQuiz(sessionId string) (models.OngoingQuiz, error)
	DeleteSession(sessionId string) error
	GetAllSessions() (map[string]models.OngoingQuiz, error)

	SetQuestionIndex(sessionId string, questionIdx int) error
	GetQuestionIndex(sessionId string) (int, error)

	RecordAnswer(sessionID, userID string, question int, answer models.UserAnswer) error
	GetAllAnswers(sessionId string) (map[string][]models.UserAnswer, error)
}
