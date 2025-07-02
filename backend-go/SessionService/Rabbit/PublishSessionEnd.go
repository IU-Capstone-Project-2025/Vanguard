package Rabbit

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"xxx/shared"
)

func (r *Rabbit) PublishSessionEnd(ctx context.Context, SessionCode string, payload interface{}) error {
	body, err := json.Marshal(SessionCode)
	if err != nil {
		return err
	}
	err = r.channel.PublishWithContext(ctx,
		shared.SessionExchange,      // exchange
		shared.SessionEndRoutingKey, // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}
	return nil
}
