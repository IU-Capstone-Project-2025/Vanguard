package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"xxx/real_time/rabbit"
)

func main() {
	brokerConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	broker, err := rabbit.NewRealTimeRabbit(brokerConn)
	go broker.ConsumeSessionStart()
	select {}
}
