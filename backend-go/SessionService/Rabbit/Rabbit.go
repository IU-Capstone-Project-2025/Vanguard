package Rabbit

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"xxx/shared"
)

type Rabbit struct {
	Conn    *amqp.Connection
	channel *amqp.Channel
}

type Broker interface {
	PublishQuestionStart(ctx context.Context, SessionCode string, payload interface{}) error
	PublishSessionEnd(ctx context.Context, SessionCode string, payload interface{}) error
	PublishSessionStart(ctx context.Context, payload interface{}) error
}

func NewRabbit(rmq_host string) (*Rabbit, error) {
	conn, err := amqp.Dial(rmq_host)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
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
		return nil, err
	}

	return &Rabbit{conn, ch}, nil
}
