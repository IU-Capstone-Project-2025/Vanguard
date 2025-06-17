package Rabbit

import amqp "github.com/rabbitmq/amqp091-go"

type Rabbit struct {
	Conn    *amqp.Connection
	channel *amqp.Channel
}

type Broker interface {
	PublishEvent(interface{}) error
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
