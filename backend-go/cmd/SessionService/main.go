package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	_ "xxx/SessionService/docs"
	"xxx/SessionService/httpServer"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func getEnvFilePath() string {
	root, err := filepath.Abs("..")
	if err != nil {
		log.Fatal("failed to find project root dir")
	}
	return filepath.Join(root, ".env")
}

// @title           Пример API
// @version         1.0
// @description     Это пример API с gorilla/mux и swaggo
// @host            localhost:8081
// @BasePath        /
func main() {
	// Load environment variables file, if running in development
	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		fmt.Println("LOADING .ENV")
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			log.Fatalf("Error: could not load .env file: %v", err)
		}
	}
	host := os.Getenv("SESSION_SERVICE_HOST")
	port := os.Getenv("SESSION_SERVICE_PORT")

	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))

	redisUrl := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	fmt.Println("env:")
	fmt.Println(host, port, rabbitUrl, redisUrl)
	//time.Sleep(30 * time.Second)

	log := setupLogger(envLocal)
	server, err := httpServer.InitHttpServer(log, host, port, rabbitUrl, redisUrl)
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
