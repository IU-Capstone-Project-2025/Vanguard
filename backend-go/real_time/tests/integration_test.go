//go:build integration

package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
	"xxx/integration_tests/utils"
	"xxx/real_time/ws"
	"xxx/shared"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestWithTestContainers(t *testing.T) {
	ctx := context.Background()

	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		fmt.Println("LOADING .ENV")
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			log.Fatalf("Error: could not load .env file: %v", err)
		}
	}

	amqpURL, closeRabbit := utils.StartRabbit(ctx, t)
	defer closeRabbit()
	redisUrl, closeRedis := utils.StartRedis(ctx, t)
	defer closeRedis()

	t.Log("rabbit running on ", amqpURL)
	t.Log("redis running on ", redisUrl)

	// 2. Start Redis container similarly if needed

	// 3. Start RealTime service in a goroutine or exec.Command, configuring it to connect to amqpURL and Redis.
	//    For brevity, assume RealTime service can be started in-process or as a subprocess, reading env vars:
	// Start your RealTime main in a goroutine if possible, or exec binary.
	wgRTS := &sync.WaitGroup{}
	wgRTS.Add(1)
	cancel := utils.StartRealTimeServer(t, wgRTS, amqpURL, redisUrl)

	wg := &sync.WaitGroup{}

	sessionIds := []string{"В4ФЛ3Р", "CO6AK4"}

	for _, s := range sessionIds {
		wg.Add(1)
		go func(sessionId string) {
			defer wg.Done()
			adminId := "admin"
			users := []string{"ginger"}
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

			usersAnswers := [][]int{
				{2, 2, 3},
			}
			adminToken := generateJWT(t, sessionId, adminId, shared.RoleAdmin)

			var adminConn *websocket.Conn
			var usersConn []*websocket.Conn

			// 5. Publish session.start
			publishSessionStart(t, amqpURL, sessionId, quiz)

			// ========================================================================
			// START OF THE DUMMY FIX
			// ========================================================================
			//
			// Give the server a moment to consume the 'session.start' message and
			// create the dynamic consumer for 'question.<session_id>.start'.
			// This value may need to be increased if the CI runner is slow.
			t.Log("Waiting for 5 seconds for the server to set up session consumers...")
			time.Sleep(5 * time.Second)
			//
			// ========================================================================
			// END OF THE DUMMY FIX
			// ========================================================================

			// 6. Admin WS connection
			adminConn = utils.ConnectWs(t, adminToken)

			// 7. Read welcome
			utils.ReadWs(t, adminConn)

			// join users
			for _, userId := range users {
				userToken := generateJWT(t, sessionId, userId, shared.RoleParticipant)
				conn := utils.ConnectWs(t, userToken)
				usersConn = append(usersConn, conn)

				utils.ReadWs(t, conn)
			}

			// 8. Start question flow
			for i, q := range quiz.Questions {
				t.Log("trigger question ", i, q)

				publishQuestionStart(t, amqpURL, sessionId)

				questionPayload := utils.ReadWs(t, adminConn)
				t.Log("checking question payload:")
				require.Equal(t, q.Text, questionPayload.Text)
				require.Equal(t, ws.MessageTypeQuestion, questionPayload.Type)
				require.Equal(t, i+1, questionPayload.QuestionIdx)
				require.Equal(t, q.Options, questionPayload.Options)

				for _, user := range usersConn {
					utils.ReadWs(t, user)
				}

				for j, user := range usersConn {
					option := usersAnswers[j][i]
					t.Logf("user %s send answer: %d", users[j], option)
					msg := ws.ClientMessage{Option: option}
					user.WriteMessage(websocket.TextMessage, msg.Bytes())

					resp := utils.ReadWs(t, user)
					require.Equal(t, ws.MessageTypeAnswer, resp.Type)
					require.Equal(t, i+1, resp.QuestionIdx)
					require.Equal(t, q.Options[option].IsCorrect, resp.Correct)
				}
			}

			// trigger session end
			publishSessionEnd(t, amqpURL, sessionId)
			t.Log("---- Admin received end message:")
			// readWs(t, adminConn)

			t.Log("---- Users received end message:")
			for _, user := range usersConn {
				utils.ReadWs(t, user)
				// t.Log(lb.Payload)
				// ans, ok := lb.Payload.(map[string]interface{})
				// require.Equal(t, true, ok)

				// userChosen, ok := ans[users[i]].([]interface{})
				// require.Equal(t, true, ok)

				// for j, isCorrectInter := range userChosen {
				// 	chosenIdx := usersAnswers[i][j]

				// 	isCorrect, ok := isCorrectInter.(bool)
				// 	require.Equal(t, true, ok)

				// 	require.Equal(t, quiz.Questions[j].Options[chosenIdx].IsCorrect, isCorrect)
				// }
			}

			t.Log("Close connections:")
			// Ensure all connections will be closed at end
			utils.CloseWs(adminConn)
			t.Log("Closed admin")
			for _, c := range usersConn {
				utils.CloseWs(c)
			}
		}(s)
	}

	wg.Wait()
	cancel()
	wgRTS.Wait()
}

func publishSessionStart(t *testing.T, amqpURL, sessionId string, quiz shared.Quiz) {
	rabCon, err := amqp.Dial(amqpURL)
	if err != nil {
		t.Fatalf("Dial RabbitMQ: %v", err)
	}
	ch, err := rabCon.Channel()
	if err != nil {
		t.Fatalf("Open channel: %v", err)
	}
	evt := shared.QuizMessage{
		SessionId: sessionId,
		Quiz:      quiz,
	}
	body, _ := json.Marshal(evt)
	ch.Publish(shared.SessionExchange, "session.start", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	rabCon.Close()
}

func publishSessionEnd(t *testing.T, amqpURL, sessionId string) {
	rabCon, err := amqp.Dial(amqpURL)
	if err != nil {
		t.Fatalf("Dial RabbitMQ: %v", err)
	}
	ch, err := rabCon.Channel()
	if err != nil {
		t.Fatalf("Open channel: %v", err)
	}
	body, _ := json.Marshal(sessionId)
	ch.Publish(shared.SessionExchange, "session.end", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	rabCon.Close()
}

func publishQuestionStart(t *testing.T, amqpURL, sessionId string) {
	rabCon, err := amqp.Dial(amqpURL)
	if err != nil {
		t.Fatalf("Dial RabbitMQ: %v", err)
	}
	ch, err := rabCon.Channel()
	if err != nil {
		t.Fatalf("Open channel: %v", err)
	}
	ch.Publish(shared.SessionExchange, fmt.Sprintf("question.%s.start", sessionId),
		false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        nil,
		})
	rabCon.Close()
}

func generateJWT(t *testing.T, session, user string, role shared.UserRole) string {
	claims := shared.UserToken{
		UserId:    user,
		UserType:  role,
		SessionId: session,
		Exp:       10000000,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Unix(0, 10000000)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rawJwt, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	require.NoError(t, err)

	return rawJwt
}

func getEnvFilePath() string {
	root, err := filepath.Abs("../../..")
	if err != nil {
		log.Fatal("failed to find project root dir")
	}
	return filepath.Join(root, ".env")
}
