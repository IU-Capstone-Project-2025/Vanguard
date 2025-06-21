package main

import (
	"fmt"
	"log/slog"
	"os"
	"xxx/SessionService/httpServer"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	log := setupLogger(envLocal)
	server, err := httpServer.InitHttpServer(log, "localhost", "8000", "amqp://guest:guest@localhost:5672/", "localhost:6379")
	if err != nil {
		log.Error("error creating http server", "error", err)
		return
	}
	server.Start()
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
