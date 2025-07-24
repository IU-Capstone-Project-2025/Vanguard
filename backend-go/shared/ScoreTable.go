package shared

const (
	ProgressUp   = "Up"
	ProgressDown = "Down"
	ProgressSame = "Same"
)

type UserScore struct {
	UserId        string `json:"user_id"`
	Place         int    `json:"place"`
	PreviousPlace int    `json:"previous_place"`
	Progress      string `json:"progress"`
	TotalScore    int    `json:"total_score"`
}

type UserCurrentPoint struct {
	UserId string `json:"user_id"`
	Score  int    `json:"score"`
}

type ScoreTable struct {
	SessionCode string      `json:"session_code"`
	Users       []UserScore `json:"users"`
}

type PopularAns struct {
	SessionCode string         `json:"session_code"`
	Answers     map[string]int `json:"answers"`
}

type BoardResponse struct {
	SessionCode string     `json:"session_code"`
	Table       ScoreTable `json:"table"`
	Popular     PopularAns `json:"popular"`
}
