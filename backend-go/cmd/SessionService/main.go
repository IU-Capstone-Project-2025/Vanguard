package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
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
	host := flag.String("host", "localhost", "HTTP server host")
	port := flag.String("port", "8000", "HTTP server port")
	rabbitMQ := flag.String("rabbitmq", "amqp://guest:guest@localhost:5672/", "RabbitMQ URL")
	redis := flag.String("redis", "localhost:6379", "Redis address")
	flag.Parse()
	//if err := godotenv.Load(); err != nil {
	//	fmt.Println("Warning: .env file not loaded:", err)
	//}
	//host := os.Getenv("SESSION_SERVICE_HOST")
	//port := os.Getenv("SESSION_SERVICE_PORT")
	//rabbitMQ := os.Getenv("RABBITMQ")
	//redis := os.Getenv("REDIS")
	//fmt.Println("env:")
	fmt.Println(*host, *port, *rabbitMQ, *redis)
	log := setupLogger(envLocal)
	server, err := httpServer.InitHttpServer(log, *host, *port, *rabbitMQ, *redis)
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
