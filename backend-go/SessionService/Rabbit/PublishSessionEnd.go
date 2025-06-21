package Rabbit

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"xxx/Shared"
)

func (r *Rabbit) PublishSessionEnd(ctx context.Context, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	err = r.channel.PublishWithContext(ctx,
		Shared.SessionExchange,      // exchange
		Shared.SessionEndRoutingKey, // routing key
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
