package rabbit

// This file stores functions related to "question"-type events published to RabbitMQ

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"xxx/real_time/ws"
	"xxx/shared"
)

// CreateQuestionStartQueue declares and binds the `question_start` queue in RabbitMQ.
// The queue is utilized to receive events "start next question".
// The queue is auto delete, since it is temporary and exists only till the session is alive.
// Returns the queue object itself, or the error if failed.
func CreateQuestionStartQueue(ch *amqp.Channel) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		"question_start",
		false,
		true, // auto delete
		false,
		true,
		nil)

	if err != nil {
		return amqp.Queue{}, err
	}

	err = ch.QueueBind(
		queue.Name,
		shared.QuestionStartRoutingKey,
		shared.SessionExchange,
		true,
		nil)
	if err != nil {
		return amqp.Queue{}, err
	}

	return queue, nil
}

// ConsumeQuestionStart method listens to "next question start" events delivered to the corresponding queue.
func (r *RealTimeRabbit) ConsumeQuestionStart(sessionId string, registry *ws.ConnectionRegistry) {
	msgs, err := r.channel.Consume(
		r.QuestionStartedQs[sessionId].Name, // the name of the already created queue
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
		for range msgs { // ignore the contents in the queue, since only event itself matters
			fmt.Println("next question triggered")
		}
	}()

	wg.Wait() // defer this function termination while consuming from the queue
}
