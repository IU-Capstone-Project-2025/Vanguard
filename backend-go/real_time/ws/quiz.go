package ws

import (
	"sync"
	"xxx/real_time/models"
	"xxx/shared"
)

// QuizTracker tracks the current quiz for each session in the map: sessionId -> models.QuizGame.
// The tracker is thread-safe
type QuizTracker struct {
	mu      sync.RWMutex
	Tracker map[string]models.QuizGame
}

func NewQuestionTracker() *QuizTracker {
	qt := &QuizTracker{
		mu:      sync.RWMutex{},
		Tracker: make(map[string]models.QuizGame),
	}

	// restore map from Redis if the service was down
	qt.restoreData()

	return qt
}

// GetCurrentQuestionIdx method returns the current question index of the session [sessionId]
func (q *QuizTracker) GetCurrentQuestionIdx(sessionId string) (int, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if quiz, exists := q.Tracker[sessionId]; !exists {
		return -1, false
	} else {
		return quiz.CurrQuestionIdx, true
	}
}

// SetCurrQuestionIdx method assigns the given [questionIdx] to the session [sessionId]
func (q *QuizTracker) SetCurrQuestionIdx(sessionId string, questionIdx int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	quiz := q.Tracker[sessionId]
	quiz.CurrQuestionIdx = questionIdx

	q.Tracker[sessionId] = quiz
}

// IncQuestionIdx method increments the current question index of the session [sessionId]
func (q *QuizTracker) IncQuestionIdx(sessionId string) bool {
	idx, err := q.GetCurrentQuestionIdx(sessionId)
	if !err {
		return false
	}

	q.SetCurrQuestionIdx(sessionId, idx+1)
	return true
}

// GetCorrectOptionIdx returns the index and the object of the correct answer for the given question
func (q *QuizTracker) GetCorrectOptionIdx(sessionId string, questionIdx int) (int, shared.Option) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if quiz, exists := q.Tracker[sessionId]; !exists {
		return -1, shared.Option{}
	} else {
		question := quiz.QuizData.GetQuestion(questionIdx)
		return question.GetCorrectOption()
	}
}

// restoreData restores map data from the Redis
func (q *QuizTracker) restoreData() {
	// TODO: implement data restoring
}
