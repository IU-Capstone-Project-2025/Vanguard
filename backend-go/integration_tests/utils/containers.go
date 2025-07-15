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
	"xxx/real_time/config"
)

func StartRabbit(ctx context.Context, t *testing.T) (addr string, terminate func()) {
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

	addr = fmt.Sprintf("amqp://%s:%s@%s:%s/", config.LoadConfig().MQ.User, config.LoadConfig().MQ.Password,
		rabbitHost, rabbitPort.Port())
	t.Logf("Rabbit running at %s", addr)
	terminate = func() {
		err := rabbitC.Terminate(ctx)
		require.NoError(t, err)
	}

	return addr, terminate
}

func StartRedis(ctx context.Context, t *testing.T) (addr string, terminate func()) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := redisC.Host(ctx)
	require.NoError(t, err)
	port, err := redisC.MappedPort(ctx, "6379")
	require.NoError(t, err)

	addr = fmt.Sprintf("redis://%s:%s", host, port.Port())
	terminate = func() {
		err := redisC.Terminate(ctx)
		require.NoError(t, err)
	}
	return addr, terminate
}
