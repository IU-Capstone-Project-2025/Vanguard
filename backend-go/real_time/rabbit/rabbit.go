package rabbit

// This file stores functions related to general settings of RabbitMQ and its set up process

import amqp "github.com/rabbitmq/amqp091-go"

// RealTimeRabbit manages all manipulations with RabbitMQ within the Real-Time Service.
// Stores connection and channel objects and existing queues. Has methods to consume queues.
type RealTimeRabbit struct {
	conn              *amqp.Connection
	channel           *amqp.Channel
	SessionStartedQ   amqp.Queue            // For events from Session service for new started session
	SessionEndedQ     amqp.Queue            // For events from Session service for closed session
	QuestionStartedQs map[string]amqp.Queue // For events from Session service for starting next question.
	// Maps session code to the corresponding queue in format 'session code': 'question_start queue'
}

// NewRealTimeRabbit initializes RealTimeRabbit object. Given the connection [conn]
func NewRealTimeRabbit(conn *amqp.Connection) (*RealTimeRabbit, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	rabbit := &RealTimeRabbit{
		conn:              conn,
		channel:           ch,
		QuestionStartedQs: make(map[string]amqp.Queue), // initialize the empty map
	}

	// -------- CREATE QUEUES --------
	// create "session_start" queue
	rabbit.SessionStartedQ, err = CreateSessionStartedQueue(ch)
	if err != nil {
		return nil, err
	}

	// create "session_end" queue
	rabbit.SessionEndedQ, err = CreateSessionEndedQueue(ch)
	if err != nil {
		return nil, err
	}

	return rabbit, nil
}
