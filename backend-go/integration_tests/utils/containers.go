package utils

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
	"runtime"
	"testing"
	"time"
	"xxx/real_time/config"
)

func StartRabbit(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path")

	baseDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "rabbit")
	definitionsPath := filepath.Join(baseDir, "definitions.json")
	confPath := filepath.Join(baseDir, "rabbitmq.conf")

	// Use absolute paths
	definitionsAbs, err := filepath.Abs(definitionsPath)
	require.NoError(t, err)
	confAbs, err := filepath.Abs(confPath)
	require.NoError(t, err)

	// 1. Start RabbitMQ container
	rabbitReq := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3-management",
		ExposedPorts: []string{"5672/tcp"},
		Env: map[string]string{
			"RABBITMQ_LOAD_DEFINITIONS": "true",
			"RABBITMQ_DEFINITIONS_FILE": "/etc/rabbitmq/definitions.json",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      definitionsAbs, // will be discarded internally
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

func StartRedis(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine", // or "redis:latest"
		ExposedPorts: []string{"6379/tcp"},
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
