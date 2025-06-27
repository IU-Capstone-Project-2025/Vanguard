package shared

type Option struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type Question struct {
	Type    string   `json:"type"`
	Text    string   `json:"text"`
	Options []Option `json:"options"`
}

type Quiz struct {
	Questions []Question `json:"questions"`
}

type QuizMessage struct {
	SessionId string `json:"session_id"`
	Quiz      Quiz   `json:"quiz"`
}
