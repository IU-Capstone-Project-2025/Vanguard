package redis_test

import (
	"context"
	"testing"
	"time"
	"xxx/integration_tests/utils"
	"xxx/real_time/cache/redis"
	"xxx/shared"

	"github.com/stretchr/testify/require"
	"xxx/real_time/models"
)

func TestRedisClientIntegration(t *testing.T) {
	ctx := context.Background()
	adddr, terminate := utils.StartRedis(ctx, t)
	defer terminate()

	client := redis.NewClient(adddr, "", 0)

	// 1. Test SetSessionQuiz & GetSessionQuiz with TTL
	quiz := shared.Quiz{Questions: []shared.Question{
		{
			Type: "single_choice",
			Text: "What is the output of print(2 ** 3)?",
			Options: []shared.Option{
				{Text: "6", IsCorrect: false},
				{Text: "8", IsCorrect: true},
				{Text: "9", IsCorrect: false},
				{Text: "5", IsCorrect: false},
			},
		},
		{
			Type: "single_choice",
			Text: "Which keyword is used to create a function in Python?",
			Options: []shared.Option{
				{Text: "func", IsCorrect: false},
				{Text: "function", IsCorrect: false},
				{Text: "def", IsCorrect: true},
				{Text: "define", IsCorrect: false},
			},
		},
		{
			Type: "single_choice",
			Text: "What data type is the result of: 3 / 2 in Python 3?",
			Options: []shared.Option{
				{Text: "int", IsCorrect: false},
				{Text: "float", IsCorrect: true},
				{Text: "str", IsCorrect: false},
				{Text: "decimal", IsCorrect: false},
			},
		},
	}}
	ongoing := models.OngoingQuiz{
		CurrQuestionIdx: 0,
		QuizData:        quiz,
	}

	t.Log("---- Testing: Set and Get session quiz:")
	t.Log("Set: ", ongoing)
	err := client.SetSessionQuiz("sess1", ongoing)
	require.NoError(t, err)

	gotQuiz, err := client.GetSessionQuiz("sess1")
	require.NoError(t, err)
	t.Log("Got: ", gotQuiz)
	require.Equal(t, ongoing.CurrQuestionIdx, gotQuiz.CurrQuestionIdx)
	require.Equal(t, quiz.Questions, gotQuiz.QuizData.Questions)

	t.Log("Simulate TTL expiration by adjusting TTL nearly expired and sleeping")
	err = client.Expire("session:sess1:quiz_state", 1*time.Second)
	require.NoError(t, err)

	time.Sleep(1100 * time.Millisecond)
	_, err = client.GetSessionQuiz("sess1")
	require.Error(t, err)
	t.Log("TTL expired successfully")

	t.Log("---- Testing: Set and Get question index for sessionID = ", "sess2")
	// 2. Test SetQuestionIndex & GetQuestionIndex
	initial := models.OngoingQuiz{CurrQuestionIdx: 2, QuizData: quiz}
	err = client.SetSessionQuiz("sess2", initial)
	require.NoError(t, err)

	err = client.SetSessionQuiz("sess3", initial)
	require.NoError(t, err)

	quizzes, err := client.GetAllSessions()
	require.NoError(t, err)

	sess2, ok := quizzes["sess2"]
	require.Equal(t, true, ok)
	require.Equal(t, initial.CurrQuestionIdx, sess2.CurrQuestionIdx)
	require.Equal(t, initial.QuizData, sess2.QuizData)

	sess3, ok := quizzes["sess3"]
	require.Equal(t, true, ok)
	require.Equal(t, initial.CurrQuestionIdx, sess3.CurrQuestionIdx)
	require.Equal(t, initial.QuizData, sess3.QuizData)

	t.Log("Set: ", initial, ", index = ", initial.CurrQuestionIdx)
	err = client.SetQuestionIndex("sess2", initial.CurrQuestionIdx)
	require.NoError(t, err)

	idx, err := client.GetQuestionIndex("sess2")
	t.Log("Got: ", idx)
	require.NoError(t, err)
	require.Equal(t, initial.CurrQuestionIdx, idx)

	t.Log("---- Testing: RecordAnswer & direct GetAllAnswers")
	t.Log("Set: ", ongoing)
	err = client.RecordAnswer("sess3", "user1", 0, models.UserAnswer{
		Correct:   true,
		Timestamp: time.Time{},
	})
	require.NoError(t, err)

	t.Log("HGetAll for this user")
	answers, err := client.GetAllAnswers("sess3")
	t.Log("answers of user1 are, ", answers["user1"])
	require.NoError(t, err)
	require.Contains(t, answers, "user1")

	// 4. Test DeleteSession
	err = client.DeleteSession("sess3")
	require.NoError(t, err)

	// Keys should be gone
	sess3quiz, err := client.GetSessionQuiz("sess3")
	require.Error(t, err)
	require.Zero(t, sess3quiz)

	// 5. Test GetAllAnswers on fresh session
	all, err := client.GetAllAnswers("sess4")
	require.NoError(t, err)
	require.Empty(t, all)
}
