package shared

import "time"

type Answer struct {
	UserId    string    `json:"user_id"`
	Correct   bool      `json:"correct"` // correctness of user's answer
	Answered  bool      `json:"answered"`
	Option    string    `json:"option"`
	Timestamp time.Time `json:"timestamp"` // time when user has answered
}

type SessionAnswers struct {
	SessionCode   string   `json:"session_code"`
	OptionsAmount int      `json:"options_amount"`
	Answers       []Answer `json:"answers"`
}
