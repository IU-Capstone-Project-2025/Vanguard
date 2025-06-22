package rabbit

// This file stores functions related to "session"-type events (start new session/cancel session) published to RabbitMQ

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"xxx/real_time/ws"
	"xxx/shared"
)

// CreateSessionStartedQueue declares and binds the `session_start` queue in RabbitMQ.
// The queue is utilized to receive events "start new session".
// Returns the queue object itself, or the error if failed.
func CreateSessionStartedQueue(ch *amqp.Channel) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		"session_started",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	err = ch.QueueBind(
		queue.Name,
		shared.SessionStartRoutingKey,
		shared.SessionExchange,
		false,
		nil)

	if err != nil {
		return amqp.Queue{}, err
	}
	return queue, nil
}

// CreateSessionEndedQueue declares and binds the `session_end` queue in RabbitMQ.
// The queue is utilized to receive events "cancel existing session".
// Returns the queue object itself, or the error if failed.
func CreateSessionEndedQueue(ch *amqp.Channel) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		"session_ended",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	err = ch.QueueBind(
		queue.Name,
		shared.SessionEndRoutingKey,
		shared.SessionExchange,
		false,
		nil)

	if err != nil {
		return amqp.Queue{}, err
	}
	return queue, nil
}

// ConsumeSessionStart method listens to "session start" events delivered to the corresponding queue.
func (r *RealTimeRabbit) ConsumeSessionStart(registry *ws.ConnectionRegistry) {
	msgs, err := r.channel.Consume(
		r.SessionStartedQ.Name, // the name of the already created queue
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	// listen to messages in parallel goroutine
	go func() {
		defer wg.Done()
		for d := range msgs {
			var msg shared.RabbitSessionMsg
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				continue
			}

			fmt.Println(msg)
			registry.RegisterSession(msg.SessionId) // register new session
		}
	}()

	wg.Wait() // defer this function termination while consuming from the queue
}

func (r *RealTimeRabbit) ConsumeSessionEnd() {

}
