package rabbit

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"xxx/real_time/ws"
	"xxx/shared"
)

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

func (r *RealTimeRabbit) ConsumeSessionStart() {
	msgs, err := r.channel.Consume(
		r.SessionStartedQ.Name,
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

	go func() {
		for d := range msgs {
			var msg shared.RabbitSessionMsg
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				continue
			}

			fmt.Println(msg)
			ws.RegisterSession(msg)
		}
	}()

	select {}
}

func (r *RealTimeRabbit) ConsumeSessionEnd() {

}
