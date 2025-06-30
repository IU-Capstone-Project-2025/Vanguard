//go:build integration

package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
	"xxx/SessionService/httpServer"
	"xxx/real_time/config"
	"xxx/real_time/rabbit"
	"xxx/real_time/ws"
	"xxx/shared"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
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
	time.Sleep(10 * time.Second)

	go func() {
		startSessionService(t, amqpUrl, redisUrl)
	}()
	t.Log("------------ wait for session service -----------------")
	time.Sleep(10 * time.Second)

	// 1. Create a new session as admin
	// Use a random userId for admin
	adminUserId := fmt.Sprintf("admin-%d", time.Now().UnixNano())
	createReq := CreateSessionReq{
		QuizId: "1",
		UserId: adminUserId,
	}
	createReqBody, err := json.Marshal(createReq)
	require.NoError(t, err)

	createResp, err := http.Post(sessionServiceURL+"/sessions", "application/json", bytes.NewReader(createReqBody))
	require.NoError(t, err)
	defer createResp.Body.Close()
	require.Equal(t, http.StatusOK, createResp.StatusCode, "expected 200 from create session")

	var adminTokenResp shared.UserToken
	err = json.NewDecoder(createResp.Body).Decode(&adminTokenResp)
	require.NoError(t, err, "decoding create session response")

	// Extract the WebSocket endpoint and session code.
	wsEndpointBase := adminTokenResp.ServerWsEndpoint
	require.NotEmpty(t, wsEndpointBase, "serverWsEndpoint must be provided by create session response")

	// Determine the session code needed for join.
	sessionCode := adminTokenResp.SessionId
	require.NotEmpty(t, sessionCode, "session code must be in ID (adjust if different)")

	// 2. Two participants join:
	participantIDs := []string{
		fmt.Sprintf("user1"),
		fmt.Sprintf("user2"),
	}
	participantTokens := make([]shared.UserToken, 0, 2)
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

		var userTokenResp shared.UserToken
		err = json.NewDecoder(joinResp.Body).Decode(&userTokenResp)
		require.NoError(t, err, "decoding join response for user %s", pid)

		// The returned ServerWsEndpoint should match admin's or be same base:
		require.Equal(t, wsEndpointBase, userTokenResp.ServerWsEndpoint, "WS endpoint mismatch for participant")

		participantTokens = append(participantTokens, userTokenResp)
	}

	// 3. Connect two WebSocket clients to Real-Time Service
	// Use goroutines to connect concurrently and store connections
	type wsClient struct {
		Conn *websocket.Conn
		ID   string
	}
	wsClients := make([]wsClient, 0, 2)
	dialTimeout := 5 * time.Second

	for i, tokResp := range participantTokens {
		wsURL := tokResp.ServerWsEndpoint

		claims := shared.UserToken{
			UserId:           tokResp.UserId,
			UserType:         tokResp.UserType,
			ServerWsEndpoint: tokResp.ServerWsEndpoint,
			SessionId:        tokResp.SessionId,
			Exp:              tokResp.Exp,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Unix(0, tokResp.Exp)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		rawJwt, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
		require.NoError(t, err)
		// Parse and dial
		u, err := url.Parse(wsURL)
		require.NoError(t, err, "invalid WS endpoint URL for client %d", i)
		q := u.Query()
		q.Set("token", rawJwt)
		u.RawQuery = q.Encode()

		// Attempt connection with timeout:
		var conn *websocket.Conn
		deadline := time.Now().Add(dialTimeout)
		for {
			conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
			if err == nil {
				break
			}
			if time.Now().After(deadline) {
				t.Fatalf("WebSocket dial failed for client %s: %v", tokResp.UserId, err)
			}
			time.Sleep(200 * time.Millisecond)
		}
		wsClients = append(wsClients, wsClient{Conn: conn, ID: tokResp.UserId})
		// Optionally: read any immediate welcome message:
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, welcomeMsg, err := conn.ReadMessage()
		if err == nil {
			t.Logf("Client %s received initial message: %s", tokResp.UserId, string(welcomeMsg))
		} else {
			t.Logf("Client %s: no initial message or read timeout: %v", tokResp.UserId, err)
		}
	}
	// Ensure all connections will be closed at end
	defer func() {
		for _, c := range wsClients {
			c.Conn.Close()
		}
	}()
}

func startRealTimeService(t *testing.T, amqpUrl string) {
	// Initialize ws connections registry
	registry := ws.NewConnectionRegistry()

	// Set route handler
	http.Handle("/ws", ws.NewWebSocketHandler(registry))

	cfg := config.LoadConfig()

	go func() {
		err := http.ListenAndServe(
			fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			nil)
		t.Fatal(err)
	}()

	t.Log("Connecting to broker: ", amqpUrl)

	brokerConn, err := amqp.Dial(amqpUrl)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("Connected to broker")
	broker, err := rabbit.NewRealTimeRabbit(brokerConn)
	go broker.ConsumeSessionStart(registry)
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
