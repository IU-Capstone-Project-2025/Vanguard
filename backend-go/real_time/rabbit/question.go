package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"xxx/shared"
)

func CreateQuestionStartQueue(ch *amqp.Channel) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		"question_start",
		false,
		true,
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

func (r *RealTimeRabbit) ConsumeQuestionStart() {

}
