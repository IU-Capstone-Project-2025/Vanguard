package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"testing"
	"time"
	"xxx/SessionService/httpServer"
	"xxx/SessionService/models"
)

func Test_HttpServerHealth(t *testing.T) {
	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			t.Fatalf("could not load .env file: %v", err)
		}
	}

	host := os.Getenv("SESSION_SERVICE_HOST")
	port := os.Getenv("SESSION_SERVICE_PORT")
	rabbitC, rabbitURL := startRabbit(context.Background(), t)
	redisC, redisURL := startRedis(context.Background(), t)
	defer redisC.Terminate(context.Background())
	defer rabbitC.Terminate(context.Background())
	// Запуск канала RabbitMQ для question.{sessionID}.start
	// ⚙️ Создаем логгер и запускаем сервер
	log := setupLogger(envLocal)
	server, err := httpServer.InitHttpServer(log, host, port, rabbitURL, redisURL)
	if err != nil {
		t.Fatalf("error creating http server: %v", err)
	}
	go server.Start()
	time.Sleep(1 * time.Second)
	defer server.Stop()
	// 🛠️ Создаем сессию
	SessionServiceUrl := fmt.Sprintf("http://%s:%s/healthz", host, port)
	req := models.CreateSessionReq{
		UserId: "1",
		QuizId: "d2372184-dedf-42db-bcbd-d6bb15b0712b",
	}
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatal("error marshaling json:", err)
	}
	resp, err := http.Post(SessionServiceUrl, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatal("error making request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: got %d", resp.StatusCode)
	}
}
