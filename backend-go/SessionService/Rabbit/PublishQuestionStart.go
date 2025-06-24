package Rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"xxx/shared"
)

func (r *Rabbit) PublishQuestionStart(ctx context.Context, SessionCode string, payload interface{}) error {
	routingKey := fmt.Sprintf("question.%s.start", SessionCode)
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	err = r.channel.PublishWithContext(ctx,
		shared.SessionExchange, // exchange
		routingKey,             // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}
	return nil
}
