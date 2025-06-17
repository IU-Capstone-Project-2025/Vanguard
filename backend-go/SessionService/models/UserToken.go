package models

type UserToken struct {
	UserId      string `json:"userId"`
	CurrentQuiz string `json:"currentQuiz"`
	Exp         int64  `json:"exp"`
}
