package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
	"xxx/SessionService/httpServer"
	"xxx/SessionService/models"
	"xxx/shared"
)

func Test_HttpServerNextQuestion(t *testing.T) {
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
	rabbitMsgChan := make(chan []byte, 1)
	go func() {
		msg := consumeQuestionStartFromRabbit(t, rabbitURL, "123")
		rabbitMsgChan <- msg
	}()

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
	SessionServiceUrl := fmt.Sprintf("http://%s:%s/sessionsMock", host, port)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %v", err)
	}

	var sessionResp models.SessionCreateResponse
	if err := json.Unmarshal(body, &sessionResp); err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
	sessionID := sessionResp.SessionId

	// 📥 Начинаем слушать RabbitMQ по ключу question.{sessionID}.start

	// 📤 Отправляем POST /session/{id}/nextQuestion
	nextQuestionUrl := fmt.Sprintf("http://%s:%s/session/%s/nextQuestion", host, port, sessionID)
	resp2, err := http.Post(nextQuestionUrl, "application/json", nil)
	if err != nil {
		t.Fatalf("error sending nextQuestion request: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: got %d", resp2.StatusCode)
	}

	// ✅ Проверка сообщения из RabbitMQ
	select {
	case msg := <-rabbitMsgChan:
		fmt.Println("📦 Message received from RabbitMQ:", string(msg))
	case <-time.After(10 * time.Second):
		t.Fatal("timeout: did not receive message from RabbitMQ on question.{sessionID}.start")
	}
}

func consumeQuestionStartFromRabbit(t *testing.T, rabbitURL, sessionID string) []byte {
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
	select {
	case msg := <-msgs:
		return msg.Body
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for question.*.start message")
		return nil
	}
}
