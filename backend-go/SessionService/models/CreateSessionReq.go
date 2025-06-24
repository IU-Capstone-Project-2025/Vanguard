package models

type CreateSessionReq struct {
	UserId string `json:"userId"`
	QuizId string `json:"quizId"`
}
