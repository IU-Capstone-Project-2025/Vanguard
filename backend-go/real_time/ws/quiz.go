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
	answers map[string]map[string][]bool // sessionId -> userId -> correctness
	tracker map[string]models.QuizGame   // stores the whole quiz data for each session.
	// Includes the index of current question and all questions with answer options.
}

func NewQuizTracker() *QuizTracker {
	qt := &QuizTracker{
		mu:      sync.RWMutex{},
		answers: make(map[string]map[string][]bool),
		tracker: make(map[string]models.QuizGame),
	}

	// restore map from Redis if the service was down
	qt.restoreData()

	return qt
}

// GetCurrentQuestion method returns the current question index of the session [sessionId] and the payload of the question
func (q *QuizTracker) GetCurrentQuestion(sessionId string) (int, *shared.Question) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if quiz, exists := q.tracker[sessionId]; !exists {
		return -1, nil
	} else {
		question := quiz.QuizData.GetQuestion(quiz.CurrQuestionIdx)
		return quiz.CurrQuestionIdx, &question
	}
}

// SetCurrQuestionIdx method assigns the given [questionIdx] to the session [sessionId]
func (q *QuizTracker) SetCurrQuestionIdx(sessionId string, questionIdx int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	quiz := q.tracker[sessionId]
	quiz.CurrQuestionIdx = questionIdx

	q.tracker[sessionId] = quiz
}

// IncQuestionIdx method increments the current question index of the session [sessionId]
func (q *QuizTracker) IncQuestionIdx(sessionId string) bool {
	idx, question := q.GetCurrentQuestion(sessionId)
	if question == nil {
		return false
	}

	q.SetCurrQuestionIdx(sessionId, idx+1)
	return true
}

// GetCorrectOption returns the index and the object of the correct answer for the given question
func (q *QuizTracker) GetCorrectOption(sessionId string, questionIdx int) (int, *shared.Option) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if quiz, exists := q.tracker[sessionId]; !exists {
		return -1, nil
	} else {
		question := quiz.QuizData.GetQuestion(questionIdx)
		idx, op := question.GetCorrectOption()
		return idx, &op
	}
}

// RecordAnswer stores whether a user’s answer was correct.
func (q *QuizTracker) RecordAnswer(sessionId, userId string, correct bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.answers[sessionId][userId]; !ok {
		q.answers[sessionId][userId] = make([]bool, 0) // create array with length = the amount of questions
	}
	q.answers[sessionId][userId] = append(q.answers[sessionId][userId], correct)
}

// NewSession adds new session and links corresponding quiz object to it
func (q *QuizTracker) NewSession(sessionId string, quizData shared.Quiz) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.tracker[sessionId]; !exists {
		q.tracker[sessionId] = models.QuizGame{
			CurrQuestionIdx: -1,
			QuizData:        quizData,
		}
		q.answers[sessionId] = make(map[string][]bool)
	}
}

// GetLeaderboard returns a simple map of userId -> correctFlag
func (q *QuizTracker) GetLeaderboard(sessionId string) map[string][]bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	// copy to avoid races
	copyMap := make(map[string][]bool, len(q.answers[sessionId]))
	for user, answers := range q.answers[sessionId] {
		copyMap[user] = answers
	}
	return copyMap
}

// restoreData restores map data from the Redis
func (q *QuizTracker) restoreData() {
	// TODO: implement data restoring
}
