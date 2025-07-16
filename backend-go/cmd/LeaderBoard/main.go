package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"xxx/LeaderBoardService/HttpServer"
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
	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		fmt.Println("LOADING .ENV")
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			log.Fatalf("Error: could not load .env file: %v", err)
		}
	}
	host := os.Getenv("LEADERBOARD_SERVICE_HOST")
	port := os.Getenv("LEADERBOARD_SERVICE_PORT")

	redisUrl := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	//time.Sleep(30 * time.Second)
	log := setupLogger(envLocal)
	server, err := HttpServer.InitHttpServer(log, host, port, redisUrl)
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
