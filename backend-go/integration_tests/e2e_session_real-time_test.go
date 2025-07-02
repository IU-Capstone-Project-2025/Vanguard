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
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
	"xxx/SessionService/httpServer"
	"xxx/SessionService/models"
	"xxx/real_time/app"
	"xxx/real_time/config"
	"xxx/real_time/rabbit"
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
	if os.Getenv("ENV") != "production" {
		// Assume .env is at project root; adapt path as needed
		if err := godotenv.Load("../../.env"); err != nil {
			// Not fatal if .env missing, but log
			t.Logf("No .env loaded: %v", err)
		}
	}
}

func startRabbit(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	defenitionsAbs, err := filepath.Abs(filepath.Join("..", "..", "rabbit", "definitions.json"))
	require.NoError(t, err)
	confAbs, err := filepath.Abs(filepath.Join("..", "..", "rabbit", "rabbitmq.conf"))
	require.NoError(t, err)

	// 1. Start RabbitMQ container
	rabbitReq := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3-management",
		ExposedPorts: []string{"5672:5672/tcp", "15672:15672/tcp"},
		Env: map[string]string{
			"RABBITMQ_LOAD_DEFINITIONS": "true",
			"RABBITMQ_DEFINITIONS_FILE": "/etc/rabbitmq/definitions.json",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      defenitionsAbs, // will be discarded internally
				ContainerFilePath: "/etc/rabbitmq/definitions.json",
				FileMode:          644,
			},

			{
				HostFilePath:      confAbs, // will be discarded internally
				ContainerFilePath: "/etc/rabbitmq/rabbitmq.conf",
				FileMode:          644,
			},
		},
		WaitingFor: wait.ForLog("Server startup complete"),
	}
	rabbitC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: rabbitReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start RabbitMQ container: %v", err)
	}

	rabbitHost, err := rabbitC.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	rabbitPort, err := rabbitC.MappedPort(ctx, "5672")
	if err != nil {
		t.Fatal(err)
	}

	u := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.LoadConfig().MQ.User, config.LoadConfig().MQ.Password,
		rabbitHost, rabbitPort.Port())
	t.Logf("Rabbit running at %s", u)
	return rabbitC, u
}

func startRedis(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine", // or "redis:latest"
		ExposedPorts: []string{"6379:6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(30 * time.Second),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start Redis container: %v", err)
	}

	host, err := redisC.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get Redis container host: %v", err)
	}
	mappedPort, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("failed to get Redis mapped port: %v", err)
	}
	addr := fmt.Sprintf("%s:%s", host, mappedPort.Port())
	t.Logf("Started Redis container at %s", addr)
	return redisC, addr
}

func TestSessionServiceToRealTime_E2E(t *testing.T) {
	loadEnv(t)
	ctx := context.Background()

	_, amqpUrl := startRabbit(ctx, t)
	_, redisUrl := startRedis(ctx, t)

	go func() {
		startRealTimeService(t, amqpUrl)
	}()
	t.Log("------------ wait for real time service -----------------")
	time.Sleep(2 * time.Second)

	go func() {
		startSessionService(t, amqpUrl, redisUrl)
	}()
	t.Log("------------ wait for session service -----------------")
	time.Sleep(2 * time.Second)

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

	// 2. Two participants join:
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
	adminConn = connectWs(t, adminResp.Jwt)

	// 7. Read welcome, send ping, etc.
	readWs(t, adminConn)

	// join users
	for _, userToken := range participantTokens {
		conn := connectWs(t, userToken)
		usersConn = append(usersConn, conn)

		readWs(t, conn)
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
		questionPayload := readWs(t, adminConn)
		t.Logf("checking question %d payload: %v", questionPayload.QuestionIdx, questionPayload)

		if questionPayload.QuestionIdx == 1 {
			for _, user := range usersConn {
				readWs(t, user)
			}
		}

		for j, user := range usersConn {
			option := usersAnswers[j][questionPayload.QuestionIdx-1]
			t.Logf("user %s sending answer: %d", participantIDs[j], option)
			msg := ws.ClientMessage{Option: option}
			user.WriteMessage(websocket.TextMessage, msg.Bytes())

			resp := readWs(t, user)
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
			readWs(t, adminConn)

			t.Log("---- Users received leaderboards:")
			for _, user := range usersConn {
				readWs(t, user)
			}
			break
		}
	}

	// Ensure all connections will be closed at end
	defer func() {
		adminConn.Close()
		for _, c := range usersConn {
			c.Close()
		}
	}()
}

func readWs(t *testing.T, conn *websocket.Conn) ws.ServerMessage {
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)

	var serverMsg ws.ServerMessage
	err = json.Unmarshal(msg, &serverMsg)

	t.Logf("Received: %s", msg)
	return serverMsg
}

func connectWs(t *testing.T, token string) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws", RawQuery: "token=" + token}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("WebSocket dial failed: %v", err)
	}

	return conn
}

func startRealTimeService(t *testing.T, amqpUrl string) {
	cfg := config.LoadConfig()

	manager := app.NewManager()

	// Connect to the rabbit MQ
	t.Log("Connecting to broker")
	err := manager.ConnectRabbitMQ(amqpUrl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Connected to broker")

	handlerDeps := ws.HandlerDeps{
		Tracker:  manager.QuizTracker,
		Registry: manager.ConnectionRegistry,
	}

	// SetCurrQuestionIdx route handler
	http.Handle("/ws", ws.NewWebSocketHandler(handlerDeps))

	go func() {
		err := http.ListenAndServe(
			fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			nil)
		t.Fatal(err)
	}()

	broker, err := rabbit.NewRealTimeRabbit(manager.Rabbit)
	t.Log("Service is up!")

	go broker.ConsumeSessionStart(manager.ConnectionRegistry, manager.QuizTracker)
	go broker.ConsumeSessionEnd(manager.ConnectionRegistry, manager.QuizTracker)
	select {}
}

func startSessionService(t *testing.T, amqpUrl, redisUrl string) {
	host := os.Getenv("SESSION_SERVICE_HOST")
	port := os.Getenv("SESSION_SERVICE_PORT")

	log := setupLogger()
	server, err := httpServer.InitHttpServer(log, host, port, amqpUrl, redisUrl)
	if err != nil {
		t.Fatal("error creating http server", "error", err)
		return
	}
	server.Start()
}

func setupLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	return log
}
