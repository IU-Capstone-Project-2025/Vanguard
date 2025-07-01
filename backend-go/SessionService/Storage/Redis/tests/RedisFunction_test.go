package tests

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
	"xxx/SessionService/Storage/Redis"
	"xxx/SessionService/models"
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
	envPath := filepath.Join("..", "..", "..", "..", "..", ".env") // сдвигаемся на 4 уровня вверх из integration_tests
	absPath, err := filepath.Abs(envPath)
	if err != nil {
		log.Fatal(err)
	}
	return absPath
}
func Test_RedisFunction(t *testing.T) {
	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			t.Fatalf("could not load .env file: %v", err)
		}
	}
	redisC, redisUrl := startRedis(context.Background(), t)
	defer redisC.Terminate(context.Background())
	ctx := context.Background()
	RedisConfig := models.Config{
		Addr:        redisUrl,
		Password:    "",
		DB:          0,
		MaxRetries:  0,
		DialTimeout: 0,
		Timeout:     0,
	}
	redis, err := Redis.NewRedisClient(ctx, RedisConfig)
	if err != nil {
		t.Fatal("error creating redis client", "error", err)
	}
	err = redis.SaveSession(&shared.Session{
		ID:               "testRedis",
		Code:             "testRedis",
		State:            "testRedis",
		ServerWsEndpoint: "testRedis",
	})
	if err != nil {
		t.Fatal("error creating redis client", "error", err)
	}
	fmt.Println("SaveSession success")
	session, err := redis.LoadSession("testRedis")
	if err != nil {
		t.Fatal("error loading session", "error", err)
	}
	if session.Code != "testRedis" {
		t.Error("error loading session", "error", err)
	}
	fmt.Println("LoadSession success", session)
	err = redis.EditSessionState("testRedis", "testRedis2")
	if err != nil {
		t.Fatal("error editing session", "error", err)
	}
	session, err = redis.LoadSession("testRedis")
	if err != nil {
		t.Fatal("error loading session", "error", err)
	}
	if session.State != "testRedis2" {
		t.Error("error loading session", "error", err)
	}
	fmt.Println("EditSessionState success", session)
	err = redis.DeleteSession("testRedis")
	if err != nil {
		t.Fatal("error deleting session", "error", err)
	}
	session, err = redis.LoadSession("testRedis")
	if err != nil {
		fmt.Println("DeleteSession success")
	}

}
