package ws

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"xxx/real_time/cache"
	"xxx/real_time/cache/redis"
	"xxx/real_time/leaderboard"
	"xxx/real_time/models"
	"xxx/shared"
)

// QuizTracker tracks the current quiz for each session in the map: sessionId -> models.OngoingQuiz.
// The tracker is thread-safe
type QuizTracker struct {
	mu      sync.Mutex
	answers map[string]map[string][]models.UserAnswer // sessionId -> userId -> [models.UserAnswer]
	tracker map[string]models.OngoingQuiz             // stores the whole quiz data for each session.
	// Includes the index of current question and all questions with answer options.
	cache cache.Cache // cache (e.g. Redis storage manager) to store copy of states from quiz tracker
	lb    *leaderboard.Client
}

func NewQuizTracker(leaderboardUrl string) *QuizTracker {
	qt := &QuizTracker{
		mu:      sync.Mutex{},
		answers: make(map[string]map[string][]models.UserAnswer),
		tracker: make(map[string]models.OngoingQuiz),
		cache:   &redis.Client{},
		lb:      leaderboard.NewClient(leaderboardUrl),
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
	q.mu.Lock()
	defer q.mu.Unlock()

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
	quiz, ok := q.tracker[sessionId]
	if !ok {
		return false
	}
	if quiz.CurrQuestionIdx+1 >= quiz.QuizData.Len() {
		return false
	}

	quiz.CurrQuestionIdx++
	q.tracker[sessionId] = quiz
	_ = q.cache.SetQuestionIndex(sessionId, quiz.CurrQuestionIdx)
	return true
}

// GetCorrectOption returns the index and the object of the correct answer for the given question
func (q *QuizTracker) GetCorrectOption(sessionId string, questionIdx int) (int, *shared.Option) {
	q.mu.Lock()
	defer q.mu.Unlock()

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

// AddParticipant initializes new user's answers array with default values.
// []models.UserAnswer array must be initialized, since user can leave question without answer recording,
// and then it will be marked just an 'not answered'
func (q *QuizTracker) AddParticipant(sessionId, userId string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.answers[sessionId][userId]; !ok {
		q.answers[sessionId][userId] = make([]models.UserAnswer, q.tracker[sessionId].QuizData.Len()) // create array with length = the amount of questions
	}
}

// GetAnswers returns the correctness of all answers given by users
func (q *QuizTracker) GetAnswers(sessionId string) map[string][]models.UserAnswer {
	q.mu.Lock()
	defer q.mu.Unlock()

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
func (q *QuizTracker) GetLeaderboard(sessionId string) (shared.BoardResponse, error) {
	q.mu.Lock()

	qid := q.tracker[sessionId].CurrQuestionIdx
	if qid > 0 {
		qid-- // Leader board is counted for previous question
	}

	q.mu.Unlock()

	currQuestionAnswers := make([]shared.Answer, 0, len(q.answers[sessionId]))
	for user, answers := range q.answers[sessionId] {
		ans := answers[qid]

		lbAns := shared.Answer{
			UserId:    user,
			Correct:   ans.Correct,
			Answered:  ans.Answered,
			Option:    strconv.Itoa(ans.Option + 1), // 1-based option index
			Timestamp: ans.Timestamp,
		}
		currQuestionAnswers = append(currQuestionAnswers, lbAns)
	}

	fmt.Println("currQuestionAnswers: ", currQuestionAnswers)

	_, question := q.GetCurrentQuestion(sessionId)

	board, err := q.lb.GetResults(context.Background(), sessionId, currQuestionAnswers, len(question.Options))
	if err != nil {
		return shared.BoardResponse{}, err
	}

	return board, nil
}

func (q *QuizTracker) GetQuizLen(sessionId string) int {
	q.mu.Lock()
	defer q.mu.Unlock()

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
