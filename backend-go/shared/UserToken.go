package shared

type UserToken struct {
	UserId           string `json:"userId"`
	UserType         string `json:"userType"`
	ServerWsEndpoint string `json:"serverWsEndpoint"`
	CurrentQuiz      string `json:"currentQuiz"`
	Exp              int64  `json:"exp"`
}
