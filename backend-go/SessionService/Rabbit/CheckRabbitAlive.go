package Rabbit

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

func (r *Rabbit) CheckRabbitAlive() error {
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))
	conn, err := amqp.Dial(rabbitUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return nil
}
