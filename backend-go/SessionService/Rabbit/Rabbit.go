package Rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	Conn    *amqp.Connection
	channel *amqp.Channel
}

type Broker interface {
	PublishEvent(interface{}) error
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
	return &Rabbit{conn, ch}, nil
}

// PublishEvent will publish to brocker some event
func (r *Rabbit) PublishEvent(payload interface{}) error {
	//TODO implement function
	return nil
}

// SessionCreated will sent message that session created
func (r *Rabbit) SessionCreated(payload interface{}) error {
	//TODO implement function
	return nil
}
