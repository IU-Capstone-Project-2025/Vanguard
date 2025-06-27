package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
	"xxx/SessionService/Rabbit"
	"xxx/shared"
)

func getEnvFilePath() string {
	envPath := filepath.Join("..", "..", "..", "..", ".env") // сдвигаемся на 4 уровня вверх из integration_tests
	absPath, err := filepath.Abs(envPath)
	if err != nil {
		log.Fatal(err)
	}
	return absPath
}
func Test_PublishQuestionStart(t *testing.T) {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			t.Fatalf("could not load .env file: %v", err)
		}
	}
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))

	rabbit, err := Rabbit.NewRabbit(rabbitURL)
	if err != nil {
		t.Errorf("Failed to open Rabbit: %s", err)
	}
	done := make(chan shared.Session)
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		t.Errorf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		t.Errorf("Failed to open a channel: %s", err)
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
		t.Errorf("Failed to declare an exchange: %s", err)
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
		t.Errorf("Failed to declare a queue: %s", err)
	}
	err = ch.QueueBind(
		q.Name,
		shared.QuestionStartRoutingKey, // например, "session.start"
		shared.SessionExchange,         // например, "session.events"
		false,
		nil,
	)
	if err != nil {
		t.Errorf("Failed to bind a queue: %s", err)
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
				t.Errorf("Failed to unmarshal: %s", err)
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
		t.Errorf("Failed to publish session start: %s", err)
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
