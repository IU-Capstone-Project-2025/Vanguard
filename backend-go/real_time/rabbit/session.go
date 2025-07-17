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
func (r *RealTimeRabbit) ConsumeSessionStart(
	registry *ws.ConnectionRegistry, tracker *ws.QuizTracker, ready chan struct{}) {
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
		close(ready)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	// listen to messages in parallel goroutine
	fmt.Println("Listen for new messages in session.start queue")
	// ✅ Signal that the queue consumer is ready
	close(ready)
	go func() {
		defer wg.Done()
		for d := range msgs {
			fmt.Println("RECEIVED SESSION START")

			var msg shared.QuizMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				continue
			}

			fmt.Println("Rabbit msg from Real Time", msg)
			registered := registry.RegisterSession(msg.SessionId) // register new session
			if registered {
				tracker.NewSession(msg.SessionId, msg.Quiz)

				go r.ConsumeQuestionStart(registry, tracker, msg.SessionId)
			}

		}
	}()

	wg.Wait() // defer this function termination while consuming from the queue
}

// ConsumeSessionEnd method listens to "session end" events delivered to the corresponding queue.
func (r *RealTimeRabbit) ConsumeSessionEnd(registry *ws.ConnectionRegistry, tracker *ws.QuizTracker, ready chan struct{}) {
	msgs, err := r.channel.Consume(
		r.SessionEndedQ.Name, // the name of the already created queue
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
		close(ready)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	// listen to messages in parallel goroutine
	fmt.Println("Listen for new messages in session.end queue")
	// ✅ Signal that the queue consumer is ready
	close(ready)
	go func() {
		defer wg.Done()
		for d := range msgs {
			fmt.Println("RECEIVED SESSION END")

			var sessionId string
			if err := json.Unmarshal(d.Body, &sessionId); err != nil {
				continue
			}

			err = r.CleanupQuestionConsumer(sessionId)
			if err != nil {
				fmt.Println(err)
			}

			gameEndAck := ws.ServerMessage{
				Type: ws.MessageTypeEnd,
			}
			registry.BroadcastToSession(sessionId, gameEndAck.Bytes(), false)

			registry.UnregisterSession(sessionId) // unregister new session
		}
	}()

	wg.Wait() // defer this function termination while consuming from the queue
}
