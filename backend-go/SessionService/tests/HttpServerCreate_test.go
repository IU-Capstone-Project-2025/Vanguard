package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"
	"xxx/SessionService/Storage/Redis"
	"xxx/SessionService/httpServer"
	"xxx/SessionService/models"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Test_HttpServerCreate(t *testing.T) {
	log := setupLogger(envLocal)

	server, err := httpServer.InitHttpServer(log, "localhost", "8000", "amqp://guest:guest@localhost:5672/", "localhost:6379")
	if err != nil {
		log.Error("error creating http server", "error", err)
		return
	}

	go server.Start()
	time.Sleep(1 * time.Second) // Даем серверу стартануть

	resp, err := http.Get("http://localhost:8000/create?userId=1")
	if err != nil {
		t.Error("error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("error reading response body: %s", err.Error())
		return
	}

	if len(body) == 0 {
		t.Errorf("response body is empty")
		return
	}

	fmt.Println("get body:")
	fmt.Println(string(body))

	var token models.UserToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		t.Errorf("error unmarshalling response: %s", err.Error())
		return
	}

	RedisConfig := models.Config{
		Addr:        "localhost:6379",
		Password:    "",
		DB:          0,
		MaxRetries:  0,
		DialTimeout: 0,
		Timeout:     0,
	}
	redis, err := Redis.NewRedisClient(context.Background(), RedisConfig)
	if err != nil {
		t.Error("error creating Redis client:", err)
		return
	}

	session, err := redis.LoadSession(token.CurrentQuiz)
	if err != nil {
		t.Error("error loading session from Redis:", err)
		return
	}

	fmt.Println("loaded session from Redis:")
	fmt.Println(session)
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
