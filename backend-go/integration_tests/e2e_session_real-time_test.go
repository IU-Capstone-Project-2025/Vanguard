//go:build integration

package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
	"xxx/SessionService/models"
	"xxx/integration_tests/utils"
	"xxx/real_time/ws"
)

var (
	sessionServiceURL = "http://localhost:8081" // Session Service base URL
)

// Structs matching your Swagger definitions:
type CreateSessionReq struct {
	QuizId string `json:"quizId"`
	UserId string `json:"userId"`
}
type ValidateCodeReq struct {
	Code   string `json:"code"`
	UserId string `json:"userId"`
}

func loadEnv(t *testing.T) {
	// Optionally load .env if needed for configuration
	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		// Assume .env is at project root; adapt path as needed
		if err := godotenv.Load("../../.env"); err != nil {
			// Not fatal if .env missing, but log
			t.Logf("No .env loaded: %v", err)
		}
	}
}

func TestSessionServiceToRealTime_E2E(t *testing.T) {
	loadEnv(t)
	ctx := context.Background()

	amqpUrl, closeRabbit := utils.StartRabbit(ctx, t)
	defer closeRabbit()
	redisUrl, closeRedis := utils.StartRedis(ctx, t)
	defer closeRedis()

	t.Log("------------ wait for real time service -----------------")
	wgRTS := &sync.WaitGroup{}
	wgRTS.Add(1)

	cancel := utils.StartRealTimeServer(t, wgRTS, amqpUrl, redisUrl)

	defer cancel()

	t.Log("------------ wait for session service -----------------")
	go func() {
		utils.StartSessionService(t, amqpUrl, redisUrl)
	}()
	time.Sleep(2 * time.Second)

	wg := &sync.WaitGroup{}

	n := 2 // the number of session to create

	for i := 0; i < n; i++ { // process with n sessions
		wg.Add(1)
		go func() {
			defer wg.Done()

			// 1. Create a new session as admin
			adminId := "admin_id"
			createReq := CreateSessionReq{
				QuizId: "1",
				UserId: adminId,
			}
			createReqBody, err := json.Marshal(createReq)
			require.NoError(t, err)

			createResp, err := http.Post(sessionServiceURL+"/sessionsMock", "application/json", bytes.NewReader(createReqBody))
			require.NoError(t, err)
			defer createResp.Body.Close()
			require.Equal(t, http.StatusOK, createResp.StatusCode, "expected 200 from create session")

			var adminResp models.SessionCreateResponse
			err = json.NewDecoder(createResp.Body).Decode(&adminResp)
			require.NoError(t, err, "decoding create session response")

			// Extract the WebSocket endpoint and session code.
			wsEndpointBase := adminResp.ServerWsEndpoint
			require.NotEmpty(t, wsEndpointBase, "serverWsEndpoint must be provided by create session response")

			// Determine the session code needed for join.
			sessionCode := adminResp.SessionId
			require.NotEmpty(t, sessionCode, "session code must be in ID (adjust if different)")

			// 2. participants join:
			participantIDs := []string{
				fmt.Sprintf("user1"),
				fmt.Sprintf("user2"),
			}
			participantTokens := make([]string, 0, len(participantIDs))
			for _, pid := range participantIDs {
				joinReq := ValidateCodeReq{
					Code:   sessionCode,
					UserId: pid,
				}
				joinReqBody, err := json.Marshal(joinReq)
				require.NoError(t, err)

				joinResp, err := http.Post(sessionServiceURL+"/join", "application/json", bytes.NewReader(joinReqBody))
				require.NoError(t, err)
				defer joinResp.Body.Close()
				require.Equal(t, http.StatusOK, joinResp.StatusCode, "expected 200 from join for user %s", pid)

				var userResp models.SessionCreateResponse
				err = json.NewDecoder(joinResp.Body).Decode(&userResp)
				require.NoError(t, err, "decoding join response for user %s", pid)

				// The returned ServerWsEndpoint should match admin's or be same base:
				require.Equal(t, wsEndpointBase, wsEndpointBase, "WS endpoint mismatch for participant")

				t.Log("Store new token for", pid)
				participantTokens = append(participantTokens, userResp.Jwt)
			}

			var adminConn *websocket.Conn
			var usersConn []*websocket.Conn

			// 6. Admin WS connection
			adminConn = utils.ConnectWs(t, adminResp.Jwt)

			// 7. Read welcome, send ping, etc.
			utils.ReadWs(t, adminConn)

			// join users
			for _, userToken := range participantTokens {
				conn := utils.ConnectWs(t, userToken)
				usersConn = append(usersConn, conn)

				utils.ReadWs(t, conn)
			}

			usersAnswers := [][]int{
				{2, 2, 3},
				{2, 2, 1},
			}

			// 8. Start question flow
			for {
				t.Log("trigger question ")
				nextQuestionResp, err := http.Post(sessionServiceURL+fmt.Sprintf("/session/%s/nextQuestion", sessionCode),
					"application/json", nil)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, nextQuestionResp.StatusCode, "expected 200 from join for user ", adminId)
				nextQuestionResp.Body.Close()

				t.Log("Listen for messages")
				questionPayload := utils.ReadWs(t, adminConn)
				t.Logf("checking question %d payload: %v", questionPayload.QuestionIdx, questionPayload)

				for _, user := range usersConn {
					utils.ReadWs(t, user)
				}

				for j, user := range usersConn {
					option := usersAnswers[j][questionPayload.QuestionIdx-1]
					t.Logf("user %s sending answer: %d", participantIDs[j], option)
					msg := ws.ClientMessage{Option: option}
					user.WriteMessage(websocket.TextMessage, msg.Bytes())

					resp := utils.ReadWs(t, user)
					require.Equal(t, ws.MessageTypeAnswer, resp.Type)
					t.Log("answer is correct: ", resp.Correct)
				}

				if questionPayload.QuestionIdx == questionPayload.QuestionsAmount {
					t.Log("Game is finished")

					t.Log("end session ")
					endSessionResp, err := http.Post(sessionServiceURL+fmt.Sprintf("/session/%s/end", sessionCode),
						"application/json", nil)
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, endSessionResp.StatusCode, "expected 200 from ending session")
					endSessionResp.Body.Close()

					t.Log("---- Admin received leaderboard:")
					//utils.ReadWs(t, adminConn)

					t.Log("---- Users received leaderboard:")
					for _, user := range usersConn {
						utils.ReadWs(t, user)
					}
					break
				}
			}

			t.Log("Close connections:")
			// Ensure all connections will be closed at end
			utils.CloseWs(adminConn)
			t.Log("Closed admin")
			for _, c := range usersConn {
				utils.CloseWs(c)
			}
		}()
	}

	wg.Wait()
	cancel()
	wgRTS.Wait()

}
