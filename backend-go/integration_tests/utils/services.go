package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"testing"
	"xxx/LeaderBoardService/HttpServer"
	"xxx/SessionService/httpServer"
	"xxx/real_time/app"
	"xxx/real_time/config"
	"xxx/real_time/ws"
)

func StartRealTimeServer(t *testing.T, wg *sync.WaitGroup, amqpUrl, redisUrl string) (ctxCancel context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.LoadConfig()
	manager := app.NewManager()

	// Connect to the rabbit MQ
	t.Log("Connecting to broker...")
	broker, err := manager.ConnectRabbitMQ(amqpUrl)

	if err != nil {
		t.Fatal(err)
	}
	t.Log("Connected to broker")

	t.Log("Connecting to Redis...")
	err = manager.ConnectRedis(redisUrl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Connected to Redis")

	handlerDeps := ws.HandlerDeps{
		Tracker:  manager.QuizTracker,
		Registry: manager.ConnectionRegistry,
	}
	mux := http.NewServeMux()
	mux.Handle("/ws", ws.NewWebSocketHandler(handlerDeps))

	srv := &http.Server{Addr: cfg.Host + ":" + cfg.Port, Handler: mux}
	go func() {
		t.Log("HTTP server starting")
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("http listen: %v", err)
		}
	}()

	sessionStartReady := make(chan struct{})
	sessionEndReady := make(chan struct{})

	go broker.ConsumeSessionStart(manager.ConnectionRegistry, manager.QuizTracker, sessionStartReady)
	go broker.ConsumeSessionEnd(manager.ConnectionRegistry, manager.QuizTracker, sessionEndReady)

	go func(t *testing.T, wg *sync.WaitGroup) {
		defer wg.Done()

		<-ctx.Done()
		t.Log("Shutting down HTTP server...")
		_ = srv.Shutdown(context.Background())
	}(t, wg)

	<-sessionStartReady
	<-sessionEndReady

	t.Log("Real-time server fully up")
	return cancel
}

func StartSessionService(t *testing.T, amqpUrl, redisUrl string) {
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

func StartLeaderBoardService(t *testing.T, redisPort string) {
	//host := os.Getenv("LEADERBOARD_SERVICE_HOST")
	//port := os.Getenv("LEADERBOARD_SERVICE_PORT")
	host := "localhost"
	port := "8082"

	redisUrl := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), redisPort)

	//time.Sleep(30 * time.Second)
	log := setupLogger()
	server, err := HttpServer.InitHttpServer(log, host, port, redisUrl)
	if err != nil {
		t.Error("error creating http server", "error", err)
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
