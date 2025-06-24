package tests

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
	"time"
	"xxx/SessionService/Rabbit"
	"xxx/shared"
)

func Test_PublishSessionStart(t *testing.T) {
	rabbit, err := Rabbit.NewRabbit("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Errorf("Failed to open Rabbit: %s", err)
	}
	done := make(chan shared.Session)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
		shared.SessionStartRoutingKey, // например, "session.start"
		shared.SessionExchange,        // например, "session.events"
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
	err = rabbit.PublishSessionStart(context.Background(), shared.Session{
		ID:               "Start",
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
