package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io/ioutil"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
	"xxx/SessionService/Storage/Redis"
	"xxx/SessionService/httpServer"
	"xxx/SessionService/models"
	"xxx/real_time/config"
	"xxx/shared"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func startRabbit(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	defenitionsAbs, err := filepath.Abs(filepath.Join("..", "..", "..", "rabbit", "definitions.json"))
	require.NoError(t, err)
	confAbs, err := filepath.Abs(filepath.Join("..", "..", "..", "rabbit", "rabbitmq.conf"))
	require.NoError(t, err)
	// 1. Start RabbitMQ container
	rabbitReq := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3-management",
		ExposedPorts: []string{"5672/tcp"},
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
		ExposedPorts: []string{"6379/tcp"},
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

func getEnvFilePath() string {
	envPath := filepath.Join("..", "..", "..", ".env") // сдвигаемся на 4 уровня вверх из integration_tests
	absPath, err := filepath.Abs(envPath)
	if err != nil {
		log.Fatal(err)
	}
	return absPath
}

func Test_HttpServerCreate(t *testing.T) {
	cwd, _ := os.Getwd()
	fmt.Println("Working dir:", cwd)

	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			t.Fatalf("could not load .env file: %v", err)
		}
	}

	host := os.Getenv("SESSION_SERVICE_HOST")
	port := os.Getenv("SESSION_SERVICE_PORT")

	rabbitC, rabbitURL := startRabbit(context.Background(), t)
	redisC, redisURL := startRedis(context.Background(), t)
	defer redisC.Terminate(context.Background())
	defer rabbitC.Terminate(context.Background())
	// 🧵 Канал для получения Rabbit-сообщения
	rabbitMsgChan := make(chan []byte, 1)
	go func() {
		msg := consumeSessionStartFromRabbit(t, rabbitURL)
		rabbitMsgChan <- msg
	}()

	log := setupLogger(envLocal)
	server, err := httpServer.InitHttpServer(log, host, port, rabbitURL, redisURL)
	if err != nil {
		t.Fatalf("error creating http server: %v", err)
	}
	go server.Start()
	time.Sleep(2 * time.Second) // дать серверу запуститься
	defer server.Stop()
	// 🔨 Делаем запрос на создание сессии
	SessionServiceUrl := fmt.Sprintf("http://%s:%s/sessionsMock", host, port)
	req := models.CreateSessionReq{
		UserName: "user1",
		QuizId:   "d2372184-dedf-42db-bcbd-d6bb15b0712b",
	}
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatal("error marshaling json:", err)
	}

	resp, err := http.Post(SessionServiceUrl, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatal("error making request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %v", err)
	}

	var token models.SessionCreateResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

	// ✅ Проверяем Redis
	RedisConfig := models.Config{
		Addr: redisURL,
	}
	redis, err := Redis.NewRedisClient(context.Background(), RedisConfig)
	if err != nil {
		t.Fatal("error creating Redis client:", err)
	}
	session, err := redis.LoadSession(token.SessionId)
	if err != nil {
		t.Fatal("error loading session from Redis:", err)
	}
	if len(session.ID) <= 1 {
		t.Fatal("session id should not be empty, sessionId:", session.ID)
	}

	// ✅ Проверяем сообщение из RabbitMQ
	select {
	case msg := <-rabbitMsgChan:
		var event shared.QuizMessage
		err := json.Unmarshal(msg, &event)
		if err != nil {
			t.Fatalf("invalid JSON in RabbitMQ message: %v", err)
		}
		if event.SessionId != token.SessionId {
			t.Errorf("unexpected SessionId in RabbitMQ: got %s, want %s", event.SessionId, token.SessionId)
		}
		t.Log(event)
	case <-time.After(10 * time.Second):
		t.Fatal("timeout: did not receive message from RabbitMQ")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	file, err := os.OpenFile("cmd/SessionService/session.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file")
	}
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
func consumeSessionStartFromRabbit(t *testing.T, rabbitURL string) []byte {
	const exchange = "session.events"
	const routingKey = "session.start"

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		t.Fatalf("RabbitMQ connect error: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("RabbitMQ channel error: %v", err)
	}

	err = ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Exchange declare error: %v", err)
	}

	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		t.Fatalf("Queue declare error: %v", err)
	}

	err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		t.Fatalf("Queue bind error: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, true, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to consume: %v", err)
	}

	select {
	case msg := <-msgs:
		return msg.Body
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for message")
		return nil
	}
}
