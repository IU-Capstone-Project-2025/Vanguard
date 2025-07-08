package models

type CreateSessionReq struct {
	UserName string `json:"userName"`
	QuizId   string `json:"quizId"`
}
