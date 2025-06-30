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

func (q Question) IsCorrectOption() {

}

func (q Question) GetCorrectOption() (int, Option) {
	for i, op := range q.Options {
		if op.IsCorrect {
			return i, op
		}
	}
	return 0, Option{}
}

type Quiz struct {
	Questions []Question `json:"questions"`
}

func (q Quiz) GetQuestion(idx int) Question {
	if idx < 0 || idx >= len(q.Questions) {
		return Question{}
	}
	return q.Questions[idx]
}

func (q Quiz) Len() int {
	return len(q.Questions)
}

type QuizMessage struct {
	SessionId string `json:"session_id"`
	Quiz      Quiz   `json:"quiz"`
}
