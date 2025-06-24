package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"time"
	_ "xxx/SessionService/docs"
	"xxx/SessionService/httpServer"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title           Пример API
// @version         1.0
// @description     Это пример API с gorilla/mux и swaggo
// @host            localhost:8000
// @BasePath        /
func main() {
	time.Sleep(30 * time.Second)
	_ = godotenv.Load()
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	rabbitMQ := os.Getenv("RABBITMQ")
	redis := os.Getenv("REDIS")
	fmt.Println(host, port, rabbitMQ, redis)
	log := setupLogger(envLocal)
	server, err := httpServer.InitHttpServer(log, host, port, rabbitMQ, redis)
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
