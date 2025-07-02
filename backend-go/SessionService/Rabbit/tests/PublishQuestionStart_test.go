package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
	"xxx/SessionService/Rabbit"
	"xxx/real_time/config"
	"xxx/shared"
)

func startRabbit(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	defenitionsAbs, err := filepath.Abs(filepath.Join("..", "..", "..", "..", "rabbit", "definitions.json"))
	require.NoError(t, err)
	confAbs, err := filepath.Abs(filepath.Join("..", "..", "..", "..", "rabbit", "rabbitmq.conf"))
	require.NoError(t, err)

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
				HostFilePath:      defenitionsAbs, // will be discarded internally
				ContainerFilePath: "/etc/rabbitmq/definitions.json",
				FileMode:          644,
			},

			{
				HostFilePath:      confAbs, // will be discarded internally
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

	rabbitHost, err := rabbitC.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	rabbitPort, err := rabbitC.MappedPort(ctx, "5672")
	if err != nil {
		t.Fatal(err)
	}

	u := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.LoadConfig().MQ.User, config.LoadConfig().MQ.Password,
		rabbitHost, rabbitPort.Port())
	t.Logf("Rabbit running at %s", u)
	return rabbitC, u
}

func startRedis(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine", // or "redis:latest"
		ExposedPorts: []string{"6379:6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(30 * time.Second),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start Redis container: %v", err)
	}

	host, err := redisC.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get Redis container host: %v", err)
	}
	mappedPort, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("failed to get Redis mapped port: %v", err)
	}
	addr := fmt.Sprintf("%s:%s", host, mappedPort.Port())
	t.Logf("Started Redis container at %s", addr)
	return redisC, addr
}

func getEnvFilePath() string {
	envPath := filepath.Join("..", "..", "..", "..", ".env") // сдвигаемся на 4 уровня вверх из integration_tests
	absPath, err := filepath.Abs(envPath)
	if err != nil {
		log.Fatal(err)
	}
	return absPath
}
func Test_PublishQuestionStart(t *testing.T) {
	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			t.Fatalf("could not load .env file: %v", err)
		}
	}
	_, rabbitURL := startRabbit(context.Background(), t)
	rabbit, err := Rabbit.NewRabbit(rabbitURL)
	if err != nil {
		t.Fatalf("Failed to open Rabbit: %s", err)
	}
	done := make(chan shared.Session)
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(
		shared.SessionExchange, // имя
		"topic",                // тип
		true,                   // durable
		false,                  // auto-delete
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		t.Fatalf("Failed to declare an exchange: %s", err)
	}
	q, err := ch.QueueDeclare(
		"",    // пустое имя = сгенерировать уникальное
		false, // durable
		true,  // auto-delete
		true,  // exclusive
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to declare a queue: %s", err)
	}
	err = ch.QueueBind(
		q.Name,
		shared.QuestionStartRoutingKey, // например, "session.start"
		shared.SessionExchange,         // например, "session.events"
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to bind a queue: %s", err)
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	go func() {
		var s shared.Session
		for m := range msgs {
			err := json.Unmarshal(m.Body, &s)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %s", err)
			}
			done <- s
			return
		}
	}()
	err = rabbit.PublishQuestionStart(context.Background(), "Abc123", shared.Session{
		ID:               "Abc123",
		Code:             "123",
		State:            "123",
		ServerWsEndpoint: "123",
	})
	if err != nil {
		t.Fatalf("Failed to publish session start: %s", err)
	}
	var s shared.Session
	select {
	case s = <-done:
		fmt.Println("done", s.ID)
	case <-time.After(10 * time.Second):
		fmt.Println("Failed to get session")
		t.FailNow()
	}
}
