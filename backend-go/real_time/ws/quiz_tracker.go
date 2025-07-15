package ws

import (
	"fmt"
	"sync"
	"xxx/real_time/cache"
	"xxx/real_time/cache/redis"
	"xxx/real_time/models"
	"xxx/shared"
)

// QuizTracker tracks the current quiz for each session in the map: sessionId -> models.OngoingQuiz.
// The tracker is thread-safe
type QuizTracker struct {
	mu      sync.RWMutex
	answers map[string]map[string][]models.UserAnswer // sessionId -> userId -> [1st question correctness, 2nd, etc.]
	tracker map[string]models.OngoingQuiz             // stores the whole quiz data for each session.
	// Includes the index of current question and all questions with answer options.
	cache cache.Cache // cache (e.g. Redis storage manager) to store copy of states from quiz tracker
}

func NewQuizTracker() *QuizTracker {
	qt := &QuizTracker{
		mu:      sync.RWMutex{},
		answers: make(map[string]map[string][]models.UserAnswer),
		tracker: make(map[string]models.OngoingQuiz),
		cache:   &redis.Client{},
	}

	return qt
}

// SetCache sets cache field assigning the given one
func (q *QuizTracker) SetCache(cache cache.Cache) {
	q.cache = cache
	// restore map from Redis if the service was down
	q.restoreData()
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

	err := q.cache.SetQuestionIndex(sessionId, questionIdx)
	fmt.Println("Redis err: ", err)
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
func (q *QuizTracker) RecordAnswer(sessionId, userId string, answer models.UserAnswer) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.answers[sessionId][userId]; !ok {
		q.answers[sessionId][userId] = make([]models.UserAnswer, q.tracker[sessionId].QuizData.Len()) // create array with length = the amount of questions
	}

	qid := q.tracker[sessionId].CurrQuestionIdx
	q.answers[sessionId][userId][qid] = answer
	q.cache.RecordAnswer(sessionId, userId, qid, answer)
}

// GetAnswers returns the correctness of all answers given by users
func (q *QuizTracker) GetAnswers(sessionId string) map[string][]models.UserAnswer {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.answers[sessionId]
}

// NewSession adds new session and links corresponding quiz object to it
func (q *QuizTracker) NewSession(sessionId string, quizData shared.Quiz) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.tracker[sessionId]; !exists {
		q.tracker[sessionId] = models.OngoingQuiz{
			CurrQuestionIdx: -1, // before starting the first question (0-th index), the index is -1
			QuizData:        quizData,
		}
		q.answers[sessionId] = make(map[string][]models.UserAnswer)
		q.cache.SetSessionQuiz(sessionId, q.tracker[sessionId])
	}
}

// DeleteSession deletes session from tracker
func (q *QuizTracker) DeleteSession(sessionId string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.tracker[sessionId]; exists {
		delete(q.answers, sessionId)
		delete(q.tracker, sessionId)

		q.cache.DeleteSession(sessionId)
	}
}

// GetLeaderboard returns a simple map of userId -> correctFlag
func (q *QuizTracker) GetLeaderboard(sessionId string) map[string][]models.UserAnswer {
	q.mu.RLock()
	defer q.mu.RUnlock()
	// copy to avoid races
	copyMap := make(map[string][]models.UserAnswer, len(q.answers[sessionId]))
	for user, answers := range q.answers[sessionId] {
		copyMap[user] = answers
	}
	return copyMap
}

func (q *QuizTracker) GetQuizLen(sessionId string) int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.tracker[sessionId].QuizData.Len()
}

// restoreData restores map data from the Redis
func (q *QuizTracker) restoreData() {
	quizzes, err := q.cache.GetAllSessions()
	if err != nil {
		fmt.Println("failed to restore data from Redis: ", err)
	}

	q.tracker = quizzes
}
