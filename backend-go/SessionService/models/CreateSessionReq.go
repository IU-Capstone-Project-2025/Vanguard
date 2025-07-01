package models

type CreateSessionReq struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	QuizId   string `json:"quizId"`
}
