package main

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"xxx/real_time/app"
	"xxx/real_time/config"
	"xxx/real_time/ws"
)

func getEnvFilePath() string {
	root, err := filepath.Abs("..")
	if err != nil {
		log.Fatal("failed to find project root dir")
	}
	return filepath.Join(root, ".env")
}

func main() {
	// Load environment variables file, if running in development
	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		fmt.Println("LOADING .ENV from ", getEnvFilePath())
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			log.Fatalf("Error: could not load .env file: %v", err)
		}
	}

	cfg := config.LoadConfig()

	manager := app.NewManager(cfg.LB.Host, cfg.LB.Port)

	// Connect to the rabbit MQ
	fmt.Println("Connecting to broker...")
	broker, err := manager.ConnectRabbitMQ(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.MQ.User, cfg.MQ.Password, cfg.MQ.Host, cfg.MQ.Port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to broker")

	fmt.Println("Connecting to Redis...")
	err = manager.ConnectRedis(fmt.Sprintf("redis://%s:%s", cfg.Redis.Host, cfg.Redis.Port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Redis")

	handlerDeps := ws.HandlerDeps{
		Tracker:  manager.QuizTracker,
		Registry: manager.ConnectionRegistry,
	}

	// SetCurrQuestionIdx route handler
	http.Handle("/ws", ws.NewWebSocketHandler(handlerDeps))

	go func() {
		err := http.ListenAndServe(
			fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			nil)
		log.Fatal(err)
	}()

	fmt.Println("Service is up!")

	sessionStartReady := make(chan struct{})
	sessionEndReady := make(chan struct{})

	go broker.ConsumeSessionStart(manager.ConnectionRegistry, manager.QuizTracker, sessionStartReady)
	go broker.ConsumeSessionEnd(manager.ConnectionRegistry, manager.QuizTracker, sessionEndReady)

	<-sessionStartReady
	<-sessionEndReady
	select {}
}
