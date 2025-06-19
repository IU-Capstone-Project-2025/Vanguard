package rabbit

import amqp "github.com/rabbitmq/amqp091-go"

type RealTimeRabbit struct {
	conn              *amqp.Connection
	channel           *amqp.Channel
	SessionStartedQ   amqp.Queue            // For events from Session service for new started session
	SessionEndedQ     amqp.Queue            // For events from Session service for closed session
	QuestionStartedQs map[string]amqp.Queue // For events from Session service for starting next question
	// in format 'session code': 'question_start queue'
}

func NewRealTimeRabbit(conn *amqp.Connection) (*RealTimeRabbit, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	rabbit := &RealTimeRabbit{
		conn:              conn,
		channel:           ch,
		QuestionStartedQs: make(map[string]amqp.Queue),
	}

	rabbit.SessionStartedQ, err = CreateSessionStartedQueue(ch)
	if err != nil {
		return nil, err
	}

	rabbit.SessionEndedQ, err = CreateSessionEndedQueue(ch)
	if err != nil {
		return nil, err
	}

	return rabbit, nil
}
