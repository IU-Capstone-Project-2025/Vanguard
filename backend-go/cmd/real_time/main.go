package main

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"xxx/real_time/config"
	"xxx/real_time/rabbit"
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
	if os.Getenv("ENV") != "production" {
		fmt.Println("LOADING .ENV")
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			log.Fatalf("Error: could not load .env file: %v", err)
		}
	}

	// Initialize ws connections registry
	registry := ws.NewConnectionRegistry()

	// Set route handler
	http.Handle("/ws", ws.NewWebSocketHandler(registry))

	cfg := config.LoadConfig()

	go func() {
		err := http.ListenAndServe(
			fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			nil)
		log.Fatal(err)
	}()

	fmt.Println("Connecting to broker")
	brokerConn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s/",
			cfg.MQ.User, cfg.MQ.Password, cfg.MQ.Host, cfg.MQ.Port),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to broker")
	broker, err := rabbit.NewRealTimeRabbit(brokerConn)
	go broker.ConsumeSessionStart(registry)
	select {}
}
