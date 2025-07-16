package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/net/context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
	"xxx/LeaderBoardService/HttpServer"
	"xxx/shared"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

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

func Test_GetLeaderBoardResults(t *testing.T) {
	cwd, _ := os.Getwd()
	fmt.Println("Working dir:", cwd)

	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			t.Fatalf("could not load .env file: %v", err)
		}
	}
	host := os.Getenv("LEADER_BOARD_HOST")
	port := os.Getenv("LEADER_BOARD_PORT")

	redisC, redisURL := startRedis(context.Background(), t)
	defer redisC.Terminate(context.Background())
	log := setupLogger(envLocal)
	server, err := HttpServer.InitHttpServer(log, host, port, redisURL)
	if err != nil {
		t.Fatalf("error creating http server: %v", err)
	}
	go server.Start()
	time.Sleep(2 * time.Second)
	defer server.Stop()
	reqURL := fmt.Sprintf("http://%s:%s/get-results", host, port)
	sessionData := shared.SessionAnswers{
		SessionCode:   "ABC123",
		OptionsAmount: 4,
		Answers: []shared.Answer{
			{
				UserId:    "user_001",
				Correct:   true,
				Answered:  true,
				Option:    "1",
				Timestamp: time.Now().Add(-2 * time.Second),
			},
			{
				UserId:    "user_002",
				Correct:   true,
				Answered:  true,
				Option:    "3",
				Timestamp: time.Now().Add(-1 * time.Second),
			},
			{
				UserId:    "user_003",
				Correct:   false,
				Answered:  false,
				Option:    "",
				Timestamp: time.Now().Add(-2 * time.Second),
			},
			{
				UserId:    "user_004",
				Correct:   true,
				Answered:  true,
				Option:    "3",
				Timestamp: time.Now(),
			},
		},
	}
	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		t.Fatalf("Error marshaling request data: %v", err)
	}

	// Выполнение POST-запроса
	resp, err := http.Post(reqURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	defer resp.Body.Close()

	// Проверка ответа
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", resp.StatusCode)
	}
	var response shared.BoardResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal("error decoding response:", err)
	}
	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("error marshalling response for log: %v", err)
	}
	t.Logf("Task status response:\n%s", data)
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
