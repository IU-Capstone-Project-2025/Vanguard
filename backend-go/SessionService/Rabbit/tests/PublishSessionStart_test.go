package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"testing"
	"time"
	"xxx/SessionService/Rabbit"
	"xxx/shared"
)

func Test_PublishSessionStart(t *testing.T) {
	if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "test" {
		if err := godotenv.Load(getEnvFilePath()); err != nil {
			t.Fatalf("could not load .env file: %v", err)
		}
	}
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))
	rabbit, err := Rabbit.NewRabbit(rabbitURL)
	if err != nil {
		t.Fatalf("Failed to open Rabbit: %s", err)
	}
	done := make(chan shared.QuizMessage)
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
		shared.SessionStartRoutingKey, // например, "session.start"
		shared.SessionExchange,        // например, "session.events"
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
		var s shared.QuizMessage
		for m := range msgs {
			err := json.Unmarshal(m.Body, &s)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %s", err)
			}
			t.Logf("Received a message: %v", s)
			done <- s
			return
		}
	}()
	quiz := shared.Quiz{Questions: []shared.Question{
		{
			Type: "single_choice",
			Text: "What is the output of print(2 ** 3)?",
			Options: []shared.Option{
				{Text: "6", IsCorrect: false},
				{Text: "8", IsCorrect: true},
				{Text: "9", IsCorrect: false},
				{Text: "5", IsCorrect: false},
			},
		},
		{
			Type: "single_choice",
			Text: "Which keyword is used to create a function in Python?",
			Options: []shared.Option{
				{Text: "func", IsCorrect: false},
				{Text: "function", IsCorrect: false},
				{Text: "def", IsCorrect: true},
				{Text: "define", IsCorrect: false},
			},
		},
		{
			Type: "single_choice",
			Text: "What data type is the result of: 3 / 2 in Python 3?",
			Options: []shared.Option{
				{Text: "int", IsCorrect: false},
				{Text: "float", IsCorrect: true},
				{Text: "str", IsCorrect: false},
				{Text: "decimal", IsCorrect: false},
			},
		},
	}}
	err = rabbit.PublishSessionStart(context.Background(), shared.QuizMessage{
		SessionId: "1",
		Quiz:      quiz,
	})
	if err != nil {
		t.Fatalf("Failed to publish session start: %s", err)
	}
	var s shared.QuizMessage
	select {
	case s = <-done:
		fmt.Println("done", s.SessionId)
	case <-time.After(10 * time.Second):
		fmt.Println("Failed to get session")
		t.FailNow()
	}
}
