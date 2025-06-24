//go:build integration

package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
	"xxx/real_time/config"
	"xxx/real_time/rabbit"
	"xxx/real_time/ws"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithTestContainers(t *testing.T) {
	ctx := context.Background()

	if os.Getenv("ENV") != "production" {
		fmt.Println("LOADING .ENV")
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			log.Fatalf("Error: could not load .env file: %v", err)
		}
	}

	defenitionsAbs, err := filepath.Abs(filepath.Join("..", "..", "..", "rabbit", "definitions.json"))
	require.NoError(t, err)
	confAbs, err := filepath.Abs(filepath.Join("..", "..", "..", "rabbit", "rabbitmq.conf"))
	require.NoError(t, err)

	fmt.Println(defenitionsAbs, confAbs)

	// 1. Start RabbitMQ container
	rabbitReq := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3-management",
		ExposedPorts: []string{"5672:5672/tcp", "15672:15672/tcp"},
		Env: map[string]string{
			"RABBITMQ_LOAD_DEFINITIONS": "true",
			"RABBITMQ_DEFINITIONS_FILE": "/etc/rabbitmq/definitions.json",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      defenitionsAbs,
				ContainerFilePath: "/etc/rabbitmq/definitions.json",
				FileMode:          644,
			},

			{
				HostFilePath:      confAbs,
				ContainerFilePath: "/etc/rabbitmq/rabbitmq.conf",
				FileMode:          644,
			},
		},
		WaitingFor: wait.ForLog("Server startup complete"),
	}
	rabbitC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: rabbitReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start RabbitMQ container: %v", err)
	}
	defer rabbitC.Terminate(ctx)

	rabbitHost, err := rabbitC.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	rabbitPort, err := rabbitC.MappedPort(ctx, "5672")
	if err != nil {
		t.Fatal(err)
	}
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.LoadConfig().MQ.User, config.LoadConfig().MQ.Password,
		rabbitHost, rabbitPort.Port())

	fmt.Println("rabbit running on ", amqpURL)

	// 2. Start Redis container similarly if needed

	// 3. Start RealTime service in a goroutine or exec.Command, configuring it to connect to amqpURL and Redis.
	//    For brevity, assume RealTime service can be started in-process or as a subprocess, reading env vars:
	// Start your RealTime main in a goroutine if possible, or exec binary.
	go func() {
		if err := startRealTimeServer(); err != nil {
			log.Fatal(err)
		}
	}()

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJTZXNzaW9uSWQiOiJDT0Q1MiIsIlJvbGUiOiJwYXJ0aWNpcGFudCIsIlVzZXJJRCI6ImdhbWFyam9iYS04PT09PT09RCJ9.bQnfDAs-Z2THdSyR1a0hzbe4PliKJB9fartXA7lKPI8"

	// 4. Wait for service readiness (e.g., dial WS until success)
	deadline := time.Now().Add(30 * time.Second)
	var wsConn *websocket.Conn
	if time.Now().After(deadline) {
		t.Fatal("WebSocket endpoint not ready in time")
	}
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws", RawQuery: "token=" + token}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err == nil {
		conn.Close()
	}
	time.Sleep(500 * time.Millisecond)

	// 5. Publish session.start
	rabCon, err := amqp.Dial(amqpURL)
	if err != nil {
		t.Fatalf("Dial RabbitMQ: %v", err)
	}
	ch, err := rabCon.Channel()
	if err != nil {
		t.Fatalf("Open channel: %v", err)
	}
	exchange := "session.events"
	ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	evt := map[string]string{"SessionId": "COD52"}
	body, _ := json.Marshal(evt)
	ch.Publish(exchange, "session.start", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	rabCon.Close()

	// 6. Wait a bit, then connect WS with a valid JWT for session "sess123"
	time.Sleep(1 * time.Second)
	u = url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws", RawQuery: "token=" + token}
	wsConn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("WebSocket dial failed: %v", err)
	}
	defer wsConn.Close()

	// 7. Read welcome, send ping, etc.
	wsConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, msg, err := wsConn.ReadMessage()
	assert.NoError(t, err)
	t.Logf("Received: %s", msg)
}

func getEnvFilePath() string {
	root, err := filepath.Abs("../../..")
	if err != nil {
		log.Fatal("failed to find project root dir")
	}
	return filepath.Join(root, ".env")
}

func startRealTimeServer() error {
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
