package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"xxx/shared"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"xxx/real_time/models"
)

// setupRedisContainer starts a Redis container and returns its address and a termination function.
func setupRedisContainer(tst *testing.T) (addr string, terminate func()) {
	t := tst
	t.Helper()
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := redisC.Host(ctx)
	require.NoError(t, err)
	port, err := redisC.MappedPort(ctx, "6379/tcp")
	require.NoError(t, err)

	addr = fmt.Sprintf("%s:%s", host, port.Port())
	terminate = func() {
		err := redisC.Terminate(ctx)
		require.NoError(t, err)
	}
	return addr, terminate
}

func TestRedisClientIntegration(t *testing.T) {
	addr, terminate := setupRedisContainer(t)
	defer terminate()

	client := NewClient(addr, "", 0)

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

	err := client.SetSessionQuiz("sess1", ongoing)
	require.NoError(t, err)

	gotQuiz, err := client.GetSessionQuiz("sess1")
	require.NoError(t, err)
	require.Equal(t, ongoing.CurrQuestionIdx, gotQuiz.CurrQuestionIdx)
	require.Equal(t, quiz.Questions, gotQuiz.QuizData.Questions)

	// simulate TTL expiration by adjusting TTL nearly expired and sleeping
	err = client.rdb.Expire(client.ctx, "session:sess1:quiz_state", 1*time.Second).Err()
	require.NoError(t, err)

	time.Sleep(1100 * time.Millisecond)
	_, err = client.GetSessionQuiz("sess1")
	require.Error(t, err)

	// 2. Test SetQuestionIndex & GetQuestionIndex
	initial := models.OngoingQuiz{CurrQuestionIdx: 0}
	err = client.SetSessionQuiz("sess2", initial)
	require.NoError(t, err)

	err = client.SetQuestionIndex("sess2", 5)
	require.NoError(t, err)

	idx, err := client.GetQuestionIndex("sess2")
	require.NoError(t, err)
	require.Equal(t, 5, idx)

	// 3. Test RecordAnswer & direct HGETALL
	err = client.RecordAnswer("sess3", "user1", 0, models.UserAnswer{
		Correct:   true,
		Timestamp: time.Time{},
	})
	require.NoError(t, err)

	// HGetAll for this user
	hashKey := fmt.Sprintf("session:%s:user:%s:answers", "sess3", "user1")
	hashData, err := client.rdb.HGetAll(client.ctx, hashKey).Result()
	require.NoError(t, err)
	require.Contains(t, hashData, "0")

	var ua models.UserAnswer
	require.NoError(t, json.Unmarshal([]byte(hashData["0"]), &ua))
	require.True(t, ua.Correct)

	// 4. Test DeleteSession
	err = client.DeleteSession("sess3")
	require.NoError(t, err)

	// Keys should be gone
	_, err = client.rdb.Get(client.ctx, "session:sess3:quiz_state").Result()
	require.Error(t, err)
	exists, err := client.rdb.Exists(client.ctx, hashKey).Result()
	require.NoError(t, err)
	require.Zero(t, exists)

	// 5. Test GetAllAnswers on fresh session
	all, err := client.GetAllAnswers("sess4")
	require.NoError(t, err)
	require.Empty(t, all)
}
